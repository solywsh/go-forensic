package pathx

import (
	"os"
	"path/filepath"
	"strings"
)

func PathJoin(arg ...string) string {
	return filepath.ToSlash(filepath.Join(arg...))
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return true
}

func FilterSubDir(paths []string) []string {
	var result []string
	for i, path := range paths {
		isSubdirectory := false
		for j, otherPath := range paths {
			if i != j && strings.HasPrefix(path, otherPath) && len(path) > len(otherPath) {
				isSubdirectory = true
				break
			}
		}
		if !isSubdirectory {
			result = append(result, path)
		}
	}
	return result
}
