package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func IsSubDir(dir, pattern string) (bool, string) {
	if dir == pattern {
		return true, dir
	}

	if strings.Contains(pattern, "*") {

		str := pattern
		str = strings.Replace(str, "**", "##WILDCARD:ANY##", -1)
		str = strings.Replace(str, "*", "##WILDCARD##", -1)
		str = regexp.QuoteMeta(str)
		str = strings.Replace(str, "##WILDCARD:ANY##", ".*", -1)
		str = strings.Replace(str, "##WILDCARD##", "[^/]*", -1)
		str = regexp.MustCompile("/+$").ReplaceAllString(str, "")
		str = fmt.Sprintf("^(%s)(/.*)?$", str)
		expr := regexp.MustCompile(str)

		if expr.MatchString(dir) {
			return true, expr.ReplaceAllString(dir, "$1")
		}

		return false, ""
	}

	return strings.HasPrefix(dir, pattern+"/"), dir
}

func Abs(dir string) string {
	home := os.Getenv("HOME")

	if strings.HasPrefix(dir, "~") {
		return strings.Replace(dir, "~", home, 1)
	}

	return dir
}
