//go:build windows
// +build windows

package osx

import (
	"fmt"
	"os/exec"
)

// RemoveAll
// If we export iOS files on Windows, the file name limitation may cause the deletion of files to fail.
// https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file#naming-conventions
// https://learn.microsoft.com/en-us/windows/win32/fileio/creating-and-opening-files
func RemoveAll(dirPath string) error {
	if dirPath == "" {
		return nil
	}
	// cmd
	//cmd := exec.Command("rd", "/s", "/q", fmt.Sprintf(`\\?\%s`, dirPath))
	// powershell
	cmd := exec.Command("powershell", "remove-Item", "-LiteralPath", fmt.Sprintf(`\\?\%s`, dirPath), "-Recurse", "-Force")
	return cmd.Run()
}
