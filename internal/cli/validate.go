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
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/OpenSLO/oslo/pkg/discoverfiles"
	"github.com/OpenSLO/oslo/pkg/validate"
)

// NewValidateCmd returns a new cobra.Command for the validate command.
func NewValidateCmd() *cobra.Command {
	var passedFilePaths []string
	var recursive bool

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates your yaml file against the OpenSLO spec.",
		Long:  `Validates your yaml file against the OpenSLO spec.`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			discoveredFilePaths, err := discoverfiles.DiscoverFilePaths(passedFilePaths, recursive)
			if err != nil {
				return err
			}
			if err := validate.Files(discoveredFilePaths); err != nil {
				return err
			}
			fmt.Println("Valid!")
			return nil
		},
	}
	discoverfiles.RegisterFileRelatedFlags(validateCmd, &passedFilePaths, &recursive)
	return validateCmd
}
