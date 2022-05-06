package yamlutils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
func Parse(fileContent []byte, filename string) ([]manifest.OpenSLOKind, error) {
	var m manifest.ObjectGeneric

	// unmarshal here to get the APIVersion so we can process the file correctly
	if err := yaml.Unmarshal(fileContent, &m); err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	var allErrors []string
	var parsedStructs []manifest.OpenSLOKind
	switch m.APIVersion {
	// This is where we add new versions of the OpenSLO spec.
	case v1alpha.APIVersion:
		// unmarshal again to get the v1alpha struct
		var o v1alpha.ObjectGeneric
		if err := yaml.Unmarshal(fileContent, &o); err != nil {
			return nil, fmt.Errorf("in file %q: %w", filename, err)
		}

		content, e := v1alpha.Parse(fileContent, o, filename)
		if e != nil {
			allErrors = append(allErrors, e.Error())
		}
		parsedStructs = append(parsedStructs, content)
	case v1.APIVersion:
		// unmarshal again to get the v1 struct
		var o v1.ObjectGeneric
		if err := yaml.Unmarshal(fileContent, &o); err != nil {
			return nil, fmt.Errorf("in file %q: %w", filename, err)
		}

		content, e := v1.Parse(fileContent, o, filename)
		if e != nil {
			allErrors = append(allErrors, e.Error())
		}
		parsedStructs = append(parsedStructs, content)
	default:
		allErrors = append(allErrors, fmt.Sprintf("Unsupported API Version in file %s", filename))
	}
	if len(allErrors) > 0 {
		return nil, errors.New(strings.Join(allErrors, "\n"))
	}

	return parsedStructs, nil
}
