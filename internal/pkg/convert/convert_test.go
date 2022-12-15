/*
Package convert provides a command to convert from openslo to other formats.

# Copyright © 2022 OpenSLO Team

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
package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	v1 "github.com/OpenSLO/oslo/pkg/manifest/v1"
)

func Test_getCountMetrics(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		ratioMetric v1.RatioMetric
		want        string
		wantErr     bool
	}{
		{
			name: "Fail with Bad",
			ratioMetric: v1.RatioMetric{
				Bad:   &v1.MetricSourceHolder{},
				Total: v1.MetricSourceHolder{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "RatioMetric",
			ratioMetric: v1.RatioMetric{
				Good: &v1.MetricSourceHolder{
					MetricSource: v1.MetricSource{
						Type: "Prometheus",
						MetricSourceSpec: map[string]string{
							"promql": "sum(rate(container_cpu_usage_seconds_total{container_name!=\"POD\"}[1m])) by (container_name)",
						},
					},
				},
				Total: v1.MetricSourceHolder{
					MetricSource: v1.MetricSource{
						Type: "Prometheus",
						MetricSourceSpec: map[string]string{
							"promql": "sum(rate(container_cpu_usage_seconds_total{container_name!=\"POD\"}[1m])) by (container_name)",
						},
					},
				},
			},
			want: `incremental: false
good:
    prometheus:
        promql: sum(rate(container_cpu_usage_seconds_total{container_name!="POD"}[1m])) by (container_name)
total:
    prometheus:
        promql: sum(rate(container_cpu_usage_seconds_total{container_name!="POD"}[1m])) by (container_name)
`,
		},
	}

	for _, tt := range tests {
		tt := tt // https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getN9CountMetrics(tt.ratioMetric)

			if (err != nil) != tt.wantErr {
				t.Errorf("getCountMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				y, _ := yaml.Marshal(got)
				assert.Equal(t, tt.want, string(y))
			}
		})
	}
}

func Test_getMetricSource(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    v1.MetricSource
		want    string
		wantErr bool
	}{
		{
			name: "Unsupported",
			args: v1.MetricSource{
				Type: "Unsupported",
			},
			wantErr: true,
		},
		{
			name: "Prometheus",
			args: v1.MetricSource{
				Type: "Prometheus",
				MetricSourceSpec: map[string]string{
					"promql": "sum(rate(container_cpu_usage_seconds_total{container_name!=\"POD\"}[1m])) by (container_name)",
				},
			},
			want: `prometheus:
    promql: sum(rate(container_cpu_usage_seconds_total{container_name!="POD"}[1m])) by (container_name)
`,
		},
		{
			name: "AmazonPrometheus",
			args: v1.MetricSource{
				Type: "AmazonPrometheus",
				MetricSourceSpec: map[string]string{
					"promql": "myapp_server_requestMsec{host=\"*\",job=\"nginx\"}",
				},
			},
			want: `amazonPrometheus:
    promql: myapp_server_requestMsec{host="*",job="nginx"}
`,
		},
		{
			name: "Datadog",
			args: v1.MetricSource{
				Type: "Datadog",
				MetricSourceSpec: map[string]string{
					"query": "sum:rate:container.cpu{container_name!=\"POD\"} by {container_name}",
				},
			},
			want: `datadog:
    query: sum:rate:container.cpu{container_name!="POD"} by {container_name}
`,
		},
		{
			name: "NewRelic",
			args: v1.MetricSource{
				Type: "NewRelic",
				MetricSourceSpec: map[string]string{
					"nrql": "SELECT sum(duration) FROM Transaction WHERE name = 'WebTransaction'",
				},
			},
			want: `newRelic:
    nrql: SELECT sum(duration) FROM Transaction WHERE name = 'WebTransaction'
`,
		},
		{
			name: "ThousandEyes",
			args: v1.MetricSource{
				Type: "ThousandEyes",
				MetricSourceSpec: map[string]string{
					"TestID":   "1234",
					"TestType": "mytype",
				},
			},
			want: `thousandEyes:
    testID: 1234
    testType: mytype
`,
		},
		{
			name: "AppDynamics",
			args: v1.MetricSource{
				Type: "AppDynamics",
				MetricSourceSpec: map[string]string{
					"applicationName": "myapp",
					"metricPath":      "mypath",
				},
			},
			want: `appDynamics:
    applicationName: myapp
    metricPath: mypath
`,
		},
		{
			name: "Splunk",
			args: v1.MetricSource{
				Type: "Splunk",
				MetricSourceSpec: map[string]string{
					"query": "mysplunkquery",
				},
			},
			want: `splunk:
    query: mysplunkquery
`,
		},
		{
			name: "Lightstep",
			args: v1.MetricSource{
				Type: "Lightstep",
				MetricSourceSpec: map[string]string{
					"streamId":   "mystreamid",
					"typeOfData": "mytypeofdata",
					"percentile": "0.96",
				},
			},
			want: `lightstep:
    streamId: mystreamid
    typeOfData: mytypeofdata
    percentile: 0.96
`,
		},
		{
			name: "SplunkObservability",
			args: v1.MetricSource{
				Type: "SplunkObservability",
				MetricSourceSpec: map[string]string{
					"program": "myprogram",
				},
			},
			want: `splunkObservability:
    program: myprogram
`,
		},
		{
			name: "Dynatrace",
			args: v1.MetricSource{
				Type: "Dynatrace",
				MetricSourceSpec: map[string]string{
					"metricSelector": "mymetricselector",
				},
			},
			want: `dynatrace:
    metricSelector: mymetricselector
`,
		},
		{
			name: "Elasticsearch",
			args: v1.MetricSource{
				Type: "Elasticsearch",
				MetricSourceSpec: map[string]string{
					"index": "myindex",
					"query": "myquery",
				},
			},
			want: `elasticsearch:
    index: myindex
    query: myquery
`,
		},
		{
			name: "CloudWatch",
			args: v1.MetricSource{
				Type: "CloudWatch",
				MetricSourceSpec: map[string]string{
					"namespace":  "mynamespace",
					"metricName": "mymetricname",
					"region":     "myregion",
					"stat":       "mystat",
					"dimensions": "name:mydimensions,value:myvalue;name:mydimensions2,value:myvalue2",
					"sql":        "myquery",
					"json":       "myjson",
				},
			},
			want: `cloudWatch:
    region: myregion
    namespace: mynamespace
    metricName: mymetricname
    stat: mystat
    dimensions:
        - name: mydimensions
          value: myvalue
        - name: mydimensions2
          value: myvalue2
    sql: myquery
    json: myjson
`,
		},
		{
			name: "Redshift",
			args: v1.MetricSource{
				Type: "Redshift",
				MetricSourceSpec: map[string]string{
					"query":        "myquery",
					"region":       "myregion",
					"clusterId":    "myclusterid",
					"databaseName": "mydatabasename",
				},
			},
			want: `redshift:
    region: myregion
    clusterId: myclusterid
    databaseName: mydatabasename
    query: myquery
`,
		},
		{
			name: "SumoLogic",
			args: v1.MetricSource{
				Type: "SumoLogic",
				MetricSourceSpec: map[string]string{
					"type":         "mytype",
					"query":        "myquery",
					"quantization": "myquantization",
					"rollup":       "myrollup",
				},
			},
			want: `sumoLogic:
    type: mytype
    query: myquery
    quantization: myquantization
    rollup: myrollup
`,
		},
		{
			name: "Instana",
			args: v1.MetricSource{
				Type: "Instana",
				MetricSourceSpec: map[string]string{
					"metricType":                            "mymetrictype",
					"infrastructure.metricRetrievalMethod":  "myInfrastructureMetricRetrivalMethod",
					"infrastructure.query":                  "myInfrastructureQuery",
					"infrastructure.snapshotId":             "myInfrastructureSnapshotId",
					"infrastructure.metricId":               "myInfrastructureMetricId",
					"infrastructure.pluginId":               "myInfrastructurePluginId",
					"application.metricId":                  "myapplicationMetricId",
					"application.aggregation":               "myapplicationAggregation",
					"application.groupBy.tag":               "myapplicationGroupByTag",
					"application.groupBy.tagEntity":         "myapplicationGroupByTagEntity",
					"application.groupBy.tagSecondLevelKey": "myapplicationTagSecondLevelKey",
					"application.apiQuery":                  "myapplicationApiQuery",
					"application.includeInternal":           "true",
					"application.includeSynthetic":          "false",
				},
			},
			want: `instana:
    metricType: mymetrictype
    infrastructure:
        metricRetrievalMethod: myInfrastructureMetricRetrivalMethod
        query: myInfrastructureQuery
        snapshotId: myInfrastructureSnapshotId
        metricId: myInfrastructureMetricId
        pluginId: myInfrastructurePluginId
    application:
        metricId: myapplicationMetricId
        aggregation: myapplicationAggregation
        groupBy:
            tag: myapplicationGroupByTag
            tagEntity: myapplicationGroupByTagEntity
            tagSecondLevelKey: myapplicationTagSecondLevelKey
        apiQuery: myapplicationApiQuery
`,
		},
		{
			name: "Pingdom",
			args: v1.MetricSource{
				Type: "Pingdom",
				MetricSourceSpec: map[string]string{
					"checkId":   "mycheckid",
					"checkType": "mychecktype",
					"status":    "mystatus",
				},
			},
			want: `pingdom:
    checkId: mycheckid
    checkType: mychecktype
    status: mystatus
`,
		},
		{
			name: "Graphite",
			args: v1.MetricSource{
				Type: "Graphite",
				MetricSourceSpec: map[string]string{
					"metricPath": "mymetricpath",
				},
			},
			want: `graphite:
    metricPath: mymetricpath
`,
		},
		{
			name: "BigQuery",
			args: v1.MetricSource{
				Type: "BigQuery",
				MetricSourceSpec: map[string]string{
					"projectId": "myprojectid",
					"query":     "myquery",
					"location":  "mylocation",
				},
			},
			want: `bigQuery:
    query: myquery
    projectId: myprojectid
    location: mylocation
`,
		},
		{
			name: "OpenTSDB",
			args: v1.MetricSource{
				Type: "OpenTSDB",
				MetricSourceSpec: map[string]string{
					"query": "myquery",
				},
			},
			want: `opentsdb:
    query: myquery
`,
		},
		{
			name: "GrafanaLoki",
			args: v1.MetricSource{
				Type: "GrafanaLoki",
				MetricSourceSpec: map[string]string{
					"logql": "mylogql",
				},
			},
			want: `grafanaLoki:
    logql: mylogql
`,
		},
		{
			name: "GoogleCloudMonitoring",
			args: v1.MetricSource{
				Type: "GoogleCloudMonitoring",
				MetricSourceSpec: map[string]string{
					"projectId": "myprojectid",
					"query":     "myquery",
				},
			},
			want: `gcm:
    query: myquery
    projectId: myprojectid
`,
		},
	}
	for _, tt := range tests {
		tt := tt // https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getN9MetricSource(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMetricSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				y, _ := yaml.Marshal(got)
				assert.Equal(t, tt.want, string(y))
			}
		})
	}
}

func Test_getN9Indicator(t *testing.T) {
	t.Parallel()
	t.Run("unknown", func(t *testing.T) {
		t.Parallel()
		var s v1.SLIInline
		ind := getN9Indicator(&s, "foo")
		assert.Equal(t, "ChangeMe", ind.MetricSource.Name)
	})
	t.Run("Total MetricSource exists", func(t *testing.T) {
		t.Parallel()
		s := v1.SLIInline{
			Spec: v1.SLISpec{
				RatioMetric: &v1.RatioMetric{
					Total: v1.MetricSourceHolder{
						MetricSource: v1.MetricSource{
							MetricSourceRef: "baz",
						},
					},
				},
			},
		}
		ind := getN9Indicator(&s, "foo")
		assert.Equal(t, "baz", ind.MetricSource.Name)
	})
	t.Run("Only Good MetricSource exists", func(t *testing.T) {
		t.Parallel()
		// if total doesn't exists, use default behavior and generate ChangeMe indicator
		s := v1.SLIInline{
			Spec: v1.SLISpec{
				RatioMetric: &v1.RatioMetric{
					Good: &v1.MetricSourceHolder{
						MetricSource: v1.MetricSource{
							MetricSourceRef: "baz",
						},
					},
				},
			},
		}
		ind := getN9Indicator(&s, "foo")
		assert.Equal(t, "ChangeMe", ind.MetricSource.Name)
	})
}
