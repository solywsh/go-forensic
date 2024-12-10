package osx

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
)

func TarDecompression(filePath, destinationDir string) error {
	cmd := exec.Command("tar", "-xf", filePath, "-C", destinationDir)
	cmd.Run()
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("Failed to delete local tar file: %s\n", err)
	}
	return nil
}

func IsTTY() bool {
	// some output needs to be run in the terminal of the tty, such as spinner
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
