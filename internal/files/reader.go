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
package files

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/OpenSLO/go-sdk/pkg/openslo"
	"github.com/OpenSLO/go-sdk/pkg/openslosdk"
)

// ReadObjects reads [openslo.Object] from the provided sources.
func ReadObjects(sources []string) ([]openslo.Object, error) {
	var allObjects []openslo.Object
	for _, src := range sources {
		data, err := readRawSchema(src)
		if err != nil {
			return nil, err
		}
		objects, err := readObjectsFromRawData(data)
		if err != nil {
			return nil, err
		}
		allObjects = append(allObjects, objects...)
	}
	return allObjects, nil
}

// readObjectsFromRawData reads [openslo.Object] from a byte slice.
func readObjectsFromRawData(data []byte) ([]openslo.Object, error) {
	format := openslosdk.FormatYAML
	if isJSONBuffer(data) {
		format = openslosdk.FormatJSON
	}
	return openslosdk.Decode(bytes.NewReader(data), format)
}

// readRawSchema reads raw OpenSLO schema from file path, HTTP address or stdin (path "-") to a byte slice.
func readRawSchema(path string) ([]byte, error) {
	switch {
	case isStdin(path):
		return io.ReadAll(os.Stdin)
	case isURL(path):
		// nolint: gosec,noctx
		// #nosec G107
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer func() { _ = resp.Body.Close() }()
		return io.ReadAll(resp.Body)
	default:
		return os.ReadFile(filepath.Clean(path))
	}
}

func isStdin(p string) bool {
	return p == "-"
}

func isURL(p string) bool {
	return strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://")
}

var jsonBufferRegex = regexp.MustCompile(`^\s*\[?\s*{`)

// isJSONBuffer scans the provided buffer, looking for an open brace indicating this is JSON.
// While a simple list like ["a", "b", "c"] is still a valid JSON,
// it does not really concern us when processing complex objects.
func isJSONBuffer(buf []byte) bool {
	return jsonBufferRegex.Match(buf)
}
