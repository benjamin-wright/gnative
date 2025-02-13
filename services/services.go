package services

import (
	"io/ioutil"
	"log"
)

func Write(namespace string, serviceNames []string) {
	for _, service := range serviceNames {
		if service == "mongo" {
			writeMongoOutputs(namespace)
		}
		
		if service == "redis" {
			writeRedisOutputs(namespace)
		}
	}
}

func writeMongoOutputs(namespace string) {
	for _, output := range getMongoOutputs(namespace) {
		ioutil.WriteFile("manifests/"+output.file, []byte(output.data), 0644)
		log.Print("Written manifest " + output.file)
	}
}

func writeRedisOutputs(namespace string) {
	for _, output := range getRedisOutputs(namespace) {
		ioutil.WriteFile("manifests/"+output.file, []byte(output.data), 0644)
		log.Print("Written manifest " + output.file)
	}
}

