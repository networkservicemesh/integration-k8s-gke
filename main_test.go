// Copyright (c) 2021-2022 Doc.ai and/or its affiliates.
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
	"testing"

	"github.com/stretchr/testify/suite"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // This is required for GKE authentication

	"github.com/networkservicemesh/integration-tests/suites/basic"
	"github.com/networkservicemesh/integration-tests/suites/features"
	"github.com/networkservicemesh/integration-tests/suites/heal"
	"github.com/networkservicemesh/integration-tests/suites/memory"
	"github.com/networkservicemesh/integration-tests/suites/observability"
)

func TestRunBasicSuite(t *testing.T) {
	suite.Run(t, new(basic.Suite))
}

func TestRunMemorySuite(t *testing.T) {
	suite.Run(t, new(memory.Suite))
}

func TestRunObservabilitySuite(t *testing.T) {
	suite.Run(t, new(observability.Suite))
}

// Disabled tests:
// TestVl3_nscs_death - https://github.com/networkservicemesh/integration-k8s-gke/issues/327
// TestVl3_nse_death  - https://github.com/networkservicemesh/integration-k8s-gke/issues/327
type healSuite struct {
	heal.Suite
}

func (s *healSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case
		"TestVl3_nscs_death",
		"TestVl3_nse_death":
		s.T().Skip()
	}
	s.Suite.BeforeTest(suiteName, testName)
}

func TestRunHealSuite(t *testing.T) {
	suite.Run(t, new(healSuite))
}

// Disabled tests:
// TestVl3_basic           - https://github.com/networkservicemesh/integration-k8s-gke/issues/327
// TestVl3_scale_from_zero - https://github.com/networkservicemesh/integration-k8s-gke/issues/327
type featuresSuite struct {
	features.Suite
}

func (s *featuresSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case
		"TestVl3_basic",
		"TestVl3_scale_from_zero":
		s.T().Skip()
	}
	s.Suite.BeforeTest(suiteName, testName)
}

func TestRunFeatureSuiteCalico(t *testing.T) {
	suite.Run(t, new(featuresSuite))
}
