package routing

import (
	"gnative/config"
)

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
