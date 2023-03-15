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
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"

	nobl9manifest "github.com/OpenSLO/oslo/internal/pkg/manifest/nobl9"
	nobl9v1alpha "github.com/OpenSLO/oslo/internal/pkg/manifest/nobl9/v1alpha"
	"github.com/OpenSLO/oslo/internal/pkg/yamlutils"
	"github.com/OpenSLO/oslo/pkg/manifest"
	v1 "github.com/OpenSLO/oslo/pkg/manifest/v1"
)

// RemoveDuplicates to remove duplicate string from a slice.
func RemoveDuplicates(s []string) []string {
	result := make([]string, 0, len(s))
	m := make(map[string]struct{})

	for _, v := range s {
		if _, exists := m[v]; !exists {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

func getParsedObjects(filenames []string) (parsed []manifest.OpenSLOKind, err error) {
	for _, filename := range filenames {
		// Get the file contents.
		content, err := yamlutils.ReadConf(filename)
		if err != nil {
			return nil, fmt.Errorf("issue reading content: %w", err)
		}

		// Parse the byte arrays to OpenSLOKind objects.
		var p []manifest.OpenSLOKind
		p, err = yamlutils.Parse(content, filename)
		if err != nil {
			return nil, fmt.Errorf("issue parsing content: %w", err)
		}

		parsed = append(parsed, p...)
	}
	return parsed, nil
}

// getObjectsByKind function that that returns an object by Kind from a list of OpenSLOKinds.
func getObjectByKind(kind string, objects []manifest.OpenSLOKind) []manifest.OpenSLOKind {
	var found []manifest.OpenSLOKind
	for _, o := range objects {
		if o.Kind() == kind {
			found = append(found, o)
		}
	}
	return found
}

//------------------------------------------------------------------------------
//
//  Nobl9 Conversion
//

const (
	n9KindAnnotation = "nobl9.com/indicator-kind"

	kindAgent  string = "Agent"
	kindDirect string = "Direct"
)

/*
Nobl9 converts the provided file to Nobl9 yaml.

Nobl9 currently supports the following kinds:
- AlertMethod
- AlertPolicy
- Annotation
- DataExport
- Objective
- Project
- RoleBinding
- Service
- SLO

However OpenSLO doesn't support Annotation, DataExport, Project or RoleBinding so we can only convert to
AlertMethod, AlertPolicy, Objective, Service and SLO.
*/
func Nobl9(out io.Writer, filenames []string, project string) error {
	var rval []interface{}
	// These are used to track the names of the objects that
	// we have, in order to warn the user if they need to add
	// the objects to Nobl9 separately.
	var serviceNames []string
	var alertPolicyNames []string

	parsed, err := getParsedObjects(filenames)
	if err != nil {
		return fmt.Errorf("issue parsing content: %w", err)
	}

	// Get the service objects.
	if err = getN9ServiceObjects(parsed, &rval, &serviceNames, project); err != nil {
		return fmt.Errorf("issue getting service objects: %w", err)
	}

	// Get the alertPolicy objects.
	if err = getN9AlertPolicyObjects(parsed, &rval, &alertPolicyNames, project); err != nil {
		return fmt.Errorf("issue getting alertPolicy objects: %w", err)
	}

	// Get the SLO objects.
	if err = getN9SLObjects(parsed, &rval, serviceNames, alertPolicyNames, project); err != nil {
		return fmt.Errorf("issue getting SLO objects: %w", err)
	}

	// Print out all of our objects.
	for _, s := range rval {
		err = printYaml(out, s)
		if err != nil {
			return fmt.Errorf("issue printing content: %w", err)
		}
	}

	return nil
}

// Constructs Nobl9 SLO objects from our list of OpenSLOKinds.
func getN9SLObjects(
	parsed []manifest.OpenSLOKind,
	rval *[]interface{},
	serviceNames,
	alertPolicies []string,
	project string,
) error {
	objects := getObjectByKind("SLO", parsed)

	if len(objects) == 0 {
		return nil
	}

	// For each SLO object, create a Nobl9 SLO object.
	for _, slo := range objects {
		s := slo.(v1.SLO)

		// Check that the service name is in the list of service names, and warn the user if it isn't.
		if s.Spec.Service != "" && !stringInSlice(s.Spec.Service, serviceNames) {
			_ = printWarning(
				fmt.Sprintf(
					"Service %s is not in the list of services for SLO %s. "+
						"You will need verify that it is present in Nobl9 before applying.",
					s.Spec.Service,
					s.Metadata.DisplayName,
				),
			)
		}

		// Check that the alert policy name is in the list of alert policies, and warn the user if it isn't.
		for _, ap := range s.Spec.AlertPolicies {
			if !stringInSlice(ap, alertPolicies) {
				_ = printWarning(
					fmt.Sprintf(
						"Alert policy %s is not in the list of alert policies for SLO %s. "+
							"You will need verify that it is present in Nobl9 before applying.",
						ap,
						s.Metadata.DisplayName,
					),
				)
			}
		}

		tw, err := getN9TimeWindow(s.Spec.TimeWindow)
		if err != nil {
			return fmt.Errorf("issue getting time window: %w", err)
		}

		// Get the Objectives, aka Thresholds
		indicator := getN9SLISpec(s.Spec, parsed)
		thresholds, err := getN9Thresholds(s.Spec.Objectives, indicator)
		if err != nil {
			return fmt.Errorf("issue getting thresholds: %w", err)
		}

		indicatorMetadata := getN9IndicatorMetadata(s.Spec)
		n9Indicator := getN9Indicator(indicator, indicatorMetadata, project)

		*rval = append(*rval, nobl9v1alpha.SLO{
			ObjectHeader: getN9ObjectHeader("SLO", s.Metadata.Name, s.Metadata.DisplayName, project, s.Metadata.Labels),
			Spec: nobl9v1alpha.SLOSpec{
				Indicator:       n9Indicator,
				Description:     s.Spec.Description,
				BudgetingMethod: s.Spec.BudgetingMethod,
				Service:         s.Spec.Service,
				AlertPolicies:   s.Spec.AlertPolicies,
				TimeWindows:     tw,
				Thresholds:      thresholds,
			},
		})
	}

	return nil
}

func getN9MetricSourceName(msh v1.MetricSourceHolder) (string, bool) {
	name := msh.MetricSource.MetricSourceRef
	if name != "" {
		return name, true
	}

	return "", false
}

func getN9IndicatorMetadata(sloSpec v1.SLOSpec) (metadata v1.Metadata) {
	if sloSpec.Indicator != nil {
		return sloSpec.Indicator.Metadata
	}
	return metadata
}

// returns nobl9 indicator base on discovery and assumptions.
//
//nolint:gocognit,cyclop
func getN9Indicator(sliSpec v1.SLISpec, metadata v1.Metadata, project string) nobl9v1alpha.Indicator {
	// check to make sure we have a project, and use default if not
	metricSourceProject := "default"
	if project != "" {
		metricSourceProject = project
	}

	// check to make sure that we have an indicator
	//nolint:nestif
	if !reflect.ValueOf(sliSpec).IsZero() {
		var name string
		// check to see if we have a ThresholdMetric, and use that to set the MetricSource
		if !reflect.ValueOf(sliSpec.ThresholdMetric).IsZero() {
			if n, ok := getN9MetricSourceName(*sliSpec.ThresholdMetric); ok {
				name = n
			} else {
				_ = printWarning(
					"Threshold MetricSource was set but the MetricSourceRef was not, " +
						"so setting the name to Changeme for the MetricSource. Please update accordingly",
				)
				name = "Changeme"
			}
		}

		// check to see if we have a RawMetric and use that to set the MetricSource
		if !reflect.ValueOf(sliSpec.RatioMetric).IsZero() {
			// try all the possible MetricSourceHolders that a ratio metric might have
			if !reflect.ValueOf(sliSpec.RatioMetric.Good).IsZero() {
				if n, ok := getN9MetricSourceName(*sliSpec.RatioMetric.Good); ok {
					name = n
				}
			}
			if !reflect.ValueOf(sliSpec.RatioMetric.Bad).IsZero() {
				if n, ok := getN9MetricSourceName(*sliSpec.RatioMetric.Bad); ok {
					name = n
				}
			}
			if !reflect.ValueOf(sliSpec.RatioMetric.Total).IsZero() {
				if n, ok := getN9MetricSourceName(sliSpec.RatioMetric.Total); ok {
					name = n
				}
			}
		}
		if name != "" {
			return nobl9v1alpha.Indicator{
				MetricSource: nobl9v1alpha.MetricSourceSpec{
					Project: metricSourceProject,
					Name:    name,
					Kind:    getKindFromAnnotations(metadata),
				},
			}
		}
	}

	// Default return.  This handles any issues we might have found
	_ = printWarning(
		"No indicator found, or missing either a ThresholdMetric or RatioMetric, " +
			"so using a default Indicator and MetricSource.  Please update accordingly.",
	)
	return nobl9v1alpha.Indicator{
		MetricSource: nobl9v1alpha.MetricSourceSpec{
			Project: metricSourceProject,
			Name:    "Changeme",
			Kind:    getKindFromAnnotations(metadata),
		},
	}
}

func getKindFromAnnotations(metadata v1.Metadata) string {
	if value, ok := metadata.Annotations[n9KindAnnotation]; ok {
		switch strings.ToLower(value) {
		case "direct":
			return kindDirect
		case "agent":
			return kindAgent
		}
	}
	_ = printWarning(
		"We set as default MetricSource Kind: Agent if you want to change it to Direct use can use annotation " +
			n9KindAnnotation,
	)
	return kindAgent
}

// Return a list of nobl9v1alpha.Thresholds from a list of v1.Objectives.
func getN9Thresholds(o []v1.Objective, indicator v1.SLISpec) ([]nobl9v1alpha.Threshold, error) {
	var t []nobl9v1alpha.Threshold //nolint:prealloc
	for _, v := range o {
		v := v // local copy
		// if the operator isn't nil, then assign it, otherwise use default
		var operator *string
		defaultOp := "lt"
		if v.Op != "" {
			operator = &v.Op
		} else {
			operator = &defaultOp
		}

		// if v.Value is set use that, otherwise default to 1
		var value float64
		if v.Value != 0 {
			value = v.Value
		} else {
			value = 1
		}

		// start building our object
		th := nobl9v1alpha.Threshold{
			ThresholdBase: nobl9v1alpha.ThresholdBase{
				DisplayName: v.DisplayName,
				Value:       value,
			},
			Operator:     operator,
			BudgetTarget: &v.Target,
		}

		// Get CountMetrics if we have a ratioMetric
		if indicator.RatioMetric != nil {
			c, err := getN9CountMetrics(*indicator.RatioMetric)
			if err != nil {
				return nil, fmt.Errorf("issue getting count metrics: %w", err)
			}
			th.CountMetrics = &c
		}

		// Get thresholdMetrics
		if indicator.ThresholdMetric != nil {
			r, err := getN9RawMetrics(*indicator.ThresholdMetric)
			if err != nil {
				return nil, fmt.Errorf("issue getting raw metrics: %w", err)
			}
			th.RawMetric = &r
		}

		t = append(t, th)
	}
	return t, nil
}

func getN9RawMetrics(r v1.MetricSourceHolder) (nobl9v1alpha.RawMetricSpec, error) {
	raw, err := getN9MetricSource(r.MetricSource)
	if err != nil {
		return nobl9v1alpha.RawMetricSpec{}, fmt.Errorf("issue getting raw metric source: %w", err)
	}

	rm := nobl9v1alpha.RawMetricSpec{
		MetricQuery: &raw,
	}

	return rm, nil
}

func getN9CountMetrics(r v1.RatioMetric) (nobl9v1alpha.CountMetricsSpec, error) {
	// Error if Bad is not nil, since Nobl9 doesn't support it.
	if r.Bad != nil {
		return nobl9v1alpha.CountMetricsSpec{}, fmt.Errorf("metric spec Bad is not supported in Nobl9")
	}

	good, err := getN9MetricSource(r.Good.MetricSource)
	if err != nil {
		return nobl9v1alpha.CountMetricsSpec{}, fmt.Errorf("issue getting good metric source: %w", err)
	}

	total, err := getN9MetricSource(r.Total.MetricSource)
	if err != nil {
		return nobl9v1alpha.CountMetricsSpec{}, fmt.Errorf("issue getting total metric source: %w", err)
	}

	cm := nobl9v1alpha.CountMetricsSpec{
		Incremental: &r.Counter,
		GoodMetric:  &good,
		TotalMetric: &total,
	}
	return cm, nil
}

// Disabling the lint for this since theres not a really good way of doing this without a big switch statement.
//
//nolint:cyclop
func getN9MetricSource(m v1.MetricSource) (nobl9v1alpha.MetricSpec, error) {
	// Nobl9 supported metric sources.
	supportedMetricSources := map[string]string{
		"AmazonPrometheus":      "AmazonPrometheus",
		"AppDynamics":           "AppDynamics",
		"BigQuery":              "BigQuery",
		"CloudWatch":            "CloudWatch",
		"CloudWatchMetric":      "CloudWatchMetric",
		"Datadog":               "Datadog",
		"Dynatrace":             "Dynatrace",
		"Elasticsearch":         "Elasticsearch",
		"GoogleCloudMonitoring": "GoogleCloudMonitoring",
		"GrafanaLoki":           "GrafanaLoki",
		"Graphite":              "Graphite",
		"Instana":               "Instana",
		"Lightstep":             "Lightstep",
		"NewRelic":              "NewRelic",
		"OpenTSDB":              "OpenTSDB",
		"Pingdom":               "Pingdom",
		"Prometheus":            "Prometheus",
		"Redshift":              "Redshift",
		"Splunk":                "Splunk",
		"SplunkObservability":   "SplunkObservability",
		"SumoLogic":             "SumoLogic",
		"ThousandEyes":          "ThousandEyes",
	}
	var ms nobl9v1alpha.MetricSpec
	switch m.Type {
	case supportedMetricSources["Datadog"]:
		query := m.MetricSourceSpec["query"]
		ms = nobl9v1alpha.MetricSpec{
			Datadog: &nobl9v1alpha.DatadogMetric{
				Query: &query,
			},
		}
	case supportedMetricSources["Prometheus"]:
		query := m.MetricSourceSpec["promql"]
		ms = nobl9v1alpha.MetricSpec{
			Prometheus: &nobl9v1alpha.PrometheusMetric{
				PromQL: &query,
			},
		}
	case supportedMetricSources["AmazonPrometheus"]:
		query := m.MetricSourceSpec["promql"]
		ms = nobl9v1alpha.MetricSpec{
			AmazonPrometheus: &nobl9v1alpha.AmazonPrometheusMetric{
				PromQL: &query,
			},
		}
	case supportedMetricSources["NewRelic"]:
		query := m.MetricSourceSpec["nrql"]
		ms = nobl9v1alpha.MetricSpec{
			NewRelic: &nobl9v1alpha.NewRelicMetric{
				NRQL: &query,
			},
		}
	case supportedMetricSources["ThousandEyes"]:
		id := m.MetricSourceSpec["TestID"]
		testType := m.MetricSourceSpec["TestType"]

		// convert id (which is a string) to int64
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nobl9v1alpha.MetricSpec{}, fmt.Errorf("issue converting test id to int64: %w", err)
		}

		ms = nobl9v1alpha.MetricSpec{
			ThousandEyes: &nobl9v1alpha.ThousandEyesMetric{
				TestID:   &idInt,
				TestType: &testType,
			},
		}
	case supportedMetricSources["AppDynamics"]:
		appName := m.MetricSourceSpec["applicationName"]
		metricPath := m.MetricSourceSpec["metricPath"]
		ms = nobl9v1alpha.MetricSpec{
			AppDynamics: &nobl9v1alpha.AppDynamicsMetric{
				ApplicationName: &appName,
				MetricPath:      &metricPath,
			},
		}
	case supportedMetricSources["Splunk"]:
		query := m.MetricSourceSpec["query"]
		ms = nobl9v1alpha.MetricSpec{
			Splunk: &nobl9v1alpha.SplunkMetric{
				Query: &query,
			},
		}
	case supportedMetricSources["Lightstep"]:
		streamID := m.MetricSourceSpec["streamId"]
		typeOfData := m.MetricSourceSpec["typeOfData"]
		percentile := m.MetricSourceSpec["percentile"]

		// convert percentile (which is a string) to float64
		percentileFloat, err := strconv.ParseFloat(percentile, 64)
		if err != nil {
			return nobl9v1alpha.MetricSpec{}, fmt.Errorf("issue converting percentile to float64: %w", err)
		}

		ms = nobl9v1alpha.MetricSpec{
			Lightstep: &nobl9v1alpha.LightstepMetric{
				StreamID:   &streamID,
				TypeOfData: &typeOfData,
				Percentile: &percentileFloat,
			},
		}
	case supportedMetricSources["SplunkObservability"]:
		program := m.MetricSourceSpec["program"]
		ms = nobl9v1alpha.MetricSpec{
			SplunkObservability: &nobl9v1alpha.SplunkObservabilityMetric{
				Program: &program,
			},
		}
	case supportedMetricSources["Dynatrace"]:
		metricSelector := m.MetricSourceSpec["metricSelector"]
		ms = nobl9v1alpha.MetricSpec{
			Dynatrace: &nobl9v1alpha.DynatraceMetric{
				MetricSelector: &metricSelector,
			},
		}
	case supportedMetricSources["Elasticsearch"]:
		index := m.MetricSourceSpec["index"]
		query := m.MetricSourceSpec["query"]

		ms = nobl9v1alpha.MetricSpec{
			Elasticsearch: &nobl9v1alpha.ElasticsearchMetric{
				Index: &index,
				Query: &query,
			},
		}
	case supportedMetricSources["CloudWatch"]:
		namespace := m.MetricSourceSpec["namespace"]
		metricName := m.MetricSourceSpec["metricName"]
		region := m.MetricSourceSpec["region"]

		cwm := getN9CloudWatchQuery(m.MetricSourceSpec)
		cwm.Namespace = &namespace
		cwm.MetricName = &metricName
		cwm.Region = &region

		ms = nobl9v1alpha.MetricSpec{
			CloudWatch: &cwm,
		}
	case supportedMetricSources["Redshift"]:
		query := m.MetricSourceSpec["query"]
		region := m.MetricSourceSpec["region"]
		clusterID := m.MetricSourceSpec["clusterId"]
		databaseName := m.MetricSourceSpec["databaseName"]

		ms = nobl9v1alpha.MetricSpec{
			Redshift: &nobl9v1alpha.RedshiftMetric{
				Query:        &query,
				Region:       &region,
				ClusterID:    &clusterID,
				DatabaseName: &databaseName,
			},
		}
	case supportedMetricSources["SumoLogic"]:
		dataType := m.MetricSourceSpec["type"]
		query := m.MetricSourceSpec["query"]
		quantization := m.MetricSourceSpec["quantization"]
		rollup := m.MetricSourceSpec["rollup"]

		ms = nobl9v1alpha.MetricSpec{
			SumoLogic: &nobl9v1alpha.SumoLogicMetric{
				Type:         &dataType,
				Query:        &query,
				Quantization: &quantization,
				Rollup:       &rollup,
			},
		}
	case supportedMetricSources["Instana"]:
		metricType := m.MetricSourceSpec["metricType"]
		metricSource, err := unflatten(m.MetricSourceSpec)
		if err != nil {
			return nobl9v1alpha.MetricSpec{}, fmt.Errorf("issue converting Instana to map: %w", err)
		}
		infrastructure := metricSource["infrastructure"].(map[string]interface{})
		metricRetrievalMethod := infrastructure["metricRetrievalMethod"].(string)
		query := infrastructure["query"].(string)
		snapshotID := infrastructure["snapshotId"].(string)

		application := metricSource["application"].(map[string]interface{})
		groupBy := application["groupBy"].(map[string]interface{})
		tag := groupBy["tag"].(string)
		tagEntity := groupBy["tagEntity"].(string)
		tagSecondLevelKey := groupBy["tagSecondLevelKey"].(string)

		ms = nobl9v1alpha.MetricSpec{
			Instana: &nobl9v1alpha.InstanaMetric{
				MetricType: metricType,
				Infrastructure: &nobl9v1alpha.InstanaInfrastructureMetricType{
					MetricRetrievalMethod: metricRetrievalMethod,
					Query:                 &query,
					SnapshotID:            &snapshotID,
					MetricID:              infrastructure["metricId"].(string),
					PluginID:              infrastructure["pluginId"].(string),
				},
				Application: &nobl9v1alpha.InstanaApplicationMetricType{
					MetricID:    application["metricId"].(string),
					Aggregation: application["aggregation"].(string),
					APIQuery:    application["apiQuery"].(string),
					GroupBy: nobl9v1alpha.InstanaApplicationMetricGroupBy{
						Tag:               tag,
						TagEntity:         tagEntity,
						TagSecondLevelKey: &tagSecondLevelKey,
					},
				},
			},
		}
	case supportedMetricSources["Pingdom"]:
		checkID := m.MetricSourceSpec["checkId"]
		checkType := m.MetricSourceSpec["checkType"]
		status := m.MetricSourceSpec["status"]

		ms = nobl9v1alpha.MetricSpec{
			Pingdom: &nobl9v1alpha.PingdomMetric{
				CheckID:   &checkID,
				CheckType: &checkType,
				Status:    &status,
			},
		}
	case supportedMetricSources["Graphite"]:
		metricPath := m.MetricSourceSpec["metricPath"]

		ms = nobl9v1alpha.MetricSpec{
			Graphite: &nobl9v1alpha.GraphiteMetric{
				MetricPath: &metricPath,
			},
		}
	case supportedMetricSources["BigQuery"]:
		query := m.MetricSourceSpec["query"]
		projectID := m.MetricSourceSpec["projectId"]
		location := m.MetricSourceSpec["location"]

		ms = nobl9v1alpha.MetricSpec{
			BigQuery: &nobl9v1alpha.BigQueryMetric{
				Query:     query,
				ProjectID: projectID,
				Location:  location,
			},
		}
	case supportedMetricSources["OpenTSDB"]:
		query := m.MetricSourceSpec["query"]

		ms = nobl9v1alpha.MetricSpec{
			OpenTSDB: &nobl9v1alpha.OpenTSDBMetric{
				Query: &query,
			},
		}
	case supportedMetricSources["GrafanaLoki"]:
		logql := m.MetricSourceSpec["logql"]

		ms = nobl9v1alpha.MetricSpec{
			GrafanaLoki: &nobl9v1alpha.GrafanaLokiMetric{
				Logql: &logql,
			},
		}
	case supportedMetricSources["GoogleCloudMonitoring"]:
		query := m.MetricSourceSpec["query"]
		projectID := m.MetricSourceSpec["projectId"]
		ms = nobl9v1alpha.MetricSpec{
			GoogleCloudMonitoring: &nobl9v1alpha.GoogleCloudMonitoringMetric{
				Query:     &query,
				ProjectID: &projectID,
			},
		}
	default:
		// get the supportedMetricSources as a string
		var supportedMetricSourcesString string
		for k := range supportedMetricSources {
			supportedMetricSourcesString += k + ", "
		}

		return ms, fmt.Errorf(
			"unsupported metric source kind %s. Supported types are %s",
			m.Type,
			supportedMetricSourcesString,
		)
	}
	return ms, nil
}

func getN9CloudWatchQuery(m map[string]string) nobl9v1alpha.CloudWatchMetric {
	switch {
	case m["sql"] != "":
		val := m["sql"]
		return nobl9v1alpha.CloudWatchMetric{
			SQL: &val,
		}
	case m["json"] != "":
		val := m["json"]
		return nobl9v1alpha.CloudWatchMetric{
			JSON: &val,
		}
	case m["dimensions"] != "":
		val, _ := getN9CloudWatchDims(m["dimensions"])
		stat := m["stat"]
		return nobl9v1alpha.CloudWatchMetric{
			Stat:       &stat,
			Dimensions: val,
		}
	}

	return nobl9v1alpha.CloudWatchMetric{}
}

func getN9CloudWatchDims(dimensions string) ([]nobl9v1alpha.CloudWatchMetricDimension, error) {
	// split the incoming dimensions string into a CloudWatchMetricDimension
	dimsPieces := strings.Split(dimensions, ";")
	dims := make([]nobl9v1alpha.CloudWatchMetricDimension, 0, len(dimsPieces))
	for _, sequence := range dimsPieces {
		// for the cloudwatch dimensions, we expect them in a single set of kv pairs, with name and value as the two keys
		// example: 'name:foo,value:"foo";name:bar,value:"bar"'
		cwDim := nobl9v1alpha.CloudWatchMetricDimension{}
		for _, dimMap := range strings.Split(sequence, ",") {
			kv := strings.Split(dimMap, ":")
			if len(kv) != 2 {
				return []nobl9v1alpha.CloudWatchMetricDimension{}, fmt.Errorf("invalid dimension: %s", dimMap)
			}

			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])

			if strings.ToLower(key) == "name" {
				cwDim.Name = &val
			}

			if strings.ToLower(key) == "value" {
				cwDim.Value = &val
			}
		}
		dims = append(dims, cwDim)
	}

	return dims, nil
}

