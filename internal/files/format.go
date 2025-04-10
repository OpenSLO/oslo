package files

import (
	"fmt"
	"io"

	"github.com/OpenSLO/go-sdk/pkg/openslosdk"
)

// Format formats multiple files and writes it to the provided writer, separated with "---".
func Format(out io.Writer, format openslosdk.ObjectFormat, sources []string) error {
	for i, src := range sources {
		if err := formatFile(out, format, src); err != nil {
			return err
		}
		if i != len(sources)-1 {
			if _, err := fmt.Fprintln(out, "---"); err != nil {
				return err
			}
		}
	}
	return nil
}

// formatFile formats a single formatFile and writes it to the provided writer.
func formatFile(out io.Writer, format openslosdk.ObjectFormat, source string) error {
	content, err := readRawSchema(source)
	if err != nil {
		return fmt.Errorf("issue reading content: %w", err)
	}
	objects, err := readObjectsFromRawData(content)
	if err != nil {
		return fmt.Errorf("issue parsing objects: %w", err)
	}

	return openslosdk.Encode(out, format, objects...)
}
