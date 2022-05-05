/*
Copyright © 2022 OpenSLO Team

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
package validate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/oslo/pkg/manifest"
	v1 "github.com/OpenSLO/oslo/pkg/manifest/v1"
	"github.com/OpenSLO/oslo/pkg/manifest/v1alpha"
)

var (
	labelRegexp               = regexp.MustCompile(`^[\p{L}]([\_\-0-9\p{L}]*[0-9\p{L}])?$`)
	hasUpperCaseLettersRegexp = regexp.MustCompile(`[A-Z]+`)
)

// ReadConf reads in filename for a yaml file, and unmarshals it.
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

// Parse takes the provided byte array, parses it, and returns a parsed struct.
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

// validateStruct takes the given struct and validates it.
func validateStruct(c []manifest.OpenSLOKind) error {
	validate := validator.New()

	_ = validate.RegisterValidation("dateWithTime", isDateWithTimeValid)
	_ = validate.RegisterValidation("timeZone", isTimeZoneValid)
	_ = validate.RegisterValidation("labels", isValidLabel)
	_ = validate.RegisterValidation("validDuration", isValidDurationString)

	var allErrors []string
	for _, ival := range c {
		if err := validate.Struct(ival); err != nil {
			for _, err := range err.(validator.ValidationErrors) { //nolint: errorlint
				allErrors = append(allErrors, err.Error())
			}
		}
	}
	if len(allErrors) > 0 {
		return errors.New(strings.Join(allErrors, "\n"))
	}
	return nil
}

// validateFiles validates the given array of filenames.
func validateFiles(files []string) error {
	var allErrors []string
	for _, ival := range files {
		c, e := ReadConf(ival)
		if e != nil {
			allErrors = append(allErrors, e.Error())
			break
		}
		content, err := Parse(c, ival)
		if err != nil {
			allErrors = append(allErrors, err.Error())
			break
		}
		if validationErrors := validateStruct(content); validationErrors != nil {
			allErrors = append(allErrors, validationErrors.Error())
		}
	}
	if len(allErrors) > 0 {
		return errors.New(strings.Join(allErrors, "\n"))
	}
	return nil
}

// NewValidateCmd returns a new cobra.Command for the validate command.
func NewValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validates your yaml file against the OpenSLO spec.",
		Long:  `Validates your yaml file against the OpenSLO spec.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if e := validateFiles(args); e != nil {
				return e
			}

			fmt.Println("Valid!")
			return nil
		},
	}
}

func isValidDurationString(fl validator.FieldLevel) bool {
	for _, s := range []string{"s", "m", "h", "d", "w", "M", "Q", "Y"} {
		duration := fl.Field().String()
		if strings.HasSuffix(duration, s) {
			return true
		}
	}
	return false
}

func isDateWithTimeValid(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		_, err := time.Parse("2006-01-02 15:04:05", fl.Field().String())
		if err != nil {
			return false
		}
	}
	return true
}

func isTimeZoneValid(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		_, err := time.LoadLocation(fl.Field().String())
		if err != nil {
			return false
		}
	}
	return true
}

func isValidLabel(fl validator.FieldLevel) bool {
	labels := fl.Field().Interface().(v1.Labels)
	for key, values := range labels {
		if !validateLabel(key) {
			return false
		}
		if duplicates(values) {
			return false
		}
		for _, val := range values {
			// Validate only if len(val) > 0, in case where we have only key labels, there is always empty val string
			// and this is not an error
			if len(val) > 0 && !validateLabel(val) {
				return false
			}
		}
	}
	return true
}

func validateLabel(value string) bool {
	if len(value) > 63 || len(value) < 1 {
		return false
	}

	if !labelRegexp.MatchString(value) {
		return false
	}
	return !hasUpperCaseLettersRegexp.MatchString(value)
}

func duplicates(list []string) bool {
	duplicateFrequency := make(map[string]int)

	for _, item := range list {
		_, exist := duplicateFrequency[item]

		if exist {
			duplicateFrequency[item]++
		} else {
			duplicateFrequency[item] = 1
		}
		if duplicateFrequency[item] > 1 {
			return true
		}
	}
	return false
}
