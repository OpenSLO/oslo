package cli

import (
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	if version == "" {
		version = "unknown"
	}
	cobra.CheckErr(newRootCmd(version).Execute())
}

func newRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "oslo",
		Short:         "Oslo is a CLI tool for the OpenSLO specification",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
	}

	rootCmd.AddCommand(NewValidateCmd())
	rootCmd.AddCommand(NewFmtCmd())

	return rootCmd
}