func getN9SLISpec(o v1.SLOSpec, parsed []manifest.OpenSLOKind) (s v1.SLISpec) {
	// if o.IndicatorRef is not nil, then return the indicator from the parsed list
	if o.IndicatorRef != nil {
		indicators := getObjectByKind("SLI", parsed)

		for _, i := range indicators {
			ind := i.(v1.SLI)
			if ind.Metadata.Name == *o.IndicatorRef {
				s = ind.Spec
				break
			}
		}
	} else {
		s = o.Indicator.Spec
	}
	return s
}

// Function that returns an nobl9v1alpha.TimeWindow from a OpenSLO TimeWindow.
func getN9TimeWindow(tw []v1.TimeWindow) ([]nobl9v1alpha.TimeWindow, error) {
	// return an error if the length of tw is greater than one, since we only support one TimeWindow
	if len(tw) > 1 {
		return nil, fmt.Errorf("OpenSLO only supports one TimeWindow")
	}

	if len(tw) < 1 {
		_ = printError("Nobl9 requires a TimeWindow defined.")
		return nil, fmt.Errorf("no TimeWindow found")
	}

	duration := tw[0].Duration

	unit, err := getDurationUnit(duration[len(duration)-1:])
	if err != nil {
		return nil, fmt.Errorf("issue getting duration unit: %w", err)
	}

	// Convert all but the last character of duration to an int.
	durationInt, err := strconv.Atoi(duration[:len(duration)-1])
	if err != nil {
		return nil, fmt.Errorf("issue converting duration to int: %w", err)
	}

	rval := []nobl9v1alpha.TimeWindow{
		{
			Unit:      unit,
			Count:     durationInt,
			IsRolling: tw[0].IsRolling,
		},
	}

	// only add Calendar for Calendar aligned
	if !tw[0].IsRolling {
		rval[0].Calendar = &nobl9v1alpha.Calendar{
			TimeZone:  tw[0].Calendar.TimeZone,
			StartTime: tw[0].Calendar.StartTime,
		}
	}

	return rval, nil
}

