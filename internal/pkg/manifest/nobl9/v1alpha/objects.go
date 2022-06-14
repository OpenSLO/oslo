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
package nobl9v1alpha

import (
	nobl9manifest "github.com/OpenSLO/oslo/internal/pkg/manifest/nobl9"
)

// apiVersion: n9/v1alpha
// kind: Service
// metadata:
//   name: # string
//   displayName: # string
//   project: # string
// spec:
//   description: # string
//   serviceType: # string
const (
	APIVersion = "n9/v1alpha"
)

// Possible values of field kind for valid Objects.
const (
	KindSLO     = "SLO"
	KindService = "Service"
)

// Labels represents a set of labels.
type Labels map[string][]string

// Metadata represents part of object which is is common for all available Objects, for internal usage.
type Metadata struct {
	Name        string `yaml:"name" validate:"required" example:"name"`
	DisplayName string `yaml:"displayName,omitempty" validate:"omitempty,min=0,max=63" example:"Prometheus Source"`
	Project     string `yaml:"project,omitempty" validate:"objectName" example:"default"`
	Labels      Labels `yaml:"labels,omitempty" validate:"omitempty,labels"`
}

// MetadataHolder is an intermediate structure that can provides metadata related
// field to other structures.
type MetadataHolder struct {
	Metadata Metadata `yaml:"metadata"`
}

// ObjectHeader is a header for all objects.
type ObjectHeader struct {
	nobl9manifest.ObjectHeader `yaml:",inline"`
	Kind                       string `yaml:"kind" validate:"required,oneof=Service SLO AlertNotificationTarget" example:"kind"` //nolint:lll
	MetadataHolder             `yaml:",inline"`
}

// Service struct which mapped one to one with kind: service yaml definition.
type Service struct {
	ObjectHeader `yaml:",inline"`
	Spec         ServiceSpec `yaml:"spec"`
}

// Kind returns the name of this type.
func (Service) Kind() string {
	return "Service"
}

// ServiceSpec represents content of Spec typical for Service Object.
type ServiceSpec struct {
	Description string `yaml:"description" validate:"max=1050" example:"Bleeding edge web app"`
}

// AlertPolicy represents a set of conditions that can trigger an alert.
type AlertPolicy struct {
	ObjectHeader `yaml:",inline"`
	Spec         AlertPolicySpec `yaml:"spec"`
}

// AlertPolicySpec represents content of AlertPolicy's Spec.
type AlertPolicySpec struct {
	Description      string           `yaml:"description" validate:"description" example:"Error budget is at risk"`
	Severity         string           `yaml:"severity" validate:"required,severity" example:"High"`
	CoolDownDuration string           `yaml:"coolDown,omitempty" validate:"omitempty,validDuration" example:"5m"` //nolint:lll
	Conditions       []AlertCondition `yaml:"conditions" validate:"required,min=1,dive"`
	AlertMethods     []string         `yaml:"alertMethods"`
}

// AlertCondition represents a condition to meet to trigger an alert.
type AlertCondition struct {
	Measurement      string      `yaml:"measurement" validate:"required,alertPolicyMeasurement" example:"BurnedBudget"`
	Value            interface{} `yaml:"value" validate:"required" swaggertype:"string" example:"0.97"`
	LastsForDuration string      `yaml:"lastsFor,omitempty" validate:"omitempty,validDuration,nonNegativeDuration" example:"15m"` //nolint:lll
	Operation        string      `yaml:"op" validate:"required,alertOperation" example:"lt"`
}

// AlertMethodAssignment represents an AlertMethod assigned to AlertPolicy.
type AlertMethodAssignment struct {
	Project string `yaml:"project,omitempty" validate:"omitempty,objectName" example:"default"`
	Name    string `yaml:"name" validate:"required,objectName" example:"webhook-alertmethod"`
}

// SLO struct which mapped one to one with kind: slo yaml definition, external usage.
type SLO struct {
	ObjectHeader `yaml:",inline"`
	Spec         SLOSpec `yaml:"spec"`
}

