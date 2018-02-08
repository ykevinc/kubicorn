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

package scp

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kris-nova/kubicorn/apis/cluster"
	"github.com/kris-nova/kubicorn/cutil/agent"
	"github.com/kris-nova/kubicorn/cutil/local"
	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	// RetryAttempts specifies the amount of retries are allowed when getting a file from a server.
	RetryAttempts = 150
	// RetrySleepSeconds specifies the time to sleep after a failed attempt to get a file form a server.
	RetrySleepSeconds = 5
)

// DownloadFile returns a file from the cluster.
func DownloadFile(existing *cluster.Cluster, sshAgent *agent.Keyring, remotePath string) ([]byte, error) {
	user := existing.SSH.User
	address := fmt.Sprintf("%s:%s", existing.KubernetesAPI.Endpoint, existing.SSH.Port)
	pubKeyPath := local.Expand(existing.SSH.PublicKeyPath)
	if existing.SSH.Port == "" {
		existing.SSH.Port = "22"
	}

	sshConfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Check for key
	if err := sshAgent.CheckKey(pubKeyPath); err != nil {
		if keyring, err := sshAgent.AddKey(pubKeyPath); err != nil {
			return nil, err
		} else {
			sshAgent = keyring
		}
	}

	if sshAgent != nil && os.Getenv("KUBICORN_FORCE_DISABLE_SSH_AGENT") == "" {
		sshConfig.Auth = append(sshConfig.Auth, sshAgent.GetAgent())
	}

	sshConfig.SetDefaults()
	conn, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	r, err := c.Open(remotePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

// RetryDownloadFile trys to get file until timeout doesn't occurs.
func RetryDownloadFile(existing *cluster.Cluster, sshAgent *agent.Keyring, remotePath string) ([]byte, error) {
	for i := 0; i <= RetryAttempts; i++ {
		file, err := DownloadFile(existing, sshAgent, remotePath)
		if err != nil {
			logger.Debug("Waiting for Kubernetes to come up.. [%v]", err)
			time.Sleep(time.Duration(RetrySleepSeconds) * time.Second)
			continue
		}
		return file, nil
	}
	return nil, fmt.Errorf("Timedout downloading file")
}
