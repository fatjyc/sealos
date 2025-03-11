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

package label

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecommended_Labels(t *testing.T) {
	tests := []struct {
		name string
		r    *Recommended
		want map[string]string
	}{
		{
			name: "empty recommended",
			r:    &Recommended{},
			want: map[string]string{},
		},
		{
			name: "all fields set",
			r: &Recommended{
				Name:      "test-app",
				Instance:  "test-instance",
				Version:   "v1.0.0",
				Component: "web",
				PartOf:    "test-system",
				ManagedBy: "sealos",
			},
			want: map[string]string{
				AppName:      "test-app",
				AppInstance:  "test-instance",
				AppVersion:   "v1.0.0",
				AppComponent: "web",
				AppPartOf:    "test-system",
				AppManagedBy: "sealos",
			},
		},
		{
			name: "partial fields set",
			r: &Recommended{
				Name:      "test-app",
				Version:   "v1.0.0",
				ManagedBy: "sealos",
			},
			want: map[string]string{
				AppName:      "test-app",
				AppVersion:   "v1.0.0",
				AppManagedBy: "sealos",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.Labels()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRecommendedLabels(t *testing.T) {
	r := &Recommended{
		Name:      "test-app",
		Instance:  "test-instance",
		Version:   "v1.0.0",
		Component: "web",
		PartOf:    "test-system",
		ManagedBy: "sealos",
	}

	want := map[string]string{
		AppName:      "test-app",
		AppInstance:  "test-instance",
		AppVersion:   "v1.0.0",
		AppComponent: "web",
		AppPartOf:    "test-system",
		AppManagedBy: "sealos",
	}

	got := RecommendedLabels(r)
	assert.Equal(t, want, got)
}
