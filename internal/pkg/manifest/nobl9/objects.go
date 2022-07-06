/*
Package manifest provides foundational structs.

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

package nobl9manifest

// ObjectHeader represents Header which is common for all available Objects.
type ObjectHeader struct {
	APIVersion string `yaml:"apiVersion" validate:"required" example:"openslo/v1alpha"`
}

// ObjectGeneric represents struct to which every Objects is parsable
// Specific types of Object have different structures as Spec.
type ObjectGeneric struct {
	ObjectHeader `yaml:",inline"`
}

// OpenSLOKind represents a type of object described by OpenSLO.
type OpenSLOKind interface {
	Kind() string
}
