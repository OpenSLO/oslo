package files_test

import (
	"bytes"
	"path/filepath"
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
			name:    "invalid file",
			files:   []string{"v0alpha/invalid-file.yaml"},
			format:  openslosdk.FormatYAML,
			wantErr: true,
			wantOut: "",
		},
		{
			name:   "invalid content",
			files:  []string{"invalid-service.yaml"},
			format: openslosdk.FormatYAML,
			wantOut: `- apiVersion: openslo/v1alpha
  kind: Service
  metadata:
    displayName: My Rad Service
    name: my-rad service
  spec:
    description: This is a great description of an even better service.
`,
		},
		{
			name:   "formats single YAML file",
			files:  []string{"valid-service.yaml"},
			format: openslosdk.FormatYAML,
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
			name:   "formats single JSON file",
			files:  []string{"valid-service.yaml"},
			format: openslosdk.FormatJSON,
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
			name:   "formats multiple YAML files",
			files:  []string{"valid-service.yaml", "valid-service.yaml"},
			format: openslosdk.FormatYAML,
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
		{
			name:   "formats a list of services",
			files:  []string{"list-of-services.yaml"},
			format: openslosdk.FormatYAML,
			wantOut: `- apiVersion: openslo/v1
  kind: Service
  metadata:
    name: my-service-1
  spec: {}
- apiVersion: openslo/v1
  kind: Service
  metadata:
    name: my-service-2
  spec: {}
`,
		},
		{
			name:   "formats two documents",
			files:  []string{"two-documents.yaml"},
			format: openslosdk.FormatYAML,
			wantOut: `- apiVersion: openslo/v1
  kind: Service
  metadata:
    name: my-service-1
  spec: {}
- apiVersion: openslo/v1
  kind: Service
  metadata:
    name: my-service-2
  spec: {}
- apiVersion: openslo/v1
  kind: Service
  metadata:
    name: my-service-3
  spec: {}
`,
		},
		{
			name:   "formats two documents in JSON",
			files:  []string{"two-documents.yaml"},
			format: openslosdk.FormatJSON,
			wantOut: `[
  {
    "apiVersion": "openslo/v1",
    "kind": "Service",
    "metadata": {
      "name": "my-service-1"
    },
    "spec": {}
  },
  {
    "apiVersion": "openslo/v1",
    "kind": "Service",
    "metadata": {
      "name": "my-service-2"
    },
    "spec": {}
  },
  {
    "apiVersion": "openslo/v1",
    "kind": "Service",
    "metadata": {
      "name": "my-service-3"
    },
    "spec": {}
  }
]
`,
		},
		{
			name:   "formats single JSON doc to YAML",
			files:  []string{"valid-service.json"},
			format: openslosdk.FormatYAML,
			wantOut: `- apiVersion: openslo/v1alpha
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
			for i, file := range tc.files {
				tc.files[i] = filepath.Join("testdata", "format", file)
			}
			err := files.Format(out, tc.format, tc.files)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.wantOut, out.String())
		})
	}
}