// Constructs Nobl9 AlertPolicy objects from our list of OpenSLOKinds.
func getN9AlertPolicyObjects(
	parsed []manifest.OpenSLOKind,
	rval *[]interface{},
	names *[]string,
	project string,
) error {
	// Get the alert policy object.
	ap := getObjectByKind("AlertPolicy", parsed)

	// Return if ap is empty.
	if len(ap) == 0 {
		return nil
	}

	// AlertCondition is required so get any from our parsed list.
	ac := getObjectByKind("AlertCondition", parsed)

	// For each AlertPolicy
	for _, o := range ap {
		// Cast to OpenSLO service objects.
		apObj, ok := o.(v1.AlertPolicy)
		if !ok {
			return fmt.Errorf("issue casting to AlertPolicy")
		}

		// Gather the alert conditions.
		var conditions []nobl9v1alpha.AlertCondition
		for _, a := range apObj.Spec.Conditions {
			err := getN9AlertCondition(a, &conditions, ac)
			if err != nil {
				return fmt.Errorf("issue getting alert condition: %w", err)
			}
		}

		// Construct the nobl9 AlertPolicy object from the OpenSLO AlertPolicy object.
		_ = printWarning("Using default serverity of 'Medium' in AlertPolicy, because we don't have an exact mapping")
		_ = printWarning("Using default CoolDownDuration in AlertPolicy, because OpenSLO doesn't support that")
		*rval = append(*rval, nobl9v1alpha.AlertPolicy{
			ObjectHeader: getN9ObjectHeader(
				"AlertPolicy",
				apObj.Metadata.Name,
				apObj.Metadata.DisplayName,
				project,
				apObj.Metadata.Labels,
			),
			Spec: nobl9v1alpha.AlertPolicySpec{
				Description:      apObj.Spec.Description,
				Conditions:       conditions,
				Severity:         "Medium",
				CoolDownDuration: "5m", // default
			},
		})

		// Add the name to our list of names.
		*names = append(*names, apObj.Metadata.Name)
	}
	return nil
}

