// Copyright (c) 2021-2022 Doc.ai and/or its affiliates.
//
// Copyright (c) 2023 Cisco and/or its affiliates.
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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/networkservicemesh/integration-tests/suites/memory"
	"github.com/stretchr/testify/suite"
)

func TestExample(t *testing.T) {
	artsDir := os.Getenv("ARTIFACTS_DIR")
	if artsDir == "" {
		artsDir = "logs"
	}

	mem := new(memory.Suite)
	mem.SetT(t)
	r := mem.Runner(".")

	cmd := exec.Command("pwd")
	stdout, _ := cmd.Output()
	fmt.Printf("pwd: %s\n", string(stdout))

	cmd = exec.Command("ls")
	stdout, _ = cmd.Output()
	fmt.Printf("ls: %s\n", string(stdout))

	stdout, err := exec.Command("kubectl", "config", "view").Output()
	fmt.Printf("kubectl config view: %s err: %v\n", string(stdout), err)

	r.Run("kubectl config view")

	var singleClusterKubeConfig = os.Getenv("KUBECONFIG")
	if singleClusterKubeConfig == "" {
		singleClusterKubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	r.Run(fmt.Sprintf("cat %s", singleClusterKubeConfig))
	fmt.Printf("singleClusterKubeconfig: %s\n", singleClusterKubeConfig)

	r.Run(fmt.Sprintf("kubectl cluster-info --kubeconfig \"%s\" dump --output-directory=logs --all-namespaces --v=9", singleClusterKubeConfig))
}

type calicoFeatureSuite struct {
	memory.Suite
}

func (s *calicoFeatureSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case
		"TestKernel2kernel",
		"TestKernel2ethernet2kernel":
		s.T().Skip()
	}
}

func TestRunMemorySuite(t *testing.T) {
	cmd := exec.Command("pwd")
	stdout, _ := cmd.Output()
	fmt.Printf("pwd: %s\n", string(stdout))

	suite.Run(t, new(calicoFeatureSuite))

	cmd = exec.Command("ls")
	stdout, _ = cmd.Output()
	fmt.Printf("ls: %s\n", string(stdout))
}
