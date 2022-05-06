package convert

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConvertCmd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr bool
	}{
		{
			name:    "Single file",
			args:    []string{"-f", "../../../test/v1/service/service.yaml"},
			wantOut: "foo",
			wantErr: false,
		},
		{
			name: "Duplicate file",
			args: []string{
				"-f", "../../../test/v1/service/service.yaml",
				"-f", "../../../test/v1/service/service.yaml",
			},
			wantOut: "foo",
			wantErr: false,
		},
		{
			name: "Multiple files",
			args: []string{
				"-f", "../../../test/v1/data-source/data-source.yaml",
				"-f", "../../../test/v1/service/service.yaml",
			},
			wantOut: "foo",
			wantErr: false,
		},
		{
			name: "Directory",
			args: []string{
				"-d", "../../../test/v1/service",
			},
			wantOut: "foo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // Parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := new(bytes.Buffer)
			root := NewConvertCmd()
			root.SetOut(actual)
			root.SetErr(actual)
			root.SetArgs(tt.args)

			if err := root.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Error executing root command: %s", err)
				return
			}

			assert.Equal(t, tt.wantOut, actual.String())
		})
	}
}
