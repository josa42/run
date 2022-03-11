package utils

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	notFound  = errors.New("Not found")
	home      = os.Getenv("HOME")
	rootFiles = []string{".git", "go.mod", "package.json"}
)

func FindUp(dir string, name string) (string, error) {
	filePath := filepath.Join(dir, name)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		if dir == "/" || dir == home {
			return "", notFound
		}

		return FindUp(filepath.Dir(dir), name)
	}

	if err != nil {
		return "", err
	}

	return dir, nil
}

func FindProjectRoot(dir string) string {
	for _, file := range rootFiles {
		if root, _ := FindUp(dir, file); root != "" {
			return root
		}
	}

	return ""
}
