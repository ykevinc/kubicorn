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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
)

// GetNodeCount returns number of nodes in the cluster.
func GetNodeCount(client *k8s.Clientset) (int, error) {
	nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return -1, err
	}
	return len(nodes.Items), nil
}

// VerifyNodeReadiness returns error in case any node is not ready.
func VerifyNodeReadiness(client *k8s.Clientset) error {
	nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, node := range nodes.Items {
		for _, c := range node.Status.Conditions {
			if c.Type != v1.NodeReady {
				return fmt.Errorf("node %s not ready", node.Name)
			}
		}
	}
	return nil
}

// VerifyComponentStatues returns error if any component is not ready.
func VerifyComponentStatues(client *k8s.Clientset) error {
	compstat, err := client.CoreV1().ComponentStatuses().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, cs := range compstat.Items {
		for _, c := range cs.Conditions {
			if c.Type != v1.ComponentHealthy {
				return fmt.Errorf("component %s not healthy", cs.Name)
			}
		}
	}
	return nil
}
