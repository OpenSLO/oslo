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
package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
  -h, --help   help for validate
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
  -h, --help   help for fmt
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
