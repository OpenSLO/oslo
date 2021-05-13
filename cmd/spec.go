/*
Copyright © 2021 OpenSLO Team

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
package cmd

type Service struct {
	ObjectHeader `yaml:",inline"`
	Spec         ServiceSpec `yaml:"spec"`
}

// ServiceSpec represents content of Spec typical for Service Object
type ServiceSpec struct {
	Description string `yaml:"description" validate:"max=1050" example:"Bleeding edge web app"`
}

// Metadata represents part of object which is is common for all available Objects, for internal usage
type Metadata struct {
	Name        string `yaml:"name" validate:"required" example:"name"`
	DisplayName string `yaml:"displayName,omitempty" validate:"omitempty,min=0,max=63" example:"Prometheus Source"`
}

// MetadataHolder is an intermediate structure that can provides metadata related
// field to other structures
type MetadataHolder struct {
	Metadata Metadata `yaml:"metadata"`
}

// ObjectHeader represents Header which is common for all available Objects
type ObjectHeader struct {
	APIVersion     string `yaml:"apiVersion" validate:"required" example:"n9/v1alpha"`
	Kind           string `yaml:"kind" validate:"required" example:"kind"`
	MetadataHolder `yaml:",inline"`
}

// ObjectGeneric represents struct to which every Objects is parsable
// Specific types of Object have different structures as Spec
type ObjectGeneric struct {
	ObjectHeader `yaml:",inline"`
}

type sloSpec struct {
	APIVersion string `validate:"nonzero,regexp=^openslo\\/[a-zA-Z0-9]*$" yaml:"apiVersion"`
	Kind       string `validate:"nonzero,regexp=SLO"`
	Metadata   struct {
		Name        string `validate:"nonzero,max=63"`
		DisplayName string `validate:"regexp=^[a-zA-Z]*$"`
	}
	Spec struct {
		Description string `validate:"max=1050"`
		Service     string `validate:"nonzero,max=63"`
		Indicator   struct {
			ThresholdMetric struct {
				Source    string `validate:"regexp=^[a-zA-Z]*$"`
				QueryType string `validate:"regexp=^[a-zA-Z]*$"`
				Query     string // TODO what kind of validation is appropriate here?
			}
		}
		TimeWindows []struct{}
	}
}
