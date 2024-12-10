package osx

import (
	"bytes"
	"fmt"
	"github.com/solywsh/go-forensic/utils"
	"os/exec"
	"strings"
)

var (
	abdUnrelatedOutPut = []string{
		"adb server is out of date",
		"daemon started successfully",
	}
)

func RunADBShellCommand(cmdArgs ...string) (string, error) {
	cmd := exec.Command("adb", append([]string{"shell", "su", "-c"}, cmdArgs...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %s, error: %s", stderr.String(), err)
	}
	output := strings.TrimSpace(out.String())
	// filter out common non-command result information
	lines := strings.Split(output, "\n")
	var filteredLines []string
	for _, line := range lines {
		if utils.ContainString(line, abdUnrelatedOutPut) {
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	return strings.Join(filteredLines, "\n"), nil
}

func RunADBCommand(cmdArgs ...string) (string, error) {
	cmd := exec.Command("adb", cmdArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %s, error: %s", stderr.String(), err)
	}
	output := strings.TrimSpace(out.String())
	lines := strings.Split(output, "\n")
	var filteredLines []string
	for _, line := range lines {
		if utils.ContainString(line, abdUnrelatedOutPut) {
			continue
		}
		filteredLines = append(filteredLines, line)
	}
	return strings.Join(filteredLines, "\n"), nil
}
