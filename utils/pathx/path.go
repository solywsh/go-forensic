package pathx

import (
	"os"
	"path/filepath"
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
