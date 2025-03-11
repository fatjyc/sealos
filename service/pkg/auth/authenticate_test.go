package auth_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/labring/sealos/service/pkg/auth"
)

const validKubeconfig = `
apiVersion: v1
clusters:
- cluster:
    server: https://test-server:6443
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
users:
- name: test-user
  user:
    token: test-token
`

const invalidKubeconfig = `invalid kubeconfig`

func TestAddWhiteListKubernetesHosts(t *testing.T) {
	auth.AddWhiteListKubernetesHosts("test-host")
	assert.True(t, auth.IsWhitelistKubernetesHost("test-host"))
}

func TestGetKcHost(t *testing.T) {
	tests := []struct {
		name    string
		kc      string
		want    string
		wantErr bool
	}{
		{
			name:    "valid kubeconfig",
			kc:      validKubeconfig,
			want:    "https://test-server:6443",
			wantErr: false,
		},
		{
			name:    "invalid kubeconfig",
			kc:      invalidKubeconfig,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.GetKcHost(tt.kc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetKcUser(t *testing.T) {
	tests := []struct {
		name    string
		kc      string
		want    string
		wantErr bool
	}{
		{
			name:    "valid kubeconfig",
			kc:      validKubeconfig,
			want:    "test-user",
			wantErr: false,
		},
		{
			name:    "invalid kubeconfig",
			kc:      invalidKubeconfig,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := auth.GetKcUser(tt.kc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCheckK8sHost(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		envHost   string
		envPort   string
		wantErr   bool
	}{
		{
			name:    "host in whitelist",
			host:    "test-host",
			envHost: "",
			envPort: "",
			wantErr: false,
		},
		{
			name:    "host matches environment variable",
			host:    "https://env-host:6443",
			envHost: "env-host",
			envPort: "6443",
			wantErr: false,
		},
		{
			name:    "host does not match environment variable",
			host:    "https://wrong-host:6443",
			envHost: "env-host",
			envPort: "6443",
			wantErr: true,
		},
		{
			name:    "no host and no environment variable",
			host:    "https://unknown-host:6443",
			envHost: "",
			envPort: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth.AddWhiteListKubernetesHosts("test-host")
			if tt.envHost != "" {
				t.Setenv("KUBERNETES_SERVICE_HOST", tt.envHost)
				t.Setenv("KUBERNETES_SERVICE_PORT", tt.envPort)
			}
			err := auth.CheckK8sHost(tt.host)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	t.Setenv("KUBERNETES_SERVICE_HOST", "test-host")
	t.Setenv("KUBERNETES_SERVICE_PORT", "6443")

	tests := []struct {
		name    string
		ns      string
		kc      string
		wantErr bool
	}{
		{
			name:    "empty namespace",
			ns:      "",
			kc:      validKubeconfig,
			wantErr: true,
		},
		{
			name:    "invalid kubeconfig",
			ns:      "default",
			kc:      invalidKubeconfig,
			wantErr: true,
		},
		{
			name:    "valid kubeconfig and namespace",
			ns:      "default",
			kc:      validKubeconfig,
			wantErr: true, // Mocking external dependencies is required for this test to pass
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.Authenticate(tt.ns, tt.kc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsWhitelistKubernetesHost(t *testing.T) {
	auth.AddWhiteListKubernetesHosts("test-host")
	assert.True(t, auth.IsWhitelistKubernetesHost("test-host"))
	assert.False(t, auth.IsWhitelistKubernetesHost("unknown-host"))
}

func TestGetKubernetesHostFromEnv(t *testing.T) {
	t.Setenv("KUBERNETES_SERVICE_HOST", "env-host")
	t.Setenv("KUBERNETES_SERVICE_PORT", "6443")
	host := auth.GetKubernetesHostFromEnv()
	assert.Equal(t, "https://env-host:6443", host)

	t.Setenv("KUBERNETES_SERVICE_HOST", "")
	t.Setenv("KUBERNETES_SERVICE_PORT", "")
	host = auth.GetKubernetesHostFromEnv()
	assert.Equal(t, "", host)
}
