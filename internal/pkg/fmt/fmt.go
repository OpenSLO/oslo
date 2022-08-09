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
	"io"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/oslo/internal/pkg/yamlutils"
)

// File formats a single file and writes it to the provided writer.
func File(out io.Writer, source string) error {
	// Get the file contents.
	content, err := yamlutils.ReadConf(source)
	if err != nil {
		return fmt.Errorf("issue reading content: %w", err)
	}

	// Parse the byte arrays to OpenSLOKind objects.
	parsed, err := yamlutils.Parse(content, source)
	if err != nil {
		return fmt.Errorf("issue parsing content: %w", err)
	}

	// New encoder that will write to the provided writer.
	enc := yaml.NewEncoder(out)
	enc.SetIndent(2)

	for _, o := range parsed {
		// Encode the object to YAML.
		if err := enc.Encode(o); err != nil {
			return err
		}
	}
	return nil
}
