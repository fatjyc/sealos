// Copyright © 2025 sealos.
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

package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/labring/sealos/controllers/job/heartbeat/internal/util"
)

func TestNewKubernetesClient(t *testing.T) {
	originalGetConfigOrDie := ctrl.GetConfigOrDie

	defer func() {
		ctrl.GetConfigOrDie = originalGetConfigOrDie
	}()

	tests := []struct {
		name        string
		setupConfig func()
		wantErr     bool
	}{
		{
			name: "successfully create client with valid config",
			setupConfig: func() {
				ctrl.GetConfigOrDie = func() *rest.Config {
					return &rest.Config{
						Host: "https://localhost:6443",
					}
				}
			},
			wantErr: false,
		},
		{
			name: "successfully create client with minimal config",
			setupConfig: func() {
				ctrl.GetConfigOrDie = func() *rest.Config {
					return &rest.Config{
						Host: "http://localhost:8080",
					}
				}
			},
			wantErr: false,
		},
		{
			name: "create client with complete config",
			setupConfig: func() {
				ctrl.GetConfigOrDie = func() *rest.Config {
					return &rest.Config{
						Host:          "https://kubernetes.default.svc",
						BearerToken:   "test-token",
						APIPath:       "/api",
						ContentConfig: rest.ContentConfig{ContentType: "application/json"},
					}
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupConfig != nil {
				tt.setupConfig()
			}

			client, err := util.NewKubernetesClient()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}
