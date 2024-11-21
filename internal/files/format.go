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
package files

import (
	"fmt"
	"io"

	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
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
