/*
Package yamlutils provides functions to parse OpenSLO manifests.

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
package yamlutils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/oslo/pkg/manifest"
	v1 "github.com/OpenSLO/oslo/pkg/manifest/v1"
	"github.com/OpenSLO/oslo/pkg/manifest/v1alpha"
)

// ReadConf reads in filename for a yaml file and returns the byte array.
func ReadConf(filename string) ([]byte, error) {
	if filename == "-" {
		return io.ReadAll(os.Stdin)
	}
	fileContent, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

// Parse takes the provided byte array, parses it, and returns an array of parsed struts.
// Ignoring the complexity linting errors for now, until we can figure
// out how to handle the complexity better.
func Parse(fileContent []byte, filename string) ( //nolint: gocognit, cyclop
	parsedStructs []manifest.OpenSLOKind,
	err error,
) {
	var m manifest.ObjectGeneric
	// unmarshal here to get the APIVersion so we can process the file correctly
	if err = yaml.Unmarshal(fileContent, &m); err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	var allErrors error
	switch m.APIVersion {
	// This is where we add new versions of the OpenSLO spec.
	case v1alpha.APIVersion:
		// unmarshal again to get the v1alpha struct
		var o v1alpha.ObjectGeneric
		if err := yaml.Unmarshal(fileContent, &o); err != nil {
			return nil, fmt.Errorf("in file %q: %w", filename, err)
		}

		// loop through and get all of the documents in the file
		decoder := yaml.NewDecoder(bytes.NewReader(fileContent))
		for {
			var i interface{}
			err := decoder.Decode(&i)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("in file %q: %w", filename, err)
			}
			c, err := yaml.Marshal(&i)
			if err != nil {
				return nil, fmt.Errorf("in file %q: %w", filename, err)
			}

			content, e := v1alpha.Parse(c, o, filename)
			if e != nil {
				allErrors = multierror.Append(allErrors, e)
			}
			parsedStructs = append(parsedStructs, content)
		}
	case v1.APIVersion:
		// unmarshal again to get the v1 struct
		var o v1.ObjectGeneric
		if err := yaml.Unmarshal(fileContent, &o); err != nil {
			return nil, fmt.Errorf("in file %q: %w", filename, err)
		}
		// loop through and get all of the documents in the file
		decoder := yaml.NewDecoder(bytes.NewReader(fileContent))
		for {
			var i interface{}
			err := decoder.Decode(&i)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("in file %q: %w", filename, err)
			}
			c, err := yaml.Marshal(&i)
			if err != nil {
				return nil, fmt.Errorf("in file %q: %w", filename, err)
			}

			kind := i.(map[string]interface{})["kind"].(string)
			content, e := v1.Parse(c, o, filename, kind)

			if e != nil {
				allErrors = multierror.Append(allErrors, fmt.Errorf("error in %q: %w", filename, e))
			}
			parsedStructs = append(parsedStructs, content)
		}
	default:
		allErrors = multierror.Append(allErrors, fmt.Errorf("unsupported API Version in file %s", filename))
	}

	return parsedStructs, allErrors
}
