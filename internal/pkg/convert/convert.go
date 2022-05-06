/*
Package convert provides a command to convert from openslo to other formats.

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
package convert

import (
	"fmt"
	"io"
	"os"
)

// RemoveDuplicates to remove duplicate string from a slice.
func RemoveDuplicates(s []string) []string {
	result := make([]string, 0, len(s))
	m := make(map[string]bool)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = true
		}
	}
	for k := range m {
		result = append(result, k)
	}
	return result
}

// ConvertFiles converts the provided file.
func ConvertFiles(out io.Writer, filenames []string) error {
	// read each file and format it
	for _, filename := range filenames {
		fmt.Fprintln(out, "Converting", filename)
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return nil
}