// SLOSpec represents content of Spec typical for SLO Object.
type SLOSpec struct {
	Description     string       `yaml:"description" validate:"description" example:"Total count of server requests"` //nolint:lll
	Indicator       Indicator    `yaml:"indicator"`
	BudgetingMethod string       `yaml:"budgetingMethod" validate:"required,budgetingMethod" example:"Occurrences"`
	Thresholds      []Threshold  `yaml:"objectives" validate:"required,dive"`
	Service         string       `yaml:"service" validate:"required,objectName" example:"webapp-service"`
	TimeWindows     []TimeWindow `yaml:"timeWindows" validate:"required,len=1,dive"`
	AlertPolicies   []string     `yaml:"alertPolicies" validate:"omitempty"`
	Attachments     []Attachment `yaml:"attachments,omitempty" validate:"omitempty,len=1,dive"`
	CreatedAt       string       `yaml:"createdAt,omitempty"`
}

// ThresholdBase base structure representing a threshold.
type ThresholdBase struct {
	DisplayName string  `yaml:"displayName" validate:"objectiveDisplayName" example:"Good"`
	Value       float64 `yaml:"value" validate:"numeric" example:"100"`
}

// Threshold represents single threshold for SLO, for internal usage.
type Threshold struct {
	ThresholdBase `yaml:",inline"`
	BudgetTarget *float64 `yaml:"target" validate:"required,numeric,gte=0,lt=1" example:"0.9"`
	TimeSliceTarget *float64          `yaml:"timeSliceTarget,omitempty" example:"0.9"`
	CountMetrics    *CountMetricsSpec `yaml:"countMetrics,omitempty"`
	RawMetric       *RawMetricSpec    `yaml:"rawMetric,omitempty"`
	Operator        *string           `yaml:"op,omitempty" example:"lte"`
	Name            *string           `yaml:"name,omitempty"`
}

// Indicator represents integration with metric source can be. e.g. Prometheus, Datadog, for internal usage.
type Indicator struct {
	MetricSource MetricSourceSpec `yaml:"metricSource" validate:"required"`
	RawMetric    MetricSpec       `yaml:"rawMetric,omitempty"`
}

// Attachment represents user defined URL attached to SLO.
type Attachment struct {
	URL         string  `yaml:"url" validate:"required,url"`
	DisplayName *string `yaml:"displayName,omitempty"`
}

type Calendar struct {
	StartTime string `yaml:"startTime" validate:"required,dateWithTime,minDateTime" example:"2020-01-21 12:30:00"`
	TimeZone  string `yaml:"timeZone" validate:"required,timeZone" example:"America/New_York"`
}

// Period represents period of time.
type Period struct {
	Begin string `yaml:"begin"`
	End   string `yaml:"end"`
}

// TimeWindow represents content of time window.
type TimeWindow struct {
	Unit      string    `yaml:"unit" validate:"required,timeUnit" example:"Week"`
	Count     int       `yaml:"count" validate:"required,gt=0" example:"1"`
	IsRolling bool      `yaml:"isRolling" example:"true"`
	Calendar  *Calendar `yaml:"calendar,omitempty"`
}

// CountMetricsSpec represents set of two time series of good and total counts.
type CountMetricsSpec struct {
	Incremental *bool       `yaml:"incremental" validate:"required"`
	GoodMetric  *MetricSpec `yaml:"good" validate:"required"`
	TotalMetric *MetricSpec `yaml:"total" validate:"required"`
}

// RawMetricSpec represents integration with a metric source for a particular threshold.
type RawMetricSpec struct {
	MetricQuery *MetricSpec `yaml:"query" validate:"required"`
}

type MetricSourceSpec struct {
	Project string `yaml:"project,omitempty" validate:"omitempty,objectName" example:"default"`
	Name    string `yaml:"name" validate:"required,objectName" example:"prometheus-source"`
	Kind    string `yaml:"kind" validate:"omitempty,metricSourceKind" example:"Agent"`
}

