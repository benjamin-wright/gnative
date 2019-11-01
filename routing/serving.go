package routing

import (
	"gnative/config"
)

func getBaseServing(namespace string, registry string, endpoint config.Endpoint, env []config.EnvironmentVariable) Service {
	envVars := []EnvironmentVariable{}

	for _, e := range env {
		envVars = append(envVars, EnvironmentVariable{
			Name: e.Name,
			Value: e.Value,
		})
	}

	return Service{
		ApiVersion: "serving.knative.dev/v1alpha1",
		Kind:       "Service",
		Metadata: Metadata{
			Name:      endpoint.Image.Name,
			Namespace: namespace,
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
							Env: envVars,
						},
					},
				},
			},
		},
	}
}
