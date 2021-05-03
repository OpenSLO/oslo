/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/validator.v2"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		type NewUserRequest struct {
			Username string `validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
			Name     string `validate:"nonzero"`
			Age      int    `validate:"min=21"`
			Password string `validate:"min=8"`
		}

		nur := NewUserRequest{Username: "something", Age: 20}
		err := validator.Validate(nur)
		if err == nil {
			color.Green("Valid!")
		} else {
			errs := err.(validator.ErrorMap)

			color.Red("Invalid")

			var errOuts []string
			for f, e := range errs {
				errOuts = append(errOuts, fmt.Sprintf("\t - %s (%v)\n", f, e))
			}

			for _, str := range errOuts {
				fmt.Print(str)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
