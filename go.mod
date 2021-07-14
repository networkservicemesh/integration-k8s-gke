module github.com/networkservicemesh/integration-k8s-gke

go 1.16

require (
	github.com/networkservicemesh/integration-tests v0.0.0-20210713144223-23171c9ef0e4
	github.com/stretchr/testify v1.7.0
	k8s.io/client-go v0.20.5
)

replace github.com/networkservicemesh/integration-tests => github.com/Mixaster995/integration-tests v0.0.0-20210714072306-6f54af8caf97
