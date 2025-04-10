package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:lll
func Test_newRootCmd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr bool
	}{
		{
			name:    "Command does not exist",
			args:    []string{"doesnotexist"},
			wantOut: "",
			wantErr: true,
		},
		{
			name: "validate command exists",
			args: []string{"validate", "--help"},
			wantOut: `Validates your yaml file against the OpenSLO spec.

Usage:
  oslo validate [flags]

Flags:
  -f, --file stringArray   The file(s) that contain the configurations.
  -h, --help               help for validate
  -R, --recursive          Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
`,
			wantErr: false,
		},
		{
			name: "fmt command exists",
			args: []string{"fmt", "--help"},
			wantOut: `Formats the provided input into the standard format.

Usage:
  oslo fmt [flags]

Flags:
  -f, --file stringArray   The file(s) that contain the configurations.
  -h, --help               help for fmt
  -o, --output string      The output format, one of [json, yaml]. (default "yaml")
  -R, --recursive          Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // Parallel testing
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := new(bytes.Buffer)
			root := newRootCmd("testVersion")
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
