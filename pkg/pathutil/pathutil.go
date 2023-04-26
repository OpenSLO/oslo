package pathutil

import "strings"

func IsStdin(p string) bool {
	return p == "-"
}

func IsURL(p string) bool {
	return strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://")
}
