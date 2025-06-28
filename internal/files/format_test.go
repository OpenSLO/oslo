package files_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/OpenSLO/go-sdk/pkg/openslosdk"
	"github.com/stretchr/testify/assert"

	"github.com/OpenSLO/oslo/internal/files"
	"github.com/OpenSLO/oslo/internal/pathutils"
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
			files:   []string{"v0alpha/invalid-file.yaml"},
			format:  openslosdk.FormatYAML,
			wantErr: true,
			wantOut: "",
		},
		{
			name:    "Invalid content",
			files:   []string{"v1alpha/invalid-service.yaml"},
			format:  openslosdk.FormatYAML,
			wantErr: false,
			wantOut: `- apiVersion: openslo/v1alpha
  kind: Service
  metadata:
    displayName: My Rad Service
    name: my-rad service
  spec:
    description: This is a great description of an even better service.
- apiVersion: openslo/v1alpha
  kind: Service
  metadata:
    name: this
  spec: {}
`,
		},
		{
			name:    "Passes single file",
			files:   []string{"v1alpha/valid-service.yaml"},
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
			files:   []string{"v1alpha/valid-service.yaml"},
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
			files:   []string{"v1alpha/valid-service.yaml", "v1alpha/valid-service.yaml"},
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
			for i, file := range tc.files {
				tc.files[i] = filepath.Join(pathutils.FindModuleRoot(), "test", "inputs", file)
			}
			if err := files.Format(out, tc.format, tc.files); (err != nil) != tc.wantErr {
				t.Errorf("fmtFile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.wantOut, out.String())
		})
	}
}
