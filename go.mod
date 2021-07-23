module github.com/networkservicemesh/integration-k8s-gke

go 1.16

require (
	github.com/networkservicemesh/integration-tests v0.0.0-20210719144646-9523da794998
	github.com/stretchr/testify v1.7.0
	k8s.io/client-go v0.20.5
)

replace github.com/networkservicemesh/integration-tests => github.com/Mixaster995/integration-tests v0.0.0-20210723081811-2c27d8ae4b06
