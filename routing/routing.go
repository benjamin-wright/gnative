package routing

import (
	"gnative/config"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func Write(c config.Config) {
	for _, route := range c.Routes {
		routeToServiceYaml(route)
		routeToDeploymentsYaml(c.Registry, route, c.Environment)
	}
}

func routeToDeploymentsYaml(registry string, route config.Route, env []config.EnvironmentVariable) {
	for _, endpoint := range route.Endpoints {
		deployment := getBaseDeployment(registry, endpoint, env)

		d, err := yaml.Marshal(&deployment)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile("manifests/"+route.Name+"_deployment_"+deployment.Metadata.Name+".yaml", d, 0644)
		if err != nil {
			panic(err)
		}

		log.Print("Written " + route.Name + " deployment " + deployment.Metadata.Name)
	}
}

func routeToServiceYaml(route config.Route) {
	service := getBaseService(route)

	for _, endpoint := range route.Endpoints {
		service.Spec.Http = append(service.Spec.Http, getHttp(endpoint))
	}

	d, err := yaml.Marshal(&service)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("manifests/"+route.Name+"_service_"+service.Metadata.Name+".yaml", d, 0644)
	if err != nil {
		panic(err)
	}

	log.Print("Written " + route.Name + " service " + service.Metadata.Name)
}
