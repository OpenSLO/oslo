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
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/OpenSLO/oslo/internal/pkg/yamlutils"
	"github.com/OpenSLO/oslo/pkg/manifest"
)

// validateStruct takes the given struct and validates it.
func validateStruct(c []manifest.OpenSLOKind) error {
	validate := validator.New()

	_ = validate.RegisterValidation("dateWithTime", isDateWithTimeValid)
	_ = validate.RegisterValidation("timeZone", isTimeZoneValid)
	_ = validate.RegisterValidation("validDuration", isValidDurationString)

	var allErrors []string
	for _, v := range c {
		if err := validate.Struct(v); err != nil {
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
	for _, file := range files {
		c, e := yamlutils.ReadConf(file)
		if e != nil {
			allErrors = append(allErrors, e.Error())
			break
		}

		content, _, err := yamlutils.Parse(c, file)
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
		_, err := time.Parse("2006-01-02T15:04:05Z", fl.Field().String())
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