// MetricSpec defines single time series obtained from data source.
type MetricSpec struct {
	Prometheus          *PrometheusMetric          `yaml:"prometheus,omitempty"`
	Datadog             *DatadogMetric             `yaml:"datadog,omitempty"`
	NewRelic            *NewRelicMetric            `yaml:"newRelic,omitempty"`
	AppDynamics         *AppDynamicsMetric         `yaml:"appDynamics,omitempty"`
	Splunk              *SplunkMetric              `yaml:"splunk,omitempty"`
	Lightstep           *LightstepMetric           `yaml:"lightstep,omitempty"`
	SplunkObservability *SplunkObservabilityMetric `yaml:"splunkObservability,omitempty"`
	Dynatrace           *DynatraceMetric           `yaml:"dynatrace,omitempty"`
	Elasticsearch       *ElasticsearchMetric       `yaml:"elasticsearch,omitempty"`
	ThousandEyes        *ThousandEyesMetric        `yaml:"thousandEyes,omitempty"`
	Graphite            *GraphiteMetric            `yaml:"graphite,omitempty"`
	BigQuery            *BigQueryMetric            `yaml:"bigQuery,omitempty"`
	OpenTSDB            *OpenTSDBMetric            `yaml:"opentsdb,omitempty"`
	GrafanaLoki         *GrafanaLokiMetric         `yaml:"grafanaLoki,omitempty"`
	CloudWatch          *CloudWatchMetric          `yaml:"cloudWatch,omitempty"`
	Pingdom             *PingdomMetric             `yaml:"pingdom,omitempty"`
	AmazonPrometheus    *AmazonPrometheusMetric    `yaml:"amazonPrometheus,omitempty"`
	Redshift            *RedshiftMetric            `yaml:"redshift,omitempty"`
	SumoLogic           *SumoLogicMetric           `yaml:"sumoLogic,omitempty"`
	Instana             *InstanaMetric             `yaml:"instana,omitempty"`
}

// PrometheusMetric represents metric from Prometheus.
type PrometheusMetric struct {
	PromQL *string `yaml:"promql" validate:"required" example:"cpu_usage_user{cpu=\"cpu-total\"}"`
}

// AmazonPrometheusMetric represents metric from Amazon Managed Prometheus.
type AmazonPrometheusMetric struct {
	PromQL *string `yaml:"promql" validate:"required" example:"cpu_usage_user{cpu=\"cpu-total\"}"`
}

// DatadogMetric represents metric from Datadog.
type DatadogMetric struct {
	Query *string `yaml:"query" validate:"required"`
}

// NewRelicMetric represents metric from NewRelic.
type NewRelicMetric struct {
	NRQL *string `yaml:"nrql" validate:"required"`
}

// ThousandEyesMetric represents metric from ThousandEyes.
type ThousandEyesMetric struct {
	TestID   *int64  `yaml:"testID" validate:"required,gte=0"`
	TestType *string `yaml:"testType" validate:"supportedThousandEyesTestType"`
}

// AppDynamicsMetric represents metric from AppDynamics.
type AppDynamicsMetric struct {
	ApplicationName *string `yaml:"applicationName" validate:"required,notEmpty"`
	MetricPath      *string `yaml:"metricPath" validate:"required,unambiguousAppDynamicMetricPath"`
}

// SplunkMetric represents metric from Splunk.
type SplunkMetric struct {
	Query *string `yaml:"query" validate:"required,notEmpty,splunkQueryValid"`
}

// LightstepMetric represents metric from Lightstep.
type LightstepMetric struct {
	StreamID   *string  `yaml:"streamId" validate:"required"`
	TypeOfData *string  `yaml:"typeOfData" validate:"required,oneof=latency error_rate good total"`
	Percentile *float64 `yaml:"percentile,omitempty"`
}

// SplunkObservabilityMetric represents metric from SplunkObservability.
type SplunkObservabilityMetric struct {
	Program *string `yaml:"program" validate:"required"`
}

// DynatraceMetric represents metric from Dynatrace.
type DynatraceMetric struct {
	MetricSelector *string `yaml:"metricSelector" validate:"required"`
}

// ElasticsearchMetric represents metric from Elasticsearch.
type ElasticsearchMetric struct {
	Index *string `yaml:"index" validate:"required"`
	Query *string `yaml:"query" validate:"required"`
}

// CloudWatchMetric represents metric from CloudWatch.
type CloudWatchMetric struct {
	Region     *string                     `yaml:"region" validate:"required,max=255"`
	Namespace  *string                     `yaml:"namespace,omitempty"`
	MetricName *string                     `yaml:"metricName,omitempty"`
	Stat       *string                     `yaml:"stat,omitempty"`
	Dimensions []CloudWatchMetricDimension `yaml:"dimensions,omitempty" validate:"max=10,uniqueDimensionNames,dive"`
	SQL        *string                     `yaml:"sql,omitempty"`
	JSON       *string                     `yaml:"json,omitempty"`
}

