/*
Package convert provides a command to convert from openslo to other formats.

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

	"github.com/spf13/cobra"

	"github.com/OpenSLO/oslo/internal/convert"
	"github.com/OpenSLO/oslo/pkg/discoverfiles"
)

// NewConvertCmd returns a new command for formatting a file.
func NewConvertCmd() *cobra.Command {
	var passedFilePaths []string
	var recursive bool
	var format string
	var project string

	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Converts from OpenSLO to another format.",
		Long: `Converts from OpenSLO to another format.

Supported output formats are:
- nobl9

Multiple files can be converted by specifying them as arguments:

	oslo convert -f file1.yaml -f file2.yaml -o nobl9

You can also convert a directory of files:

  oslo convert -d path/to/directory -o nobl9

The output is written to standard output. If you want to write to a file, you can redirect the output:

  oslo convert -f file.yaml -o nobl9 > output.yaml
`,
		Args: cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			discoveredFilePaths, err := discoverfiles.DiscoverFilePaths(passedFilePaths, recursive)
			if err != nil {
				return err
			}
			discoveredFilePaths = convert.RemoveDuplicates(discoveredFilePaths)
			switch format {
			case "nobl9":
				if err := convert.Nobl9(cmd.OutOrStdout(), discoveredFilePaths, project); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}
			return nil
		},
	}

	discoverfiles.RegisterFileRelatedFlags(convertCmd, &passedFilePaths, &recursive)
	convertCmd.Flags().StringVarP(&format, "output", "o", "", "The output format to convert to.")
	if err := convertCmd.MarkFlagRequired("output"); err != nil {
		panic(err)
	}
	convertCmd.Flags().StringVarP(
		&project,
		"project",
		"p",
		"default",
		"Used for nobl9 output. What project to assign the resources to.",
	)

	return convertCmd
}
