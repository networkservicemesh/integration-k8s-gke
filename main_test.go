// Copyright (c) 2021 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main_test

import (
	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/spire"
	"testing"

	"github.com/stretchr/testify/suite"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // This is required for GKE authentication
)

type Suite struct {
	base.Suite
	spireSuite spire.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.spireSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/basic")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete mutatingwebhookconfiguration --all` + "\n" + `kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/registry-k8s-sa \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:registry-k8s-sa`)
	r.Run(`kubectl apply -k .`)
}

func (s *Suite) TestMemif2Vxlan2Memif() {
	r := s.Runner("../deployments-k8s/examples/use-cases/Memif2Vxlan2Memif")
	s.T().Cleanup(func() {
		r.Run(`kubectl describe pod "${FWD0}" -n nsm-system`)
		r.Run(`kubectl describe pod "${FWD1}" -n nsm-system`)
		r.Run(`kubectl exec "${FWD0}" -n nsm-system -- vppctl show int addr`)
		r.Run(`kubectl exec "${FWD0}" -n nsm-system -- vppctl show fib entry`)
		r.Run(`kubectl exec "${FWD0}" -n nsm-system -- vppctl show trace max 10`)
		r.Run(`kubectl exec "${FWD1}" -n nsm-system -- vppctl show trace max 20`)
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODES=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}'))`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/nsc-memif` + "\n" + `- ../../../apps/nse-memif` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc-memif` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: memif://icmp-responder/nsm-1` + "\n" + `` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[0]}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse-memif` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nse` + "\n" + `          env:` + "\n" + `            - name: NSE_CIDR_PREFIX` + "\n" + `              value: 172.16.1.100/31` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[1]}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc-memif -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse-memif -n ${NAMESPACE}`)
	r.Run(`NSC=$(kubectl get pods -l app=nsc-memif -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`NSE=$(kubectl get pods -l app=nse-memif -n ${NAMESPACE} --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl get pods --all-namespaces`)
	r.Run(`kubectl describe nodes`)
	r.Run(`FWD0=$(kubectl get pods -n nsm-system -l app=forwarder-vpp --field-selector spec.nodeName=${NODES[0]} --template='{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`FWD1=$(kubectl get pods -n nsm-system -l app=forwarder-vpp --field-selector spec.nodeName=${NODES[1]} --template='{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod ${FWD0} -n nsm-system`)
	r.Run(`kubectl exec "${FWD0}" -n nsm-system -- vppctl trace add memif-input 10`)
	r.Run(`kubectl exec "${FWD1}" -n nsm-system -- vppctl trace add af-packet-input 20`)
	r.Run(`sleep 5`)
	r.Run(`kubectl exec "${FWD0}" -n nsm-system -- ip neigh show`)
	r.Run(`kubectl exec "${FWD0}" -n nsm-system -- ip a`)
	r.Run(`kubectl exec "${FWD0}" -n nsm-system -- ip route`)
	r.Run(`result=$(kubectl exec "${NSC}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.100 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
	r.Run(`result=$(kubectl exec "${NSE}" -n "${NAMESPACE}" -- vppctl ping 172.16.1.101 repeat 4)` + "\n" + `echo ${result}` + "\n" + `! echo ${result} | grep -E -q "(100% packet loss)|(0 sent)|(no egress interface)"`)
}

func TestRunBasicSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
