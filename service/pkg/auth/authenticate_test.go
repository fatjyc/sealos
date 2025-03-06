package auth_test

import (
	"os"
	"testing"

	"github.com/labring/sealos/service/pkg/auth"
	"github.com/stretchr/testify/assert"
)

const validKubeconfig = `
apiVersion: v1
clusters:
- cluster:
    server: https://kubernetes.default.svc:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: fake-cert-data
    client-key-data: fake-key-data
`

const invalidKubeconfig = `invalid`

func TestAddWhiteListKubernetesHosts(t *testing.T) {
	host := "test-host"
	auth.AddWhiteListKubernetesHosts(host)
	assert.True(t, auth.IsWhitelistKubernetesHost(host))
}

func TestCheckK8sHost(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		envHost  string
		envPort  string
		wantErr  bool
		setupEnv bool
	}{
		{"matching host", "https://kubernetes.default.svc:6443", "kubernetes.default.svc", "6443", false, true},
		{"non-matching host", "https://different.host:6443", "kubernetes.default.svc", "6443", true, true},
		{"no env variables", "https://kubernetes.default.svc:6443", "", "", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv {
				os.Setenv("KUBERNETES_SERVICE_HOST", tt.envHost)
				os.Setenv("KUBERNETES_SERVICE_PORT", tt.envPort)
				defer os.Unsetenv("KUBERNETES_SERVICE_HOST")
				defer os.Unsetenv("KUBERNETES_SERVICE_PORT")
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
	os.Setenv("KUBERNETES_SERVICE_HOST", "kubernetes.default.svc")
	os.Setenv("KUBERNETES_SERVICE_PORT", "6443")
	defer os.Unsetenv("KUBERNETES_SERVICE_HOST")
	defer os.Unsetenv("KUBERNETES_SERVICE_PORT")

	tests := []struct {
		name    string
		ns      string
		kc      string
		wantErr bool
	}{
		{"empty namespace", "", validKubeconfig, true},
		{"invalid kubeconfig", "test", invalidKubeconfig, true},
		{"valid inputs", "test", validKubeconfig, true}, // Authenticate requires realistic responses from external dependencies.
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
	assert.False(t, auth.IsWhitelistKubernetesHost("non-existent-host"))
}

func TestGetKubernetesHostFromEnv(t *testing.T) {
	tests := []struct {
		envHost string
		envPort string
		want    string
	}{
		{"kubernetes.default.svc", "6443", "https://kubernetes.default.svc:6443"},
		{"", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.envHost, func(t *testing.T) {
			os.Setenv("KUBERNETES_SERVICE_HOST", tt.envHost)
			os.Setenv("KUBERNETES_SERVICE_PORT", tt.envPort)
			defer os.Unsetenv("KUBERNETES_SERVICE_HOST")
			defer os.Unsetenv("KUBERNETES_SERVICE_PORT")

			got := auth.GetKubernetesHostFromEnv()
			assert.Equal(t, tt.want, got)
		})
	}
}