// returns an nobl9v1alpha.AlertCondition from an OpenSLO.AlertPolicyCondition.
func getN9AlertCondition(
	apc v1.AlertPolicyCondition,
	conditions *[]nobl9v1alpha.AlertCondition,
	ac []manifest.OpenSLOKind,
) error {
	// If we have an inline condition, we can use it.
	//nolint: nestif
	if apc.AlertConditionInline != nil {
		_ = printWarning("using the default averageBurnRate in AlertCondition, since there isn't a direct match")
		_ = printWarning("Using default operator in AlertCondition, because OpenSLO doesn't support that feature")
		*conditions = append(*conditions, nobl9v1alpha.AlertCondition{
			Measurement:      "averageBurnRate",
			Value:            apc.AlertConditionInline.Spec.Condition.Threshold,
			LastsForDuration: apc.AlertConditionInline.Spec.Condition.AlertAfter,
			Operation:        "gt",
		})
	} else {
		// Error if we don't have any, since we need at least one.
		if len(ac) == 0 {
			return fmt.Errorf("no alert conditions found. Required for alert policy")
		}

		// If we don't have an inline condition, we need to get the AlertCondition.
		for _, c := range ac {
			// Get the AlertCondition that matches the name.
			acObj, err := c.(v1.AlertCondition)
			if !err {
				return fmt.Errorf("issue casting to AlertCondition")
			}
			if apc.AlertPolicyConditionSpec.ConditionRef == acObj.Metadata.Name {
				_ = printWarning("Using default averageBurnRate in AlertCondition, since there isn't direct mapping")
				_ = printWarning("Using default operator in AlertCondition, because OpenSLO doesn't support that feature")
				*conditions = append(*conditions, nobl9v1alpha.AlertCondition{
					Measurement:      "averageBurnRate",
					Value:            acObj.Spec.Condition.Threshold,
					LastsForDuration: acObj.Spec.Condition.AlertAfter,
					Operation:        "gt",
				})
			} else {
				return fmt.Errorf("alert condition %s not found", apc.AlertPolicyConditionSpec.ConditionRef)
			}
		}
	}
	return nil
}

