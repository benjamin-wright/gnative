package routing

import (
	"gnative/config"
)

type Service struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       ServiceSpec `yaml:"spec"`
}

type ServiceSpec struct {
	Template ServiceTemplate `yaml:"template"`
}

type ServiceTemplate struct {
	Spec TemplateSpec `yaml:"spec"`
}

type TemplateSpec struct {
	Containers []Container `yaml:"containers"`
}

type Container struct {
	Image string                       `yaml:"image"`
	Env   []config.EnvironmentVariable `yaml:"env"`
	Ports []ContainerPort              `yaml:"ports"`
}

type EnvironmentVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type ContainerPort struct {
	ContainerPort int `yaml:"containerPort"`
}

func getBaseDeployment(registry string, endpoint config.Endpoint, env []config.EnvironmentVariable) Service {
	return Service{
		ApiVersion: "serving.knative.dev/v1alpha1",
		Kind:       "Service",
		Metadata: Metadata{
			Name:      endpoint.Image.Name,
			Namespace: "default",
		},
		Spec: ServiceSpec{
			Template: ServiceTemplate{
				Spec: TemplateSpec{
					Containers: []Container{
						{
							Image: registry + "/" + endpoint.Image.Name + ":" + endpoint.Image.Tag,
							Ports: []ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							Env: env,
						},
					},
				},
			},
		},
	}
}
