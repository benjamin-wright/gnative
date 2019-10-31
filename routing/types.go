package routing

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
	Name  string                       `yaml:"name"`
	Image string                       `yaml:"image"`
	Env   []EnvironmentVariable `yaml:"env"`
	Ports []ContainerPort              `yaml:"ports"`
}

type EnvironmentVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type ContainerPort struct {
	ContainerPort int `yaml:"containerPort"`
}

type VirtualService struct {
	ApiVersion string             `yaml:"apiVersion"`
	Kind       string             `yaml:"kind"`
	Metadata   Metadata           `yaml:"metadata"`
	Spec       VirtualServiceSpec `yaml:"spec"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type VirtualServiceSpec struct {
	Gateways []string `yaml:"gateways"`
	Hosts    []string `yaml:"hosts"`
	Http     []Http   `yaml:"http"`
}

type Http struct {
	Match   []MatchRule `yaml:"match"`
	Rewrite RewriteRule `yaml:"rewrite"`
	Route   []Route     `yaml:"route"`
}

type MatchRule struct {
	Uri Uri `yaml:"uri"`
}

type Uri struct {
	Prefix string `yaml:"prefix"`
}

type RewriteRule struct {
	Authority string `yaml:"authority"`
}

type Route struct {
	Destination Destination `yaml:"destination"`
}

type Destination struct {
	Host   string `yaml:"host"`
	Port   Port   `yaml:"port"`
	Weight int    `yaml:"weight"`
}

type Port struct {
	Number int `yaml:"number"`
}