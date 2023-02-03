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
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	nobl9v1alpha "github.com/OpenSLO/oslo/internal/pkg/manifest/nobl9/v1alpha"
	"github.com/OpenSLO/oslo/pkg/manifest"
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

func Test_RemoveDuplicates(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "empty array",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "All unique",
			input: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "Some dupes",
			input: []string{"a", "b", "c", "b", "a"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "All the same",
			input: []string{"a", "a", "a"},
			want:  []string{"a"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := RemoveDuplicates(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getParsedObjects(t *testing.T) {
	t.Parallel()
	// needed farther down to test an empty list
	var empty []manifest.OpenSLOKind
	// needed here so we can pass in the pointer address later
	sliName := "foo-sli"

	tests := []struct {
		name    string
		args    []string
		want    []manifest.OpenSLOKind
		wantErr bool
	}{
		{
			name: "empty list",
			args: []string{},
			want: empty,
		},
		{
			name: "Single DataSource per file",
			args: []string{"../../../test/v1/data-source/data-source.yaml"},
			want: []manifest.OpenSLOKind{
				v1.DataSource{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "DataSource",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name:        "TestDataSource",
								DisplayName: "Test Data Source",
							},
						},
					},
					Spec: v1.DataSourceSpec{
						Type: "CloudWatch",
						ConnectionDetails: map[string]string{
							"accessKeyID":     "accessKey",
							"secretAccessKey": "secretAccessKey",
						},
					},
				},
			},
		},
		{
			name: "Single Service file",
			args: []string{"../../../test/v1/service/service.yaml"},
			want: []manifest.OpenSLOKind{
				v1.Service{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "Service",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name:        "my-rad-service",
								DisplayName: "My Rad Service",
							},
						},
					},
					Spec: v1.ServiceSpec{
						Description: "This is a great description of an even better service.",
					},
				},
			},
		},
		{
			name: "Single SLI per file",
			args: []string{"../../../test/v1/sli/sli-description-threshold-metricsourceref.yaml"},
			want: []manifest.OpenSLOKind{
				v1.SLI{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "SLI",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name:        "GreatSLI",
								DisplayName: "Great SLI",
							},
						},
					},
					Spec: v1.SLISpec{
						ThresholdMetric: &v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								MetricSourceRef: "redshift-datasource",
								MetricSourceSpec: map[string]string{
									"clusterId":    "metrics-cluster",
									"databaseName": "metrics-db",
									"query":        "SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to",
									"region":       "eu-central-1",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple Files",
			args: []string{
				"../../../test/v1/service/service.yaml",
				"../../../test/v1/sli/sli-description-threshold-metricsourceref.yaml",
			},
			want: []manifest.OpenSLOKind{
				v1.Service{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "Service",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name:        "my-rad-service",
								DisplayName: "My Rad Service",
							},
						},
					},
					Spec: v1.ServiceSpec{
						Description: "This is a great description of an even better service.",
					},
				},
				v1.SLI{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "SLI",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name:        "GreatSLI",
								DisplayName: "Great SLI",
							},
						},
					},
					Spec: v1.SLISpec{
						ThresholdMetric: &v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								MetricSourceRef: "redshift-datasource",
								MetricSourceSpec: map[string]string{
									"clusterId":    "metrics-cluster",
									"databaseName": "metrics-db",
									"query":        "SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to",
									"region":       "eu-central-1",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple definitions per file",
			args: []string{"../../../test/v1/multi.yaml"},
			want: []manifest.OpenSLOKind{
				v1.SLO{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "SLO",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name:        "foo-slo",
								DisplayName: "FOO SLO",
							},
						},
					},
					Spec: v1.SLOSpec{
						Description:     "Foo SLO",
						Service:         "foo-slos",
						IndicatorRef:    &sliName,
						BudgetingMethod: "Occurrences",
						TimeWindow: []v1.TimeWindow{
							{
								Duration:  "28d",
								IsRolling: true,
							},
						},
						Objectives: []v1.Objective{
							{
								DisplayName: "Test Objective",
								Op:          "gte",
								Value:       10,
								Target:      0.98,
							},
						},
						AlertPolicies: []string{},
					},
				},
				v1.SLI{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "SLI",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name: "foo-sli",
							},
						},
					},
					Spec: v1.SLISpec{
						ThresholdMetric: &v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								MetricSourceRef: "foo-cloudwatch",
								Type:            "CloudWatch",
								MetricSourceSpec: map[string]string{
									"dimensions": "name:CanaryName,value:web-app",
									"metricName": "2xx",
									"namespace":  "CloudWatchSynthetics",
									"region":     "us-east-1",
									"stat":       "SampleCount",
								},
							},
						},
					},
				},
				v1.DataSource{
					ObjectHeader: v1.ObjectHeader{
						ObjectHeader: manifest.ObjectHeader{
							APIVersion: "openslo/v1",
						},
						Kind: "DataSource",
						MetadataHolder: v1.MetadataHolder{
							Metadata: v1.Metadata{
								Name: "foo-cloudwatch",
							},
						},
					},
					Spec: v1.DataSourceSpec{
						Type: "CloudWatch",
						ConnectionDetails: map[string]string{
							"accessKeyID":     "FOOBAR",
							"secretAccessKey": "BAZBAT",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt // https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, _, err := getParsedObjects(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getParsedObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Nobl9(t *testing.T) {
	t.Parallel()
	type args struct {
		filenames []string
		project   string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "Single Service Single file",
			args: args{
				filenames: []string{
					"../../../test/v1/service/service-with-labels.yaml",
				},
				project: "default",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: Service
metadata:
    name: my-rad-service
    displayName: My Rad Service
    project: default
    labels:
        costCentre:
            - project1
        serviceTier:
            - tier-1
        team:
            - identity
        userImpacting:
            - "true"
spec:
    description: This is a great description of an even better service.
`,
		},
		{
			name: "Single SLO single file",
			args: args{
				filenames: []string{
					"../../../test/v1/slo/slo-indicatorref-rolling-alerts.yaml",
				},
				project: "default",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: SLO
metadata:
    name: TestSLO
    displayName: Test SLO
    project: default
spec:
    description: This is a great description
    indicator:
        metricSource:
            project: default
            name: Changeme
            kind: Agent
    budgetingMethod: Occurrences
    objectives:
        - displayName: Foo Total Errors
          value: 1
          target: 0.98
          op: lt
    service: TheServiceName
    timeWindows:
        - unit: Day
          count: 1
          isRolling: true
    alertPolicies:
        - FooAlertPolicy
`,
		},
		{
			name: "Multiple Kinds Single File",
			args: args{
				filenames: []string{
					"../../../test/v1/multi.yaml",
				},
				project: "default",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: SLO
metadata:
    name: foo-slo
    displayName: FOO SLO
    project: default
spec:
    description: Foo SLO
    indicator:
        metricSource:
            project: default
            name: foo-cloudwatch
            kind: Agent
    budgetingMethod: Occurrences
    objectives:
        - displayName: Test Objective
          value: 10
          target: 0.98
          rawMetric:
            query:
                cloudWatch:
                    region: us-east-1
                    namespace: CloudWatchSynthetics
                    metricName: 2xx
                    stat: SampleCount
                    dimensions:
                        - name: CanaryName
                          value: web-app
          op: gte
    service: foo-slos
    timeWindows:
        - unit: Day
          count: 28
          isRolling: true
    alertPolicies: []
`,
		},
		{
			name: "Multiple Kinds Multiple Files",
			args: args{
				filenames: []string{
					"../../../test/v1/slo/slo-indicatorRef-rolling-cloudwatch.yaml",
					"../../../test/v1/sli/sli-threshold-cloudwatch.yaml",
					"../../../test/v1/data-source/data-source-cloudwatch.yaml",
				},
				project: "default",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: SLO
metadata:
    name: foo-openslo-slo
    displayName: FOO OPENSLO SLO
    project: default
spec:
    description: ""
    indicator:
        metricSource:
            project: default
            name: foo-cloudwatch
            kind: Agent
    budgetingMethod: Occurrences
    objectives:
        - displayName: Test Objective
          value: 10
          target: 0.98
          rawMetric:
            query:
                cloudWatch:
                    region: us-east-1
                    namespace: CloudWatchSynthetics
                    metricName: 2xx
                    stat: SampleCount
                    dimensions:
                        - name: CanaryName
                          value: web-app
          op: gte
    service: foo-slos
    timeWindows:
        - unit: Day
          count: 28
          isRolling: true
    alertPolicies: []
`,
		},
		{
			name: "Test adding annotation to chose indicator kind",
			args: args{
				filenames: []string{
					"../../../test/v1/slo/slo-with-annotations.yaml",
				},
				project: "default",
			},
			wantOut: `---
apiVersion: n9/v1alpha
kind: SLO
metadata:
    name: monthy-openslo-slo
    displayName: Python
    project: default
spec:
    description: ""
    indicator:
        metricSource:
            project: default
            name: Changeme
            kind: Direct
    budgetingMethod: Occurrences
    objectives:
        - displayName: Life of Brian
          value: 10
          target: 0.98
          op: gte
    service: foo-slos
    timeWindows:
        - unit: Day
          count: 28
          isRolling: true
    alertPolicies: []
`,
		},
	}
	for _, tt := range tests {
		tt := tt // https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out := &bytes.Buffer{}
			if err := Nobl9(out, tt.args.filenames, tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("Nobl9() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				assert.Equal(t, tt.wantOut, gotOut)
			}
		})
	}
}

func Test_getN9Indicator(t *testing.T) {
	t.Parallel()
	var nilIndicator v1.SLISpec
	type args struct {
		indicator v1.SLISpec
		project   string
	}
	tests := []struct {
		name string
		args args
		want nobl9v1alpha.Indicator
	}{
		{
			name: "nil Indicator",
			args: args{
				indicator: nilIndicator,
				project:   "default",
			},
			want: nobl9v1alpha.Indicator{
				MetricSource: nobl9v1alpha.MetricSourceSpec{
					Project: "default",
					Name:    "Changeme",
					Kind:    "Agent",
				},
			},
		},
		{
			name: "Empty Indicator",
			args: args{
				indicator: nilIndicator,
			},
			want: nobl9v1alpha.Indicator{
				MetricSource: nobl9v1alpha.MetricSourceSpec{
					Project: "default",
					Name:    "Changeme",
					Kind:    "Agent",
				},
			},
		},
		{
			name: "Empty RatioMetric",
			args: args{
				indicator: v1.SLISpec{
					RatioMetric: &v1.RatioMetric{},
				},
				project: "FooBar",
			},
			want: nobl9v1alpha.Indicator{
				MetricSource: nobl9v1alpha.MetricSourceSpec{
					Project: "FooBar",
					Name:    "Changeme",
					Kind:    "Agent",
				},
			},
		},
		{
			name: "Empty ThresholdMetric",
			args: args{
				indicator: v1.SLISpec{
					ThresholdMetric: &v1.MetricSourceHolder{},
				},
				project: "FooBar",
			},
			want: nobl9v1alpha.Indicator{
				MetricSource: nobl9v1alpha.MetricSourceSpec{
					Project: "FooBar",
					Name:    "Changeme",
					Kind:    "Agent",
				},
			},
		},
		{
			name: "Ratio Metric with good and total",
			args: args{
				indicator: v1.SLISpec{
					RatioMetric: &v1.RatioMetric{
						Good: &v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								MetricSourceRef: "foo-bar-sli",
							},
						},
						Total: v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								MetricSourceRef: "foo-bar-sli",
							},
						},
					},
				},
				project: "FooBar",
			},
			want: nobl9v1alpha.Indicator{
				MetricSource: nobl9v1alpha.MetricSourceSpec{
					Project: "FooBar",
					Name:    "foo-bar-sli",
					Kind:    "Agent",
				},
			},
		},
		{
			name: "Ratio Metric with bad and total",
			args: args{
				indicator: v1.SLISpec{
					RatioMetric: &v1.RatioMetric{
						Bad: &v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								MetricSourceRef: "foo-bar-sli",
							},
						},
						Total: v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								MetricSourceRef: "foo-bar-bad-sli",
							},
						},
					},
				},
				project: "FooBar",
			},
			want: nobl9v1alpha.Indicator{
				MetricSource: nobl9v1alpha.MetricSourceSpec{
					Project: "FooBar",
					Name:    "foo-bar-bad-sli",
					Kind:    "Agent",
				},
			},
		},
		{
			name: "Threshold Metric",
			args: args{
				indicator: v1.SLISpec{
					ThresholdMetric: &v1.MetricSourceHolder{
						MetricSource: v1.MetricSource{
							MetricSourceRef: "thresh-foo-sli",
						},
					},
				},
				project: "FooBarThresh",
			},
			want: nobl9v1alpha.Indicator{
				MetricSource: nobl9v1alpha.MetricSourceSpec{
					Project: "FooBarThresh",
					Name:    "thresh-foo-sli",
					Kind:    "Agent",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt // https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getN9Indicator(tt.args.indicator, []string{}, tt.args.project)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getN9Thresholds(t *testing.T) {
	t.Parallel()
	type args struct {
		o         []v1.Objective
		indicator v1.SLISpec
	}
	tests := []struct {
		name    string
		args    args
		want    []nobl9v1alpha.Threshold
		wantErr bool
	}{
		{
			name: "Single Objective, no Indicator",
			args: args{
				o: []v1.Objective{
					{
						DisplayName: "foo1",
						Op:          "gte",
						Value:       0.01,
						Target:      100,
					},
				},
			},
			want: []nobl9v1alpha.Threshold{
				{
					ThresholdBase: nobl9v1alpha.ThresholdBase{
						DisplayName: "foo1",
						Value:       0.01,
					},
					// anonymous function since we have to pass the address
					BudgetTarget: func() *float64 { i := float64(100); return &i }(),
					Operator:     func() *string { i := "gte"; return &i }(), //nolint:goconst
				},
			},
		},
		{
			name: "Single Objective, RatioMetric Good Indicator",
			args: args{
				o: []v1.Objective{
					{
						DisplayName: "foo1",
						Op:          "gte",
						Value:       0.01,
						Target:      100,
					},
				},
				indicator: v1.SLISpec{
					RatioMetric: &v1.RatioMetric{
						Counter: true,
						Good: &v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								Type: "Datadog",
								MetricSourceSpec: map[string]string{
									"query": "foo",
								},
							},
						},
						Total: v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								Type: "Datadog",
								MetricSourceSpec: map[string]string{
									"query": "bar",
								},
							},
						},
					},
				},
			},
			want: []nobl9v1alpha.Threshold{
				{
					ThresholdBase: nobl9v1alpha.ThresholdBase{
						DisplayName: "foo1",
						Value:       0.01,
					},
					// anonymous function since we have to pass the address
					BudgetTarget: func() *float64 { i := float64(100); return &i }(),
					CountMetrics: &nobl9v1alpha.CountMetricsSpec{
						Incremental: func() *bool { i := true; return &i }(),
						GoodMetric: &nobl9v1alpha.MetricSpec{
							Datadog: &nobl9v1alpha.DatadogMetric{
								Query: func() *string { i := "foo"; return &i }(),
							},
						},
						TotalMetric: &nobl9v1alpha.MetricSpec{
							Datadog: &nobl9v1alpha.DatadogMetric{
								Query: func() *string { i := "bar"; return &i }(),
							},
						},
					},
					Operator: func() *string { i := "gte"; return &i }(),
				},
			},
		},
		{
			name: "Single Objective, RatioMetric Bad Indicator",
			args: args{
				o: []v1.Objective{
					{
						DisplayName: "foo1",
						Op:          "gte",
						Value:       0.01,
						Target:      100,
					},
				},
				indicator: v1.SLISpec{
					RatioMetric: &v1.RatioMetric{
						Counter: true,
						Bad: &v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								Type: "Datadog",
								MetricSourceSpec: map[string]string{
									"query": "foo",
								},
							},
						},
						Total: v1.MetricSourceHolder{
							MetricSource: v1.MetricSource{
								Type: "Datadog",
								MetricSourceSpec: map[string]string{
									"query": "bar",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Single Objective, Threshold Indicator",
			args: args{
				o: []v1.Objective{
					{
						DisplayName: "foo1",
						Op:          "gte",
						Value:       0.01,
						Target:      100,
					},
				},
				indicator: v1.SLISpec{
					ThresholdMetric: &v1.MetricSourceHolder{
						MetricSource: v1.MetricSource{
							Type: "NewRelic",
							MetricSourceSpec: map[string]string{
								"nrql": "foo-bar",
							},
						},
					},
				},
			},
			want: []nobl9v1alpha.Threshold{
				{
					ThresholdBase: nobl9v1alpha.ThresholdBase{
						DisplayName: "foo1",
						Value:       0.01,
					},
					// anonymous function since we have to pass the address
					BudgetTarget: func() *float64 { i := float64(100); return &i }(),
					RawMetric: &nobl9v1alpha.RawMetricSpec{
						MetricQuery: &nobl9v1alpha.MetricSpec{
							NewRelic: &nobl9v1alpha.NewRelicMetric{
								NRQL: func() *string { i := "foo-bar"; return &i }(), //nolint:goconst
							},
						},
					},
					Operator: func() *string { i := "gte"; return &i }(),
				},
			},
		},
		{
			name: "Multiple Objectives, Threshold Indicator",
			args: args{
				o: []v1.Objective{
					{
						DisplayName: "foo1",
						Op:          "gte",
						Value:       0.98,
						Target:      200,
					},
					{
						DisplayName: "foo1",
						Op:          "gte",
						Value:       0.01,
						Target:      100,
					},
				},
				indicator: v1.SLISpec{
					ThresholdMetric: &v1.MetricSourceHolder{
						MetricSource: v1.MetricSource{
							Type: "NewRelic",
							MetricSourceSpec: map[string]string{
								"nrql": "foo-bar",
							},
						},
					},
				},
			},
			want: []nobl9v1alpha.Threshold{
				{
					ThresholdBase: nobl9v1alpha.ThresholdBase{
						DisplayName: "foo1",
						Value:       0.98,
					},
					// anonymous function since we have to pass the address
					BudgetTarget: func() *float64 { i := float64(200); return &i }(),
					RawMetric: &nobl9v1alpha.RawMetricSpec{
						MetricQuery: &nobl9v1alpha.MetricSpec{
							NewRelic: &nobl9v1alpha.NewRelicMetric{
								NRQL: func() *string { i := "foo-bar"; return &i }(),
							},
						},
					},
					Operator: func() *string { i := "gte"; return &i }(),
				},
				{
					ThresholdBase: nobl9v1alpha.ThresholdBase{
						DisplayName: "foo1",
						Value:       0.01,
					},
					// anonymous function since we have to pass the address
					BudgetTarget: func() *float64 { i := float64(100); return &i }(),
					RawMetric: &nobl9v1alpha.RawMetricSpec{
						MetricQuery: &nobl9v1alpha.MetricSpec{
							NewRelic: &nobl9v1alpha.NewRelicMetric{
								NRQL: func() *string { i := "foo-bar"; return &i }(),
							},
						},
					},
					Operator: func() *string { i := "gte"; return &i }(),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt // https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getN9Thresholds(tt.args.o, tt.args.indicator)
			if (err != nil) != tt.wantErr {
				t.Errorf("getN9Thresholds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getN9CloudWatchQuery(t *testing.T) {
	t.Parallel()
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
		want nobl9v1alpha.CloudWatchMetric
	}{
		{
			name: "has dimensions",
			args: args{
				m: map[string]string{
					"namespace":  "mynamespace",
					"metricName": "mymetricname",
					"region":     "myregion",
					"stat":       "mystat",
					"dimensions": "name:mydimensions,value:myvalue;name:mydimensions2,value:myvalue2",
				},
			},
			want: nobl9v1alpha.CloudWatchMetric{
				Stat: func() *string { i := "mystat"; return &i }(),
				Dimensions: []nobl9v1alpha.CloudWatchMetricDimension{
					{
						Name:  func() *string { i := "mydimensions"; return &i }(),
						Value: func() *string { i := "myvalue"; return &i }(),
					},
					{
						Name:  func() *string { i := "mydimensions2"; return &i }(),
						Value: func() *string { i := "myvalue2"; return &i }(),
					},
				},
			},
		},
		{
			name: "has json",
			args: args{
				m: map[string]string{
					"namespace":  "mynamespace",
					"metricName": "mymetricname",
					"region":     "myregion",
					"json":       "{\"foo\": \"bar\"}",
				},
			},
			want: nobl9v1alpha.CloudWatchMetric{
				JSON: func() *string { i := "{\"foo\": \"bar\"}"; return &i }(),
			},
		},
		{
			name: "has SQL",
			args: args{
				m: map[string]string{
					"namespace":  "mynamespace",
					"metricName": "mymetricname",
					"region":     "myregion",
					"sql":        "SELECT * FROM FOO",
				},
			},
			want: nobl9v1alpha.CloudWatchMetric{
				SQL: func() *string { i := "SELECT * FROM FOO"; return &i }(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt // https://gist.github.com/kunwardeep/80c2e9f3d3256c894898bae82d9f75d0
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getN9CloudWatchQuery(tt.args.m)
			assert.Equal(t, tt.want, got)
		})
	}
}
