package fsnotify

import (
	"os"
	"path/filepath"
	"strings"
)

var ignoredPaths = []string{
	"node_modules",
	".git",
	"vendor",
	"dist",
}

func ShouldIgnore(path string) bool {
	normalized := filepath.Clean(path)

	for _, ignore := range ignoredPaths {
		ignorePattern := string(os.PathSeparator) + ignore + string(os.PathSeparator)
		if strings.Contains(normalized, ignorePattern) ||
			strings.HasSuffix(normalized, string(os.PathSeparator)+ignore) ||
			strings.HasSuffix(normalized, ignorePattern) {
			return true
		}
	}

	return false
}
