package docker

import (
	"os"
	"fmt"
	"context"
	"strings"
	"io/ioutil"
	"errors"
	"log"

	"gnative/config"
	
	"github.com/jhoonb/archivex"
	"github.com/docker/go-connections/nat"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/pkg/jsonmessage"
)

const TMP_FILE = "/tmp/build-context.tar"
const defaultDockerAPIVersion = "v1.40"

func Build(conf config.Config) error {
	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		return err
	}

	err = buildImage(cli, conf.Source.Libraries, BASE_IMAGE_NAME, getBaseGoBuildDockerfile())
	if err != nil { return err }

	libdirs, err := getDirectories(conf.Source.Libraries)
	if err != nil { return err }

	files, err := getDirectories(conf.Source.Images)
	if err != nil { return err }

	for _, file := range files {
		contents, err := ioutil.ReadDir(conf.Source.Images + "/" + file)
		if (err != nil) { return err }

		moduleName := findModule(contents)
		if (moduleName == "") { return errors.New("Folder " + file + " in " + conf.Source.Images + " does not contain a go module") }

		err = buildImage(cli, conf.Source.Images + "/" + file, conf.Registry + "/" + file + ":0.0.1", getGoDockerfile(libdirs, conf.Registry))
		if err != nil { return err }
	} 

	return nil
}

func getDirectories(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil { return nil, err }

	dirs := []string{}

	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	return dirs, nil
}

func findModule(files []os.FileInfo) string {
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mod") {
			return file.Name()
		}
	}

	return ""
}

func buildImage(cli *client.Client, sourceLocation string, name string, dockerfile string) error {
	defer clean()

	makeArchive(sourceLocation, dockerfile)

	buildContext, err := os.Open(TMP_FILE)
	defer buildContext.Close()

	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     false,
		Tags:           []string{ name },
		Dockerfile:     "Dockerfile",
	}

	buildResponse, err := cli.ImageBuild(context.Background(), buildContext, options)
	if err != nil {
		return err
	}
	defer buildResponse.Body.Close()

	fmt.Printf("********* %s **********\n", name)
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	return jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, os.Stderr, termFd, isTerm, nil)
}

func makeArchive(directory string, dockerfile string) {
	tar := new(archivex.TarFile)
	tar.Create(TMP_FILE)
	tar.AddAll(directory, false)
	tar.Add("Dockerfile", strings.NewReader(dockerfile), nil)
	tar.Close()
}

func clean() {
	os.Remove(TMP_FILE)
}

func TestRun(service string, conf config.Config) error {
	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		return err
	}

	networkId, err := createTestNetwork(cli)
	if err != nil {
		return err
	}

	for _, service := range conf.Services {
		if service == "mongo" {
			err = createMongoContainer(cli, networkId)
			if err != nil {
				return err
			}
		}
		if service == "redis" {
			err = createRedisContainer(cli, networkId)
			if err != nil {
				return err
			}
		}
	}

	for _, initTask := range conf.Init {
		err = createInitContainer(cli, networkId, conf, initTask)
		if err != nil {
			return err
		}
	}

	return createFunctionContainer(cli, networkId, conf, service)
}

func createTestNetwork(cli *client.Client) (string, error) {
	log.Print("Creating network gnative-test-network")

	options := types.NetworkCreate{
		CheckDuplicate: true,
		Driver: 		"bridge",
	}

	res, err := cli.NetworkCreate(
		context.Background(),
		"gnative-test-network",
		options,
	)

	if err != nil {
		return "", err
	}

	if res.Warning != "" {
		return "", errors.New(res.Warning)
	}

	return res.ID, nil
}

func createRedisContainer(cli *client.Client, networkId string) error {
	return createServiceContainer(
		cli,
		networkId,
		"redis",
		"6379",
		"redis",
	)
}

func createMongoContainer(cli *client.Client, networkId string) error {
	return createServiceContainer(
		cli,
		networkId,
		"mongo",
		"27017",
		"mongo",
	)
}

