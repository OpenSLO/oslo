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
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// DataSource defines the data source for the SLI.
type DataSource struct {
	ObjectHeader `yaml:",inline"`
	Spec         DataSourceSpec `yaml:"spec" validate:"required"`
}

// DataSourceSpec defines the data source specification.
type DataSourceSpec struct {
	Type              string            `yaml:"type" validate:"required"`
	ConnectionDetails map[string]string `yaml:"connectionDetails"`
}

// SLI represents the SLI.
type SLI struct {
	ObjectHeader `yaml:",inline"`
	Spec         SLISpec `yaml:"spec" validate:"required"`
}

// SLIInline represents the SLI inline.
type SLIInline struct {
	Metadata Metadata `yaml:"metadata" validate:"required"`
	Spec     SLISpec  `yaml:"spec" validate:"required"`
}

// SLISpec defines the SLI specification.
type SLISpec struct {
	ThresholdMetric *MetricSourceHolder `yaml:"thresholdMetric,omitempty" validate:"required_without=RatioMetric"`
	RatioMetric     *RatioMetric        `yaml:"ratioMetric,omitempty" validate:"required_without=ThresholdMetric"`
}

// RatioMetric represents the ratio metric.
type RatioMetric struct {
	Counter bool                `yaml:"counter" example:"true"`
	Good    *MetricSourceHolder `yaml:"good,omitempty" validate:"required_without=Bad"`
	Bad     *MetricSourceHolder `yaml:"bad,omitempty" validate:"required_without=Good"`
	Total   MetricSourceHolder  `yaml:"total" validate:"required"`
}

// MetricSourceHolder represents the metric source holder.
type MetricSourceHolder struct {
	MetricSource MetricSource `yaml:"metricSource" validate:"required"`
}

// MetricSource represents the metric source.
type MetricSource struct {
	MetricSourceRef  string            `yaml:"metricSourceRef,omitempty" validate:"required_without=MetricSourceSpec"`
	Type             string            `yaml:"type,omitempty" validate:"required_without=MetricSourceRef"`
	MetricSourceSpec map[string]string `yaml:"spec" validate:"required_without=MetricSourceRef"`
}

// UnmarshalYAML is used to override the default unmarshal behavior.
// MetricSources can have varying structures, so this method performs the following tasks:
//  1. Extracts the MetricSourceRef and Type separately, and assigns them to the MetricSource.
//  2. Attempts to unmarshal the MetricSourceSpec, which can be either a string, a list, or a more complex structure:
//     2a. If it's a scalar string, it is added directly as a single string.
//     2b. If it's a list of scalar values, they are concatenated into a single semicolon separated string.
//     2c. If it's a more complex sequence, the values are processed and flattened appropriately.
//
// This also assumes a certain flat structure that we can revisit if the need arises.

func (m *MetricSource) UnmarshalYAML(value *yaml.Node) error {
	// temp struct to unmarshal the string values
	var tmpMetricSource struct {
		MetricSourceRef  string               `yaml:"metricSourceRef,omitempty" validate:"required_without=MetricSourceSpec"`
		Type             string               `yaml:"type,omitempty" validate:"required_without=MetricSourceRef"`
		MetricSourceSpec map[string]yaml.Node `yaml:"spec"`
	}

	if err := value.Decode(&tmpMetricSource); err != nil {
		return err
	}

	// no error with these, assign them
	m.MetricSourceRef = tmpMetricSource.MetricSourceRef
	m.Type = tmpMetricSource.Type
	// initialize this so we can assign the values later
	m.MetricSourceSpec = make(map[string]string)

	for k, v := range tmpMetricSource.MetricSourceSpec {
		// simple use case
		if v.Kind == yaml.ScalarNode {
			m.MetricSourceSpec[k] = v.Value
		}

		if v.Kind == yaml.SequenceNode {
			// top level string that we will join with a semicolon
			seqStrings := []string{}
			for _, node := range v.Content {
				if node.Kind == yaml.ScalarNode {
					seqStrings = append(seqStrings, node.Value)
				} else if node.Kind == yaml.MappingNode {
					// each of these are k/v pairs that we will join with a comma
					kvPairs := []string{}
					for i := 0; i < len(node.Content); i += 2 {
						kvPairs = append(kvPairs, fmt.Sprintf("%s:%s", node.Content[i].Value, node.Content[i+1].Value))
					}
					seqStrings = append(seqStrings, strings.Join(kvPairs, ","))
				}
			}
			m.MetricSourceSpec[k] = strings.Join(seqStrings, ";")
		}
	}

	return nil
}

// Kind returns the name of this type.
func (DataSource) Kind() string {
	return "DataSource"
}

// Kind returns the name of this type.
func (SLI) Kind() string {
	return "SLI"
}
