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

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/oslo/pkg/manifest"
)

// APIVersion is a value of valid apiVersions.
const (
	APIVersion = "openslo/v1"
)

// Possible values of field kind for valid Objects.
const (
	KindAlertCondition          = "AlertCondition"
	KindAlertNotificationTarget = "AlertNotificationTarget"
	KindAlertPolicy             = "AlertPolicy"
	KindDataSource              = "DataSource"
	KindSLI                     = "SLI"
	KindSLO                     = "SLO"
	KindService                 = "Service"
)

// Parse is responsible for parsing all structs in this apiVersion.
func Parse(fileContent []byte, m ObjectGeneric, filename, kind string) (manifest.OpenSLOKind, error) {
	switch kind {
	case KindAlertCondition:
		var content AlertCondition
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case KindAlertNotificationTarget:
		var content AlertNotificationTarget
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case KindAlertPolicy:
		var content AlertPolicy
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case KindDataSource:
		var content DataSource
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case KindService:
		var content Service
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case KindSLO:
		var content SLO
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case KindSLI:
		var content SLI
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	default:
		return nil, fmt.Errorf("unsupported kind: %s", m.Kind)
	}
}

// ----------------------------------------------------------------------------
// Object definitions
// ----------------------------------------------------------------------------

// Labels is a map of labels.
type Labels map[string]string

// Annotations is a map of annotations.
type Annotations map[string]string

// Metadata represents part of object which is is common for all available Objects, for internal usage.
type Metadata struct {
	Name        string      `yaml:"name" validate:"required" example:"name"`
	DisplayName string      `yaml:"displayName,omitempty" validate:"omitempty,min=0,max=63" example:"Prometheus Source"`
	Labels      Labels      `json:"labels,omitempty" validate:"omitempty"`
	Annotations Annotations `json:"annotations,omitempty" validate:"omitempty"`
}

// MetadataHolder is an intermediate structure that can provides metadata related
// field to other structures.
type MetadataHolder struct {
	Metadata Metadata `yaml:"metadata"`
}

// ObjectHeader represents Header which is common for all available Objects.
type ObjectHeader struct {
	manifest.ObjectHeader `yaml:",inline"`
	Kind                  string `yaml:"kind" validate:"required,oneof=Service SLO SLI AlertPolicy AlertNotificationTarget AlertCondition DataSource" example:"kind"` //nolint:lll
	MetadataHolder        `yaml:",inline"`
}

// ObjectGeneric represents struct to which every Objects is parsable
// Specific types of Object have different structures as Spec.
type ObjectGeneric struct {
	ObjectHeader `yaml:",inline"`
}

// MetricSourceSpec represents the metric source.
type MetricSourceSpec struct {
	Source    string `yaml:"source" validate:"required,alpha"`
	QueryType string `yaml:"queryType" validate:"required,alpha"`
	Query     string `yaml:"query" validate:"required"`
}

// ObjectiveBase base structure representing a threshold.
type ObjectiveBase struct {
	DisplayName string  `yaml:"displayName" validate:"max=1050" example:"Good"`
	Value       float64 `yaml:"value" validate:"numeric" example:"100"`
}

/*----- Service -----*/

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
