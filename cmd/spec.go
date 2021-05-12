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

type newUserRequest struct {
	Conf struct {
		Username string `validate:"nonzero,min=3,max=40,regexp=^[a-zA-Z]*$"`
		Name     string `validate:"nonzero"`
		Age      int    `validate:"min=21"`
		Password string `validate:"nonzero,min=8"`
	}
}

type serviceSpec struct {
	APIVersion string `validate:"nonzero,regexp=^openslo\\/[a-zA-Z0-9]*$" yaml:"apiVersion"`
	Kind       string `validate:"nonzero,regexp=Service"`
	Metadata   struct {
		Name        string `validate:"nonzero,max=63"`
		DisplayName string `validate:"regexp=^[a-zA-Z]*$"`
	}
	Spec struct {
		Description string `validate:"max=1050"`
	}
}

type sloSpec struct{}
