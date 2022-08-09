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
package fmt

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	fmtCmd "github.com/OpenSLO/oslo/internal/pkg/fmt"
)

// NewFmtCmd returns a new command for formatting a file.
func NewFmtCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fmt",
		Short: "Formats the provided input into the standard format.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := fmtCmd.File(cmd.OutOrStdout(), args[0]); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
				os.Exit(1)
			}
		},
	}
}
