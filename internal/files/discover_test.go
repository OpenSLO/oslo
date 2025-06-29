package files_test

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/OpenSLO/oslo/internal/files"
)

// TestDiscoverFilePaths tests it on real filesystem,
// content of those file doesn't matter, thus they're empty.
func TestDiscoverFilePaths(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		filePaths     []string
		recursive     bool
		want          []string
		expectedError error
	}{
		{
			name:      "path to single file",
			filePaths: []string{"testdata/discover/x.yml"},
			recursive: false,
			want:      []string{"testdata/discover/x.yml"},
		},
		{
			name:      "path to single file with recursive",
			filePaths: []string{"testdata/discover/x.yml"},
			recursive: true,
			want:      []string{"testdata/discover/x.yml"},
		},
		{
			name:      "path to multiple files and stdin",
			filePaths: []string{"testdata/discover/x.yml", "-", "testdata/discover/a/b/b1.yml"},
			recursive: false,
			want:      []string{"testdata/discover/x.yml", "-", "testdata/discover/a/b/b1.yml"},
		},
		{
			name:          "path to non-existence file and directory",
			filePaths:     []string{"testdata/discover/a/b/b1.yml", "testdata/discover/non-existing.yml"},
			recursive:     false,
			expectedError: fs.ErrNotExist,
		},
		{
			name:          "path to non-existence file and directory with recursive",
			filePaths:     []string{"testdata/discover/a/b/b1.yml", "testdata/discover/non-existing.yml"},
			recursive:     true,
			expectedError: fs.ErrNotExist,
		},
		{
			name:      "path to directory with subdirectories and additional file",
			filePaths: []string{"testdata/discover", "testdata/discover/a/b/b2.yml"},
			recursive: false,
			want:      []string{"testdata/discover/x.yml", "testdata/discover/y.yaml", "testdata/discover/a/b/b2.yml"},
		},
		{
			name:      "path to directory with subdirectories and to stdin with recursive",
			filePaths: []string{"testdata/discover", "-"},
			recursive: true,
			want: []string{
				"testdata/discover/a/a1.yml", "testdata/discover/a/a2.yml",
				"testdata/discover/a/b/b1.yml", "testdata/discover/a/b/b2.yml",
				"testdata/discover/aa/aa1.yml",
				"testdata/discover/x.yml", "testdata/discover/y.yaml",
				"-",
			},
		},
		{
			name:      "path to directory with subdirectories, stdin and URLs with recursive",
			filePaths: []string{"testdata/discover", "-", "http://example.com/file-1", "https://example.com/file-2"},
			recursive: true,
			want: []string{
				"testdata/discover/a/a1.yml", "testdata/discover/a/a2.yml",
				"testdata/discover/a/b/b1.yml", "testdata/discover/a/b/b2.yml",
				"testdata/discover/aa/aa1.yml",
				"testdata/discover/x.yml", "testdata/discover/y.yaml",
				"-",
				"http://example.com/file-1",
				"https://example.com/file-2",
			},
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.name, func(t *testing.T) {
			t.Parallel()
			res, err := files.Discover(tC.filePaths, tC.recursive)
			require.ErrorIs(t, err, tC.expectedError)
			require.Equal(t, tC.want, res)
		})
	}
}