// RedshiftMetric represents metric from Redshift.
type RedshiftMetric struct {
	Region       *string `yaml:"region" validate:"required,max=255"`
	ClusterID    *string `yaml:"clusterId" validate:"required"`
	DatabaseName *string `yaml:"databaseName" validate:"required"`
	Query        *string `yaml:"query" validate:"required,redshiftRequiredColumns"`
}

// SumoLogicMetric represents metric from Sumo Logic.
type SumoLogicMetric struct {
	Type         *string `yaml:"type" validate:"required"`
	Query        *string `yaml:"query" validate:"required"`
	Quantization *string `yaml:"quantization,omitempty"`
	Rollup       *string `yaml:"rollup,omitempty"`
	// For struct level validation refer to sumoLogicStructValidation in pkg/manifest/v1alpha/validator.go
}

// InstanaMetric represents metric from Redshift.
type InstanaMetric struct {
	MetricType     string                           `yaml:"metricType" validate:"required,oneof=infrastructure application"` //nolint:lll
	Infrastructure *InstanaInfrastructureMetricType `yaml:"infrastructure,omitempty"`
	Application    *InstanaApplicationMetricType    `yaml:"application,omitempty"`
}

type InstanaInfrastructureMetricType struct {
	MetricRetrievalMethod string  `yaml:"metricRetrievalMethod" validate:"required,oneof=query snapshot"`
	Query                 *string `yaml:"query,omitempty"`
	SnapshotID            *string `yaml:"snapshotId,omitempty"`
	MetricID              string  `yaml:"metricId" validate:"required"`
	PluginID              string  `yaml:"pluginId" validate:"required"`
}

type InstanaApplicationMetricType struct {
	MetricID         string                          `yaml:"metricId" validate:"required,oneof=calls erroneousCalls errors latency"` //nolint:lll
	Aggregation      string                          `yaml:"aggregation" validate:"required"`
	GroupBy          InstanaApplicationMetricGroupBy `yaml:"groupBy" validate:"required"`
	APIQuery         string                          `yaml:"apiQuery" validate:"required,json"`
	IncludeInternal  bool                            `yaml:"includeInternal,omitempty"`
	IncludeSynthetic bool                            `yaml:"includeSynthetic,omitempty"`
}

type InstanaApplicationMetricGroupBy struct {
	Tag               string  `yaml:"tag" validate:"required"`
	TagEntity         string  `yaml:"tagEntity" validate:"required,oneof=DESTINATION SOURCE NOT_APPLICABLE"`
	TagSecondLevelKey *string `yaml:"tagSecondLevelKey,omitempty"`
}

// CloudWatchMetricDimension represents name/value pair that is part of the identity of a metric.
type CloudWatchMetricDimension struct {
	Name  *string `yaml:"name" validate:"required,max=255,ascii,notBlank"`
	Value *string `yaml:"value" validate:"required,max=255,ascii,notBlank"`
}

// PingdomMetric represents metric from Pingdom.
type PingdomMetric struct {
	CheckID   *string `yaml:"checkId" validate:"required,notBlank,numeric" example:"1234567"`
	CheckType *string `yaml:"checkType" validate:"required,pingdomCheckTypeFieldValid" example:"uptime"`
	Status    *string `yaml:"status,omitempty" validate:"omitempty,pingdomStatusValid" example:"up,down"`
}

// GraphiteMetric represents metric from Graphite.
type GraphiteMetric struct {
	MetricPath *string `yaml:"metricPath" validate:"required,metricPathGraphite"`
}

// BigQueryMetric represents metric from BigQuery.
type BigQueryMetric struct {
	Query     string `yaml:"query" validate:"required,bigQueryRequiredColumns"`
	ProjectID string `yaml:"projectId" validate:"required"`
	Location  string `yaml:"location" validate:"required"`
}

// OpenTSDBMetric represents metric from OpenTSDB.
type OpenTSDBMetric struct {
	Query *string `yaml:"query" validate:"required"`
}

// GrafanaLokiMetric represents metric from GrafanaLokiMetric.
type GrafanaLokiMetric struct {
	Logql *string `yaml:"logql" validate:"required"`
}
