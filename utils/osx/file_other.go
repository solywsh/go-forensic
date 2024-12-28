//go:build !windows
// +build !windows

package osx

import "os"

func RemoveAll(dirPath string) error {
	return os.RemoveAll(dirPath)
}
