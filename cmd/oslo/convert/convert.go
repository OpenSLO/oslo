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
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/OpenSLO/oslo/internal/pkg/convert"
)

// NewConvertCmd returns a new command for formatting a file.
func NewConvertCmd() *cobra.Command {
	var files []string
	var directory string

	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Converts from one format to another.",
		Run: func(cmd *cobra.Command, args []string) {
			// If a directory is provided, read all files in the directory.
			if directory != "" {
				dirFiles, err := os.ReadDir(directory)
				if err != nil {
					log.Fatal(err)
				}
				for _, file := range dirFiles {
					files = append(files, directory+"/"+file.Name())
				}
			}

			// Remove duplicates from the list of files so we are only
			// processing each file once.
			files = convert.RemoveDuplicates(files)

			// Convert the files
			if err := convert.ConvertFiles(cmd.OutOrStdout(), files); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
				os.Exit(1)
			}
		},
	}

	convertCmd.Flags().StringArrayVarP(&files, "file", "f", []string{}, "The file(s) to format.")
	convertCmd.Flags().StringVarP(&directory, "directory", "d", "", "The directory to format.")

	return convertCmd
}
