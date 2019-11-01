package main

import (
	"gnative/config"
	"gnative/routing"
	"gnative/services"
	"gnative/tasks"
	"gnative/docker"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:   "template",
			Usage:  "build kubernetes configuration files",
			Action: template,
		},
		{
			Name:   "hosts",
			Usage:  "setup hosts file to redirect hosts to 127.0.0.1",
			Action: hosts,
		},
		{
			Name:   "build",
			Usage:  "build the function images",
			Action: build,
		},
		{
			Name:	"test-run",
			Usage:  "run one of the functions locally",
			Action: testRun,
		},
		{
			Name:	"test-stop",
			Usage:  "stop local containers",
			Action: testStop,
		},
	}

	err := app.Run(os.Args)
	check(err)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func template(c *cli.Context) error {
	conf, err := config.Get()
	check(err)

	os.Mkdir("manifests", 0777)

	routing.Write(conf)
	services.Write(conf.Namespace, conf.Services)
	tasks.Write(conf)

	return nil
}

func hosts(c *cli.Context) error {
	conf, err := config.Get()
	check(err)

	err = config.SetHosts(conf)

	if err != nil {
		if os.IsPermission(err) {
			log.Print("Permissions error, try rerunning with sudo")
		} else {
			log.Fatal(err)
		}
	}

	return nil
}

func build(c *cli.Context) error {
	conf, err := config.Get()
	check(err)

	err = docker.Build(conf)
	check(err)

	return nil
}

func testRun(c *cli.Context) error {
	conf, err := config.Get()
	check(err)

	function := c.Args().Get(0)
	if function == "" {
		log.Fatal("test-run requires at least one argument")
	}

	log.Print("Stopping previous test stuff...")

	err = docker.TestStop(conf)
	check(err)

	log.Print("Running '" + c.Args().Get(0) + "' locally...")

	err = docker.TestRun(function, conf)
	check(err)

	return nil
}

func testStop(c *cli.Context) error {
	conf, err := config.Get()
	check(err)

	log.Print("Cleaning up local test containers")

	err = docker.TestStop(conf)
	check(err)

	return nil
}