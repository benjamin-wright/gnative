package routing

import (
	"gnative/config"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func Write(c config.Config) {
	for _, route := range c.Routes {
		routeToVirtualServiceYaml(c.Namespace, route)
		routeToServingYaml(c.Namespace, c.Registry, route, c.Environment)
	}
}

func routeToServingYaml(namespace string, registry string, route config.Route, env []config.EnvironmentVariable) {
	for _, endpoint := range route.Endpoints {
		serving := getBaseServing(namespace, registry, endpoint, env)

		d, err := yaml.Marshal(&serving)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile("manifests/"+route.Name+"_deployment_"+serving.Metadata.Name+".yaml", d, 0644)
		if err != nil {
			panic(err)
		}

		log.Print("Written " + route.Name + " deployment " + serving.Metadata.Name)
	}
}

func routeToVirtualServiceYaml(namespace string, route config.Route) {
	service := getBaseService(namespace, route)

	for _, endpoint := range route.Endpoints {
		service.Spec.Http = append(service.Spec.Http, getHttp(namespace, endpoint))
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
