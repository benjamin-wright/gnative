package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Source struct {
	Libraries string
	Images string
}

type Image struct {
	Name string
	Tag  string
}

type Endpoint struct {
	Path  string
	Image Image
}

type Route struct {
	Name      string
	Hostname  string
	Endpoints []Endpoint
}

type EnvironmentVariable struct {
	Name  string
	Value string
}

type Config struct {
	Registry    string
	Source 			Source
	Services    []string
	Routes      []Route
	Environment []EnvironmentVariable
}

func assertExists() {
	_, err := os.Stat("gonative.yaml")
	if os.IsNotExist(err) {
		log.Fatal("Config not found, are you in the right directory?")
	}
}

func Get() (Config, error) {
	assertExists()

	config := Config{}

	configData, err := ioutil.ReadFile("gonative.yaml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

var HOSTS_START string = "# Added by gnative"
var HOSTS_END string = "# End of section"

func SetHosts(config Config) error {
	assertExists()

	hostsFile, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		return err
	}

	lines := strings.Split(string(hostsFile), "\n")
	startIndex, endIndex := findStartEnd(lines, HOSTS_START, HOSTS_END)

	outputs := []string{}

	for i := 0; i < startIndex; i++ {
		outputs = append(outputs, lines[i])
	}

	for i := endIndex + 1; i < len(lines); i++ {
		outputs = append(outputs, lines[i])
	}

	outputs = append(outputs, HOSTS_START)
	for _, route := range config.Routes {
		outputs = append(outputs, "127.0.0.1      "+route.Hostname)
	}
	outputs = append(outputs, HOSTS_END)

	outputData := []byte(strings.Join(outputs, "\n"))

	err = ioutil.WriteFile("/etc/hosts", outputData, 0644)
	return err
}

func findStartEnd(lines []string, start string, end string) (int, int) {
	startIndex := -1
	endIndex := -1

	for i, l := range lines {
		if l == start {
			startIndex = i
			continue
		}

		if l == end && startIndex >= 0 {
			endIndex = i
			break
		}
	}

	return startIndex, endIndex
}
