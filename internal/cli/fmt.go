/*
Package fmt handles formatting of the provided input.

# Copyright © 2022 OpenSLO Team

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
package cli

import (
	"fmt"

	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
	"github.com/spf13/cobra"

	"github.com/OpenSLO/oslo/internal/files"
)

// NewFmtCmd returns a new command for formatting a file.
func NewFmtCmd() *cobra.Command {
	var (
		passedFilePaths []string
		recursive       bool
		output          string
	)

	fmtCmd := &cobra.Command{
		Use:   "fmt",
		Short: "Formats the provided input into the standard format.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			discoveredFilePaths, err := files.Discover(passedFilePaths, recursive)
			if err != nil {
				return err
			}
			var format openslosdk.ObjectFormat
			switch output {
			case "json":
				format = openslosdk.FormatJSON
			case "yaml":
				format = openslosdk.FormatYAML
			default:
				return fmt.Errorf("invalid output format: %s", output)
			}
			return files.Format(cmd.OutOrStdout(), format, discoveredFilePaths)
		},
	}
	registerFileRelatedFlags(fmtCmd, &passedFilePaths, &recursive)
	fmtCmd.Flags().StringVarP(
		&output, "output", "o", "yaml",
		"The output format, one of [json, yaml].",
	)
	return fmtCmd
}
