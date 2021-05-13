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
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

// readConf reads in filename for a yaml file, and unmarshals it.
func readConf(filename string) (interface{}, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var m ObjectGeneric

	if err := yaml.Unmarshal(fileContent, &m); err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	fmt.Println(m.Kind)

	switch m.Kind {
	case "Service":
		var content Service
		if err := yaml.Unmarshal(fileContent, &content); err != nil {
			return nil, fmt.Errorf("in file %q: %w", filename, err)
		}
		return content, nil
	case "SLO":
		var content sloSpec
		if err := yaml.Unmarshal(fileContent, &content); err != nil {
			return sloSpec{}, fmt.Errorf("in file %q: %w", filename, err)
		}
		return content, nil
	default:
		return nil, fmt.Errorf("Unsupported kind: %s", m.Kind)
	}
}

func validateStruct(c interface{}) {
	validate = validator.New()

	err := validate.Struct(c)
	fmt.Println(c)

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println(err)
		}

		// from here you can create your own error messages in whatever language you wish
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
		validateStruct(c)
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
