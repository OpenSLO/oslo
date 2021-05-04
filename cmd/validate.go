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
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v3"
)

type NewUserRequest struct {
	Conf struct {
		Username string `validate:"nonzero,min=3,max=40,regexp=^[a-zA-Z]*$"`
		Name     string `validate:"nonzero"`
		Age      int    `validate:"min=21"`
		Password string `validate:"nonzero,min=8"`
	}
}

// reads in filename for a yaml file, and unmarshalls it
func readConf(filename string) (*NewUserRequest, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &NewUserRequest{}
	err = yaml.Unmarshal(buf, c)

	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}

func validate(c *NewUserRequest) {
	err := validator.Validate(c)
	if err == nil {
		fmt.Println("Valid!")
	} else {
		errs := err.(validator.ErrorMap)

		fmt.Println("Invalid")

		var errOuts []string
		for f, e := range errs {
			errOuts = append(errOuts, fmt.Sprintf("  - %s (%v)\n", f, e))
		}

		for _, str := range errOuts {
			fmt.Print(str)
		}
	}
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

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates your yaml file against the OpenSLO spec",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		validateFiles(args)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