func createServiceContainer(cli *client.Client, networkId string, image string, port string, alias string) error {
	log.Print("Updating " + image + " image")
	r, err := cli.ImagePull(
		context.Background(),
		"docker.io/library/" + image,
		types.ImagePullOptions{},
	)
	if err != nil {
		return err
	}

	ioutil.ReadAll(r)

	return createContainer(
		cli,
		networkId,
		image,
		port,
		[]string{ alias },
		[]string{},
		"gnative-" + image, 
	)
}

func createInitContainer(cli *client.Client, networkId string, conf config.Config, initTask config.InitTask) error {
	ctx := context.Background()
	env := []string{}

	for _, envVar := range conf.Environment {
		env = append(env, envVar.Name + "=" + envVar.Value)
	}

	for _, envVar := range initTask.Environment {
		env = append(env, envVar.Name + "=" + envVar.Value)
	}

	config := container.Config{
		Image: conf.Registry + "/" + initTask.Image.Name + ":" + initTask.Image.Tag,
		Env: env,
	}

	log.Print("Creating init container " + initTask.Name)

	res, err := cli.ContainerCreate(
		ctx,
		&config,
		nil,
		nil,
		"gnative-init-" + initTask.Name,
	)

	if err != nil {
		return err
	}

	if len(res.Warnings) > 0 {
		return errors.New(strings.Join(res.Warnings, ","))
	}

	containerId := res.ID

	err = cli.ContainerStart(
		ctx,
		containerId,
		types.ContainerStartOptions{},
	)

	if err != nil {
		return err
	}

	return cli.NetworkConnect(
		ctx,
		networkId,
		containerId,
		nil,
	)
}

func createFunctionContainer(cli *client.Client, networkId string, conf config.Config, service string) error {
	env := []string{}

	for _, envVar := range conf.Environment {
		env = append(env, envVar.Name + "=" + envVar.Value)
	}

	return createContainer(
		cli,
		networkId,
		conf.Registry + "/" + service + ":0.0.1",
		"8080",
		[]string{},
		env,
		"gnative-" + service, 
	)
}

func createContainer(cli *client.Client, networkId string, image string, port string, hosts []string, env []string, name string) error {
	ctx := context.Background()

	log.Print("Creating container " + name)
	
	containerPort, err := nat.NewPort("tcp", port)
	if err != nil {
		return err
	}

	config := container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			containerPort: struct{}{},
		},
		Env: env,
	}

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: port,
	}
	
	hostConfig := container.HostConfig{
		PortBindings: nat.PortMap{containerPort: []nat.PortBinding{hostBinding}},
	}

	res, err := cli.ContainerCreate(
		ctx,
		&config,
		&hostConfig,
		nil,
		name,
	)

	if err != nil {
		return err
	}

	if len(res.Warnings) > 0 {
		return errors.New(strings.Join(res.Warnings, ","))
	}

	containerId := res.ID

	err = cli.ContainerStart(
		ctx,
		containerId,
		types.ContainerStartOptions{},
	)

	if err != nil {
		return err
	}

	return cli.NetworkConnect(
		ctx,
		networkId,
		containerId,
		&network.EndpointSettings{
			Aliases: hosts,
		},
	)
}

func TestStop(conf config.Config) error {
	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		return err
	}

	err = stopTestContainers(cli)
	if err != nil {
		return err
	}

	return removeTestNetwork(cli)
}

func stopTestContainers(cli *client.Client) error {
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{ All: true })
	if err != nil {
		return err
	}

	count := 0

	for _, container := range containers {
		if !strings.HasPrefix(container.Names[0], "/gnative-") {
			continue
		}

		count = count + 1

		if (container.State == "running") {
			err = cli.ContainerStop(ctx, container.ID, nil)
			if err != nil {
				return err
			}
		}

		err = cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
		if err != nil {
			return err
		}
	}

	fmt.Printf("- stopped %d container(s)\n", count)

	return nil
}

func removeTestNetwork(cli *client.Client) error {
	ctx := context.Background()
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return err
	}

	count := 0

	for _, network := range networks {
		if !strings.HasPrefix(network.Name, "gnative-") {
			continue
		}

		count = count + 1

		err = cli.NetworkRemove(ctx, network.ID)
		if err != nil {
			return err
		}
	}

	fmt.Printf("- removed %d networks(s)\n", count)

	return nil
}