package services

import (
	"io/ioutil"
	"log"
)

func Write(serviceNames []string) {
	for _, service := range serviceNames {
		if service == "mongo" {
			writeMongoOutputs()
		}
	}
}

func writeMongoOutputs() {
	for _, output := range getMongoOutputs() {
		ioutil.WriteFile("manifests/"+output.file, []byte(output.data), 0644)
		log.Print("Written manifest " + output.file)
	}
}
