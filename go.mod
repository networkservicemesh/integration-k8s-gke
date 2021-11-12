module github.com/networkservicemesh/integration-k8s-gke

go 1.16

require (
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/networkservicemesh/integration-tests v0.0.0-20211104152816-6a9063fda6d0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/client-go v0.20.5
)

replace github.com/networkservicemesh/integration-tests v0.0.0-20211104152816-6a9063fda6d0 => github.com/glazychev-art/integration-tests v0.0.0-20211115061303-b291dcbd3eff
