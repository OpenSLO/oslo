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
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	cobra.CheckErr(newRootCmd(version).Execute())
}

func newRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "oslo",
		Short:         "Oslo is a CLI tool for the OpenSLO spec",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
	}

	rootCmd.AddCommand(NewValidateCmd())
	rootCmd.AddCommand(NewFmtCmd())

	return rootCmd
}
