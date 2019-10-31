package tasks


import (
	"gnative/config"
	"io/ioutil"
	"log"
	
	"gopkg.in/yaml.v2"
)

func Write(c config.Config) {
	for _, task := range c.Init {
		job := taskToJob(c.Registry, task, c.Environment)
		d, err := yaml.Marshal(&job)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile("manifests/init_job_"+task.Name+".yaml", d, 0644)
		if err != nil {
			panic(err)
		}

		log.Print("Written " + task.Name + " init task job")
	}
}