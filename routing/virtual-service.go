package routing

import (
	"gnative/config"
)

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

func getBaseService(route config.Route) VirtualService {
	return VirtualService{
		ApiVersion: "networking.istio.io/v1alpha3",
		Kind:       "VirtualService",
		Metadata: Metadata{
			Name:      route.Name,
			Namespace: "default",
		},
		Spec: VirtualServiceSpec{
			Gateways: []string{"knative-ingress-gateway.knative-serving.svc.cluster.local"},
			Hosts:    []string{route.Hostname},
		},
	}
}

func getHttp(endpoint config.Endpoint) Http {
	return Http{
		Match: []MatchRule{
			{
				Uri: Uri{
					Prefix: endpoint.Path,
				},
			},
		},
		Rewrite: RewriteRule{
			Authority: endpoint.Image.Name + ".default.example.com",
		},
		Route: []Route{
			{
				Destination: Destination{
					Host: "istio-ingressgateway.istio-system.svc.cluster.local",
					Port: Port{
						Number: 80,
					},
					Weight: 100,
				},
			},
		},
	}
}
