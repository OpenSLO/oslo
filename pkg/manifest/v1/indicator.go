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
	"github.com/OpenSLO/oslo/pkg/manifest"
)

// DataSource defines the data source for the SLI.
type DataSource struct {
	manifest.ObjectHeader `yaml:",inline"`
	Spec                  DataSourceSpec `yaml:"spec" validate:"required"`
}

// DataSourceSpec defines the data source specification.
type DataSourceSpec struct {
	Type              string            `yaml:"type" validate:"required"`
	ConnectionDetails map[string]string `yaml:"connectionDetails"`
}

// SLI represents the SLI.
type SLI struct {
	manifest.ObjectHeader `yaml:",inline"`
	Metadata              Metadata `yaml:"metadata" validate:"required"`
	Spec                  SLISpec  `yaml:"spec" validate:"required"`
}

// SliInline represents the SLI inline.
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

// Kind returns the name of this type.
func (DataSource) Kind() string {
	return "DataSource"
}

// Kind returns the name of this type.
func (SLI) Kind() string {
	return "SLI"
}
