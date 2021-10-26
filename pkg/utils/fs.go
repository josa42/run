package utils

import (
	"os"
	"strings"
)

func IsSubDir(dir, parent string) bool {
	if dir == parent {
		return true
	}
	return strings.HasPrefix(dir, parent+"/")
}

func Abs(dir string) string {
	home := os.Getenv("HOME")

	if strings.HasPrefix(dir, "~") {
		return strings.Replace(dir, "~", home, 1)
	}

	return dir
}
