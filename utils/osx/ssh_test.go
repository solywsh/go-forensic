package osx

import (
	"os"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	sshClient := NewSSH().SetPort("2222").SetPassword("alpine")
	res, err := sshClient.MustC().ExecuteCommand("ls -alh")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestExecuteCommandWithStream(t *testing.T) {
	sshClient := NewSSH().SetPort("2222").SetPassword("alpine")
	err := sshClient.MustC().ExecuteCommandWithStream("apt update", os.Stdout, os.Stderr)
	if err != nil {
		t.Error(err)
		return
	}
}
