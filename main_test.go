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
	"io/fs"
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

	cmd := exec.Command("pwd")
	stdout, _ := cmd.Output()
	fmt.Printf("pwd: %s\n", string(stdout))
	fmt.Printf("log path: %s\n", filepath.Join("./", artsDir, "/helloworld.txt"))
	os.Mkdir("logs", fs.ModePerm)
	os.WriteFile(filepath.Join("./", "artsDir", "/helloworld.txt"), []byte("Hello, World!"), os.ModePerm)
	os.WriteFile(filepath.Join("./", artsDir, "/helloworld.txt"), []byte("Hello, World!"), os.ModePerm)

	cmd = exec.Command("ls")
	stdout, _ = cmd.Output()
	fmt.Printf("ls: %s\n", string(stdout))

	cmd = exec.Command("ls", "..")
	stdout, _ = cmd.Output()
	fmt.Printf("ls ..: %s\n", string(stdout))

	cmd = exec.Command("ls", "../../..")
	stdout, _ = cmd.Output()
	fmt.Printf("ls ../../..: %s\n", string(stdout))

	cmd = exec.Command("ls", "../../../..")
	stdout, _ = cmd.Output()
	fmt.Printf("ls ../../../..: %s\n", string(stdout))

	var singleClusterKubeConfig = os.Getenv("KUBECONFIG")
	if singleClusterKubeConfig == "" {
		singleClusterKubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	fmt.Printf("singleClusterKubeconfig: %s\n", singleClusterKubeConfig)
	cmd = exec.Command("kubectl", "cluster-info", "dump", "--output-directory=logs1111", "--all-namespaces", fmt.Sprintf("--kubeconfig %v", singleClusterKubeConfig))
	stdout, err := cmd.Output()
	fmt.Printf("cluster dump: %s\n", string(stdout))

	if err != nil {
		fmt.Printf("cluster dump error: %s\n", err.Error())
	}
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
