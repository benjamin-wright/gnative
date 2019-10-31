package tasks

type Job struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       JobSpec 	   `yaml:"spec"`
}

type Metadata struct {
	Name      string 			`yaml:"name"`
	Namespace string 			`yaml:"namespace"`
}

type JobSpec struct {
	BackoffLimit int 		 `yaml:"backoffLimit"`
	Template     JobTemplate `yaml:"template"`
}

type JobTemplate struct {
	Spec 	 TemplateSpec 		`yaml:"spec"`
}

type TemplateSpec struct {
	Containers    []Container `yaml:"containers"`
	RestartPolicy string      `yaml:"restartPolicy"`
}

type Container struct {
	Name  string                `yaml:"name"`
	Image string                `yaml:"image"`
	Env   []EnvironmentVariable `yaml:"env"`
}

type EnvironmentVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}