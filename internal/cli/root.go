package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "oslo",
		Short:         "Oslo is a CLI tool for the OpenSLO specification",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
	}

	coreGroup := &cobra.Group{
		ID:    "core",
		Title: "Core commands:",
	}
	rootCmd.AddGroup(coreGroup)

	subCommands := []*cobra.Command{
		NewValidateCmd(),
		NewFmtCmd(),
	}
	for _, subCmd := range subCommands {
		subCmd.GroupID = coreGroup.ID
		rootCmd.AddCommand(subCmd)
	}

	return rootCmd
}
