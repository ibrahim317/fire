package core

import (
	"os"
	"path/filepath"
	"sync"
)

// DefaultMapPath is the relative path to the custom map file (under project root).
const DefaultMapPath = "maps/custom_map.json"

var (
	projectRootOnce sync.Once
	projectRoot     string
)

// ResourcePath returns the absolute path for a resource file under the project root.
func ResourcePath(rel string) string {
	return filepath.Join(resolveProjectRoot(), rel)
}

func resolveProjectRoot() string {
	projectRootOnce.Do(func() {
		if root, err := findProjectRoot(); err == nil {
			projectRoot = root
		} else {
			projectRoot = "."
		}
	})
	return projectRoot
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "resources")); statErr == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}
