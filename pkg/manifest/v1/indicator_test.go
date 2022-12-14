//nolint:gci
/*
Copyright © 2022 OpenSLO Team
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

//nolint:godot
// spell-checker:disable
func Test_MetricSource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		yaml string
		want MetricSource
	}{
		{
			name: "CloudWatch - One Dimension",
			yaml: `metricSourceRef: usix-cloudwatch
type: CloudWatch
spec:
  metricName: 2xx
  namespace: CloudWatchSynthetics
  region: us-east-1
  stat: SampleCount
  dimensions:
    - name: following
      value: bar`,
			want: MetricSource{
				MetricSourceRef: "usix-cloudwatch",
				Type:            "CloudWatch",
				MetricSourceSpec: map[string]string{
					"metricName": "2xx",
					"namespace":  "CloudWatchSynthetics",
					"region":     "us-east-1",
					"stat":       "SampleCount",
					"dimensions": "name:following,value:bar",
				},
			},
		},
		{
			name: "CloudWatch - Two Dimensions",
			yaml: `metricSourceRef: usix-cloudwatch
type: CloudWatch
spec:
  metricName: 2xx
  namespace: CloudWatchSynthetics
  region: us-east-1
  stat: SampleCount
  dimensions:
    - name: following
      value: bar
    - name: another
      value: batz`,
			want: MetricSource{
				MetricSourceRef: "usix-cloudwatch",
				Type:            "CloudWatch",
				MetricSourceSpec: map[string]string{
					"metricName": "2xx",
					"namespace":  "CloudWatchSynthetics",
					"region":     "us-east-1",
					"stat":       "SampleCount",
					"dimensions": "name:following,value:bar;name:another,value:batz",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var out MetricSource
			if err := yaml.Unmarshal([]byte(tt.yaml), &out); err != nil {
				t.Fatalf("Failed to unmarshal the yaml: %+v", err)
			}
			assert.Equal(t, tt.want, out)
		})
	}
}

// spell-checker:enable
