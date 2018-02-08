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
	"testing"

	"github.com/kris-nova/kubicorn/cloud"
	"github.com/kris-nova/kubicorn/cutil/agent"
	"github.com/kris-nova/kubicorn/e2e/tutil/healthcheck"
	"github.com/kris-nova/kubicorn/e2e/tutil/k8slogger"
	"github.com/kris-nova/kubicorn/e2e/tutil/kubernetes"
	"github.com/kris-nova/kubicorn/e2e/tutil/scp"
	"github.com/kris-nova/kubicorn/e2e/tutil/ssh"
	"github.com/mholt/archiver"
)

func TestMain(m *testing.M) {
	// Create cluster.
	cluster, reconciler, err := CreateDOUbuntuCluster()
	if err != nil {
		panic(err)
	}

	// Get kubePath.
	agnt := agent.NewAgent()
	kubeFile, err := scp.RetryDownloadFile(cluster, agnt, "/root/.kube/config")
	if err != nil {
		handleError(reconciler, err)
	}

	kubePath, err := scp.CreateTempFileFromBytes(kubeFile, ".kube", "config")
	if err != nil {
		handleError(reconciler, err)
	}

	// New Kubernetes Client.
	client, err := kubernetes.NewClient(kubePath)
	if err != nil {
		handleError(reconciler, err)
	}

	// Node readiness.
	_, err = healthcheck.RetryVerifyNodeReadiness(client)
	if err != nil {
		handleError(reconciler, err)
	}
	/*if count != 3 {
		handleError(reconciler, fmt.Errorf("node count missmatch"))
	}*/

	// Make sure componenets are ready.
	err = healthcheck.VerifyComponentStatuses(client)
	if err != nil {
		handleError(reconciler, err)
	}

	// Create Sonobuoy stuff.
	err = ssh.ExecCommandSSH(cluster, agent.NewAgent(),
		"kubectl apply -f https://raw.githubusercontent.com/heptio/sonobuoy/master/examples/quickstart.yaml")
	if err != nil {
		handleError(reconciler, err)
	}

	err = k8slogger.WaitPodLogsStream(client, "sonobuoy", "heptio-sonobuoy")
	if err != nil {
		handleError(reconciler, err)
	}

	err = ssh.ExecCommandSSH(cluster, agent.NewAgent(),
		"kubectl cp heptio-sonobuoy/sonobuoy:/tmp/sonobuoy /root/archive --namespace=heptio-sonobuoy && mv /root/archive/*sonobuoy* /root/archive/sonobuoy.tar.gz")
	if err != nil {
		handleError(reconciler, err)
	}

	sb, err := scp.RetryDownloadFile(cluster, agent.NewAgent(), "/root/archive/sonobuoy.tar.gz")
	if err != nil {
		handleError(reconciler, err)
	}

	sbPath, err := scp.CreateTempFileFromBytes(sb, "archive", "sonobuoy.tar.gz")
	if err != nil {
		handleError(reconciler, err)
	}

	err = archiver.Zip.Open(sbPath, "./")
	if err != nil {
		handleError(reconciler, err)
	}

	// Remove cluster.
	err = DestroyDOUbuntuCluster(reconciler)
	if err != nil {
		panic(err)
	}
}

func handleError(reconciler cloud.Reconciler, err error) {
	// Remove cluster.
	e := DestroyDOUbuntuCluster(reconciler)
	if err != nil {
		panic(e)
	}
	panic(err)
}
