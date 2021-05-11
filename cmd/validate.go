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
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v3"
)

type newUserRequest struct {
	Conf struct {
		Username string `validate:"nonzero,min=3,max=40,regexp=^[a-zA-Z]*$"`
		Name     string `validate:"nonzero"`
		Age      int    `validate:"min=21"`
		Password string `validate:"nonzero,min=8"`
	}
}

// readConf reads in filename for a yaml file, and unmarshals it.
func readConf(filename string) (newUserRequest, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return newUserRequest{}, err
	}

	var content newUserRequest
	if err := yaml.Unmarshal(fileContent, &content); err != nil {
		return newUserRequest{}, fmt.Errorf("in file %q: %w", filename, err)
	}

	return content, nil
}

func validate(c newUserRequest) error {
	if err := validator.Validate(c); err != nil {
		var errs validator.ErrorMap
		errors.As(err, &errs)

		fmt.Println("Invalid")

		for f, e := range errs {
			fmt.Printf("  - %s (%v)\n", f, e)
		}
		return errors.New("Error in validation")
	}
	fmt.Println("Valid!")
	return nil
}

func validateFiles(files []string) (int, error) {
	for _, ival := range files {
		c, err := readConf(ival)
		if err != nil {
			return -1, errors.New("Error in validation")
		}
		if err := validate(c); err != nil {
			return -1, errors.New("Error in validation")
		}
	}
	return 0, nil
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validates your yaml file against the OpenSLO spec",
		Long:  `TODO`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := validateFiles(args); err != nil {
				return err
			}
			return nil
		},
	}
}
