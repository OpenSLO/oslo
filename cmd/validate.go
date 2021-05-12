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
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v3"
)

// readConf reads in filename for a yaml file, and unmarshals it.
func readConf(filename string) (interface{}, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return newUserRequest{}, err
	}

	// theres probably a better way of handling this but what
	// we are going to do is unmarshal using a generic schema in
	// order to infer the kind, and use that to load the correct
	// schema and use that to unmarshal
	m := make(map[interface{}]interface{})

	if err := yaml.Unmarshal(fileContent, &m); err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	switch m["kind"] {
	case "Service":
		var content serviceSpec
		if err := yaml.Unmarshal(fileContent, &content); err != nil {
			return serviceSpec{}, fmt.Errorf("in file %q: %w", filename, err)
		}
		return content, nil
	case "SLO":
		var content sloSpec
		if err := yaml.Unmarshal(fileContent, &content); err != nil {
			return sloSpec{}, fmt.Errorf("in file %q: %w", filename, err)
		}
		return content, nil
	default:
		var content newUserRequest
		if err := yaml.Unmarshal(fileContent, &content); err != nil {
			return newUserRequest{}, fmt.Errorf("in file %q: %w", filename, err)
		}
		return content, nil
	}
}

func validate(c interface{}) {
	if err := validator.Validate(c); err != nil {
		var errs validator.ErrorMap
		errors.As(err, &errs)

		fmt.Println("Invalid")

		for f, e := range errs {
			fmt.Printf("  - %s (%v)\n", f, e)
		}
		return
	}
	fmt.Println("Valid!")
}

func validateFiles(files []string) {
	for _, ival := range files {
		c, err := readConf(ival)
		if err != nil {
			log.Fatal(err)
		}
		validate(c)
	}
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validates your yaml file against the OpenSLO spec",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			validateFiles(args)
		},
	}
}
