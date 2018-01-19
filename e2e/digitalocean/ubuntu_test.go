// Copyright © 2017 The Kubicorn Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package digitalocean

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kris-nova/kubicorn/cutil/agent"
	"github.com/kris-nova/kubicorn/e2e/tutil/healthcheck"
	"github.com/kris-nova/kubicorn/e2e/tutil/k8slogger"
	"github.com/kris-nova/kubicorn/e2e/tutil/kubeconfig"
	"github.com/kris-nova/kubicorn/e2e/tutil/kubernetes"
	"github.com/kris-nova/kubicorn/e2e/tutil/sshcmd"
)

func TestMain(m *testing.M) {
	// Create cluster.
	cluster, reconciler, err := CreateDOUbuntuCluster()
	if err != nil {
		panic(err)
	}

	// Get kubePath.
	ssh := agent.NewAgent()
	kubePath, err := kubeconfig.RetryGetConfigFilePath(cluster, ssh)
	if err != nil {
		panic(err)
	}

	// New Kubernetes Client.
	client, err := kubernetes.NewClient(kubePath)
	if err != nil {
		panic(err)
	}

	// Node count.
	c, err := healthcheck.GetNodeCount(client)
	if err != nil {
		panic(err)
	}

	fmt.Println(c)

	// Node readiness.
	err = healthcheck.VerifyNodeReadiness(client)
	if err != nil {
		panic(err)
	}

	// VerifyComponentStatues.
	err = healthcheck.VerifyComponentStatues(client)
	if err != nil {
		panic(err)
	}

	err = sshcmd.ExecCommandSSH(cluster, agent.NewAgent(),
		"kubectl apply -f https://raw.githubusercontent.com/heptio/sonobuoy/master/examples/quickstart.yaml")
	if err != nil {
		panic(err)
	}

	rc, err := k8slogger.GetPodLogsStream(client, "sonobuoy", "heptio-sonobuoy")
	if err != nil {
		panic(err)
	}

	b := new(bytes.Buffer)
	b.ReadFrom(rc)
	fmt.Println(b.String())

	// Remove cluster.
	err = DestroyDOUbuntuCluster(reconciler)
	if err != nil {
		panic(err)
	}
}
