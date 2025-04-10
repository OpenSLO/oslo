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

func TestReadConf(t *testing.T) {
	expectedContent := []byte(testInput)

	t.Run("from filepath successfully", func(t *testing.T) {
		const filePath = "./test-input"
		content, err := readRawSchema(filePath)
		require.NoErrorf(t, err, "can't read content from filepath %q", filePath)
		require.Equal(t, expectedContent, content)
	})

	t.Run("from URL successfully", func(t *testing.T) {
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
