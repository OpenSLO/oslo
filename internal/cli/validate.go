package cli

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/OpenSLO/go-sdk/pkg/openslosdk"

	"github.com/OpenSLO/oslo/internal/files"
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
			discoveredFilePaths, err := files.Discover(passedFilePaths, recursive)
			if err != nil {
				return err
			}
			objectsPerSource, err := files.ReadObjects(discoveredFilePaths)
			if err != nil {
				return err
			}
			sources := slices.Sorted(maps.Keys(objectsPerSource))

			hasErrors := false
			for _, src := range sources {
				objects := objectsPerSource[src]
				switch len(objects) {
				case 1:
					if err = objects[0].Validate(); err != nil {
						hasErrors = true
						printStderr(fmt.Errorf("Errors in %s:\n%s", src, indentString(err.Error(), 2)))
					}
				default:
					if err = openslosdk.Validate(objects...); err != nil {
						hasErrors = true
						printStderr(fmt.Errorf("Errors in %s:\n%s", src, indentString(err.Error(), 2)))
					}
				}
			}
			if !hasErrors {
				printStderr("Valid!")
				return nil
			}
			return errors.New("Configuration is not valid!")
		},
	}
	registerFileRelatedFlags(validateCmd, &passedFilePaths, &recursive)
	return validateCmd
}

func indentString(s string, i int) string {
	indent := strings.Repeat(" ", i)
	split := strings.Split(s, "\n")
	for i := range split {
		split[i] = indent + split[i]
	}
	return strings.Join(split, "\n")
}

func printStderr(err any) {
	fmt.Fprintln(os.Stderr, err)
}