// function that takes a manifest.OpenSLOKind and returns a nobl9v1alpha.ObjectHeader.
func getN9ObjectHeader(kind, name, displayName, project string, labels v1.Labels) nobl9v1alpha.ObjectHeader {
	return nobl9v1alpha.ObjectHeader{
		ObjectHeader: nobl9manifest.ObjectHeader{
			APIVersion: nobl9v1alpha.APIVersion,
		},
		Kind: kind,
		MetadataHolder: nobl9v1alpha.MetadataHolder{
			Metadata: nobl9v1alpha.Metadata{
				Name:        name,
				DisplayName: displayName,
				Project:     project,
				Labels:      getN9Labels(labels),
			},
		},
	}
}

// getN9Labels takes a v1.Labels object and maps it to a nobl9v1alpha.Labels.
func getN9Labels(labels v1.Labels) map[string][]string {
	if labels == nil {
		return nil
	}
	rval := make(map[string][]string)
	for k, v := range labels {
		rval[k] = v
	}
	return rval
}

// Constructs Nobl9 Service objects from our list of OpenSLOKinds.
func getN9ServiceObjects(parsed []manifest.OpenSLOKind, rval *[]interface{}, names *[]string, project string) error {
	// Get the service object.
	obj := getObjectByKind("Service", parsed)

	for _, o := range obj {
		// Cast to OpenSLO service objects.
		srvObj, ok := o.(v1.Service)
		if !ok {
			return fmt.Errorf("issue casting to service")
		}
		// Construct the nobl9 service object from the OpenSLO service object.
		*rval = append(*rval, nobl9v1alpha.Service{
			ObjectHeader: getN9ObjectHeader(
				"Service",
				srvObj.Metadata.Name,
				srvObj.Metadata.DisplayName,
				project,
				srvObj.Metadata.Labels,
			),
			Spec: nobl9v1alpha.ServiceSpec{
				Description: srvObj.Spec.Description,
			},
		})

		// Add the name to the list of names.
		*names = append(*names, srvObj.Metadata.Name)
	}
	return nil
}

