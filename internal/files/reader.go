package files

import (
	"bytes"
	"fmt"
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
// It returns a map where the key is a file path and the value are objects read form this file.
func ReadObjects(sources []string) (map[string][]openslo.Object, error) {
	allObjects := make(map[string][]openslo.Object)
	for _, src := range sources {
		objects, err := readObjectsFromSource(src)
		if err != nil {
			return nil, fmt.Errorf("failed to read objects from %s: %w", src, err)
		}
		allObjects[src] = objects
	}
	return allObjects, nil
}

func readObjectsFromSource(source string) ([]openslo.Object, error) {
	data, err := readRawSchema(source)
	if err != nil {
		return nil, err
	}
	return readObjectsFromRawData(data)
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
