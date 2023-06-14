package main

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/OpenSLO/oslo/pkg/discoverfiles"
	"github.com/OpenSLO/oslo/pkg/validate"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)

	filePaths := []string{filepath.Join(filepath.Dir(filename), "definitions")}
	discoveredFilePaths, err := discoverfiles.DiscoverFilePaths(filePaths, true)
	if err != nil {
		panic(err)
	}
	if err = validate.Files(discoveredFilePaths); err != nil {
		panic(err)
	}
	fmt.Println("Valid!")
}