// ------------------------------------------------------------------------------
//
// Helper functions.
func printYaml(out io.Writer, object interface{}) error {
	// Convert parsed to yaml and print to out.
	yml, err := yaml.Marshal(object)
	if err != nil {
		return fmt.Errorf("issue marshaling content: %w", err)
	}

	fmt.Fprint(out, "---\n")
	_, err = out.Write(yml)
	if err != nil {
		return fmt.Errorf("issue writing content: %w", err)
	}

	return nil
}

// Function to print warning messages to Stderr so that we can see them when
// when doing redirection in the console.
func printWarning(message string) error {
	yellow := color.New(color.FgYellow).Add(color.Bold)
	white := color.New(color.FgWhite).Add(color.Bold)

	yellow.EnableColor()
	white.EnableColor()

	if _, err := yellow.Fprint(os.Stderr, "WARNING: "); err != nil {
		return fmt.Errorf("issue printing warning: %w", err)
	}

	if _, err := white.Fprintln(os.Stderr, message); err != nil {
		return fmt.Errorf("issue printing warning: %w", err)
	}

	yellow.DisableColor()
	white.DisableColor()

	color.Unset()
	return nil
}

func printError(message string) error {
	red := color.New(color.FgRed).Add(color.Bold)
	white := color.New(color.FgWhite).Add(color.Bold)

	red.EnableColor()
	white.EnableColor()

	if _, err := red.Fprint(os.Stderr, "ERROR: "); err != nil {
		return fmt.Errorf("issue printing warning: %w", err)
	}

	if _, err := white.Fprintln(os.Stderr, message); err != nil {
		return fmt.Errorf("issue printing warning: %w", err)
	}

	red.DisableColor()
	white.DisableColor()

	color.Unset()
	return nil
}

