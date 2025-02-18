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
package files

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed test-input
var testInput string

// This test is not run in parallel with others, because
// one of the subtest modifies global variable os.Stdin.
func TestReadConf(t *testing.T) { //nolint:tparallel
	expectedContent := []byte(testInput)

	t.Run("from filepath successfully", func(t *testing.T) {
		t.Parallel()
		const filePath = "./test-input"
		content, err := readRawSchema(filePath)
		require.NoErrorf(t, err, "can't read content from filepath %q", filePath)
		require.Equal(t, expectedContent, content)
	})

	t.Run("from URL successfully", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(testInput))
			require.NoError(t, err, "http test server can't serve properly")
		}))
		defer server.Close()

		content, err := readRawSchema(server.URL)
		require.NoErrorf(t, err, "can't read content from URL of the test server: %q", server.URL)
		require.Equal(t, expectedContent, content)
	})

	t.Run("from stdin successfully", func(t *testing.T) {
		t.Parallel()
		output, input, err := os.Pipe()
		require.NoError(t, err, "failed to create a pipe for stdin mock")
		_, err = input.Write(expectedContent)
		require.NoError(t, err, "failed to write to a pipe to mock user input via stdin")
		require.NoError(t, input.Close(), "failed to close a pipe that mocks stdin")

		// Restore stdin right after the test.
		defer func(v *os.File) { os.Stdin = v }(os.Stdin)
		os.Stdin = output

		const indicateStdin = "-"
		content, err := readRawSchema(indicateStdin)
		require.NoError(t, err, "can't read content from stdin")
		require.Equal(t, expectedContent, content)
	})
}
