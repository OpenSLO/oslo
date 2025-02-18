package cli

import "github.com/spf13/cobra"

// registerFileRelatedFlags registers flags --file | -f and --recursive | -R for command
// passed as the argument and make them required.
func registerFileRelatedFlags(cmd *cobra.Command, filePaths *[]string, recursive *bool) {
	const fileFlag = "file"
	cmd.Flags().StringArrayVarP(
		filePaths, fileFlag, "f", []string{},
		"The file(s) that contain the configurations.",
	)
	if err := cmd.MarkFlagRequired(fileFlag); err != nil {
		panic(err)
	}
	cmd.Flags().BoolVarP(
		recursive, "recursive", "R", false,
		"Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.", //nolint:lll
	)
}