// Function that unflattens json into nested maps.
func unflatten(json map[string]string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, value := range json {
		keyParts := strings.Split(key, ".")
		m := result
		for i, k := range keyParts[:len(keyParts)-1] {
			v, exists := m[k]
			if !exists {
				newMap := map[string]interface{}{}
				m[k] = newMap
				m = newMap
				continue
			}

			innerMap, ok := v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("key=%v is not an object", strings.Join(keyParts[0:i+1], "."))
			}
			m = innerMap
		}

		leafKey := keyParts[len(keyParts)-1]
		if _, exists := m[leafKey]; exists {
			return nil, fmt.Errorf("key=%v already exists", key)
		}
		m[keyParts[len(keyParts)-1]] = value
	}

	return result, nil
}

// Function that takes a duration shorthand string and returns the unit of time, eg minute, hour, month.
func getDurationUnit(d string) (string, error) {
	switch d {
	case "m":
		return "Minute", nil
	case "h":
		return "Hour", nil
	case "d":
		return "Day", nil
	case "w":
		return "Week", nil
	case "M":
		return "Month", nil
	case "Q":
		return "Quarter", nil
	case "Y":
		return "Year", nil
	}

	return "", fmt.Errorf("duration unit not supported %s", d)
}

// Checks that the given string is in the given slice.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
