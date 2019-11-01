package tasks

import (
	"gnative/config"
)

func taskToJob(namespace string, registry string, task config.InitTask, globalEnv []config.EnvironmentVariable) Job {
	envVars := []EnvironmentVariable{}

	for _, env := range task.Environment {
		envVars = append(envVars, EnvironmentVariable{
			Name: env.Name,
			Value: env.Value,
		})
	}

	for _, env := range globalEnv {
		envVars = append(envVars, EnvironmentVariable{
			Name: env.Name,
			Value: env.Value,
		})
	}

	return Job{
		ApiVersion: "batch/v1",
		Kind: "Job",
		Metadata: Metadata{
			Name: task.Name,
			Namespace: namespace,
		},
		Spec: JobSpec{
			Template: JobTemplate{
				Spec: TemplateSpec{
					Containers: []Container{
						{
							Name: task.Name,
							Image: registry + "/" + task.Image.Name + ":" + task.Image.Tag,
							Env: envVars,
						},
					},
					RestartPolicy: "Never",
				},
			},
			BackoffLimit: 4,
		},
	}
}