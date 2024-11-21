package files

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

// Discover returns all file paths that come from file paths provided as the argument.
// Return directly if they are standard files. For directories list all files available in its
// root or recursively traverse all subdirectories and find files in them when the argument recursive
// is true. For "-" path that indicates standard input or starts with http:// or https://,
// it returns the path directly in the same way as other paths.
func Discover(filePaths []string, recursive bool) ([]string, error) { //nolint:cyclop
	var discoveredPaths []string
	for _, p := range filePaths {
		// Indicates that a file should be read from standard input.
		// Code that consumes file paths needs to handle "-"
		// by reading from os.Stdin in such case.
		if isStdin(p) {
			discoveredPaths = append(discoveredPaths, p)
			continue
		}
		// Content should be downloaded from this URL and treated as others.
		if isURL(p) {
			discoveredPaths = append(discoveredPaths, p)
			continue
		}

		// When path is valid and it's not a directory, use it directly.
		fInfo, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if !fInfo.IsDir() {
			discoveredPaths = append(discoveredPaths, p)
			continue
		}

		// When recursive is true and the path is a directory,
		// discover all files in it and its subdirectories.
		if recursive {
			if walkErr := filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					discoveredPaths = append(discoveredPaths, path)
				}
				return nil
			}); walkErr != nil {
				return nil, walkErr
			}
			continue
		}

		// When recursive is false and the path is a directory,
		// get only paths for files in the root of it.
		entries, err := os.ReadDir(p)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if !e.IsDir() {
				discoveredPaths = append(discoveredPaths, path.Join(p, e.Name()))
			}
		}
	}
	return discoveredPaths, nil
}
