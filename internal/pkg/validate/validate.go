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
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/OpenSLO/oslo/internal/pkg/yamlutils"
	"github.com/OpenSLO/oslo/pkg/manifest"
	v1 "github.com/OpenSLO/oslo/pkg/manifest/v1"
)

var (
	labelRegexp               = regexp.MustCompile(`^[\p{L}]([\_\-0-9\p{L}]*[0-9\p{L}])?$`)
	hasUpperCaseLettersRegexp = regexp.MustCompile(`[A-Z]+`)
)

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

// Files validates the given array of filenames.
func Files(files []string) error {
	var allErrors []string
	for _, ival := range files {
		c, e := yamlutils.ReadConf(ival)
		if e != nil {
			allErrors = append(allErrors, e.Error())
			break
		}

		// prints ival to stdout
		fmt.Println(ival)

		content, err := yamlutils.Parse(c, ival)
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

func isValidDurationString(fl validator.FieldLevel) bool {
	for _, s := range []string{"m", "h", "d", "w", "M", "Q", "Y"} {
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
