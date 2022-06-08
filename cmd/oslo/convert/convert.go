/*
Package convert provides a command to convert from openslo to other formats.

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
package convert

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/OpenSLO/oslo/internal/pkg/convert"
)

// NewConvertCmd returns a new command for formatting a file.
func NewConvertCmd() *cobra.Command {
	var files []string
	var directory string
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

The output is written to standard output.  If you want to write to a file, you can redirect the output:

  oslo convert -f file.yaml -o nobl9 > output.yaml
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If a directory is provided, read all files in the directory.
			if directory != "" {
				dirFiles, err := os.ReadDir(directory)
				if err != nil {
					return err
				}
				for _, file := range dirFiles {
					files = append(files, directory+"/"+file.Name())
				}
			}

			// Remove duplicates from the list of files so we are only
			// processing each file once.
			files = convert.RemoveDuplicates(files)

			// Convert the files to the specified format.
			switch format {
			case "nobl9":
				if err := convert.Nobl9(cmd.OutOrStdout(), files, project); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}

			return nil
		},
	}

	convertCmd.Flags().StringArrayVarP(&files, "file", "f", []string{}, "The file(s) to format.")
	convertCmd.Flags().StringVarP(&directory, "directory", "d", "", "The directory to format.")
	convertCmd.Flags().StringVarP(&format, "output", "o", "", "The output format to convert to.")
	convertCmd.Flags().StringVarP(&project, "project", "p", "default", "Used for nobl9 output. What project to assign the resources to.")

	convertCmd.MarkFlagRequired("format")

	return convertCmd
}
