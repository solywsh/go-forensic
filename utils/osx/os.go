package osx

import "os/exec"

func TarDecompression(filePath, destinationDir string) error {
	cmd := exec.Command("tar", "-xf", filePath, "-C", destinationDir)
	cmd.Run()
	return nil
}
