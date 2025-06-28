package main

import (
	"runtime/debug"
	"strings"

	"github.com/OpenSLO/oslo/internal/cli"
	"github.com/spf13/cobra"
)

// version is set during build time.
var version string

func main() {
	root := cli.NewRootCmd(getBuildVersion(version))
	cobra.CheckErr(root.Execute())
}

func getBuildVersion(version string) string {
	if version == "" {
		version = getRuntimeVersion()
	}
	return strings.TrimSpace(version)
}

func getRuntimeVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "(devel)" {
		return "0.0.0"
	}
	return strings.TrimPrefix(info.Main.Version, "v")
}
