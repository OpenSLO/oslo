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
package files_test

import (
	"bytes"
	"testing"

	"github.com/OpenSLO/go-sdk/pkg/openslosdk"
	"github.com/stretchr/testify/assert"

	"github.com/OpenSLO/oslo/internal/files"
)

func TestFormatFiles(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		files   []string
		format  openslosdk.ObjectFormat
		wantOut string
		wantErr bool
	}{
		{
			name:    "Invalid file",
			files:   []string{"test/v1alpha/invalid-file.yaml"},
			format:  openslosdk.FormatYAML,
			wantErr: true,
			wantOut: "",
		},
		{
			name:    "Invalid content",
			files:   []string{"../../test/v1alpha/invalid-service.yaml"},
			format:  openslosdk.FormatYAML,
			wantErr: true,
			wantOut: "",
		},
		{
			name:    "Passes single file",
			files:   []string{"../../test/v1alpha/valid-service.yaml"},
			format:  openslosdk.FormatYAML,
			wantErr: false,
			wantOut: `- apiVersion: openslo/v1alpha
  kind: Service
  metadata:
    displayName: My Rad Service
    name: my-rad-service
  spec:
    description: This is a great description of an even better service.
`,
		},
		{
			name:    "Passes single JSON file",
			files:   []string{"../../test/v1alpha/valid-service.yaml"},
			format:  openslosdk.FormatJSON,
			wantErr: false,
			wantOut: `[
  {
    "apiVersion": "openslo/v1alpha",
    "kind": "Service",
    "metadata": {
      "name": "my-rad-service",
      "displayName": "My Rad Service"
    },
    "spec": {
      "description": "This is a great description of an even better service."
    }
  }
]
`,
		},
		{
			name:    "Passes multiple files",
			files:   []string{"../../test/v1alpha/valid-service.yaml", "../../test/v1alpha/valid-service.yaml"},
			format:  openslosdk.FormatYAML,
			wantErr: false,
			wantOut: `- apiVersion: openslo/v1alpha
  kind: Service
  metadata:
    displayName: My Rad Service
    name: my-rad-service
  spec:
    description: This is a great description of an even better service.
---
- apiVersion: openslo/v1alpha
  kind: Service
  metadata:
    displayName: My Rad Service
    name: my-rad-service
  spec:
    description: This is a great description of an even better service.
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out := &bytes.Buffer{}
			if err := files.Format(out, tc.format, tc.files); (err != nil) != tc.wantErr {
				t.Errorf("fmtFile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.wantOut, out.String())
		})
	}
}
