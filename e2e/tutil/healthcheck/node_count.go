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

package healthcheck

import (
	"fmt"
	"time"

	"github.com/kris-nova/kubicorn/cutil/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
)

// RetryGetNodeCount waits for expected number of nodes to come up.
func RetryGetNodeCount(client *k8s.Clientset, expectedNodes int) (int, error) {
	for i := 0; i <= retryAttempts; i++ {
		cnt, err := getNodeCount(client, expectedNodes)
		if err != nil || cnt != expectedNodes {
			logger.Debug("Waiting for Nodes to be created..")
			time.Sleep(time.Duration(retrySleepSeconds) * time.Second)
			continue
		}
		return cnt, nil
	}
	return -1, fmt.Errorf("Timedout waiting nodes to be created.")
}

// getNodeCount returns number of nodes in the cluster.
func getNodeCount(client *k8s.Clientset, expectedNodes int) (int, error) {
	nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return -1, err
	}
	return len(nodes.Items), nil
}
