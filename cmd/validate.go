/*
Copyright © 2021 OpenSLO Team

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
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/oslo/pkg/manifest"
	"github.com/OpenSLO/oslo/pkg/manifest/v1alpha"
)

// readConf reads in filename for a yaml file, and unmarshals it.
func readConf(filename string) ([]byte, error) {
	if filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

// parse takes the provided byte array, parses it, and returns a parsed struct.
func parse(fileContent []byte, filename string) ([]interface{}, error) {
	var m manifest.ObjectGeneric

	if err := yaml.Unmarshal(fileContent, &m); err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	var allErrors []string
	var parsedStructs []interface{}
	switch m.APIVersion {
	case v1alpha.APIVersion:
		content, e := v1alpha.Parse(fileContent, m, filename)
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
func validateStruct(c []interface{}) error {
	validate := validator.New()

	_ = validate.RegisterValidation("dateWithTime", isDateWithTimeValid)
	_ = validate.RegisterValidation("timeZone", isTimeZoneValid)

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
		c, e := readConf(ival)
		if e != nil {
			allErrors = append(allErrors, e.Error())
			break
		}
		content, err := parse(c, ival)
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

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validates your yaml file against the OpenSLO spec",
		Long:  `TODO`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if e := validateFiles(args); e != nil {
				fmt.Println(e.Error())
				os.Exit(1)
			}
			fmt.Println("Valid!")
		},
	}
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
