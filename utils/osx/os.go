package osx

import (
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
)

func TarDecompression(filePath, destinationDir string) error {
	cmd := exec.Command("tar", "-xf", filePath, "-C", destinationDir)
	cmd.Run()
	return nil
}

func IsTTY() bool {
	// some output needs to be run in the terminal of the tty, such as spinner
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
