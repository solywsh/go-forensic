package ssh

import (
	"os"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	sshClient := NewClient(
		WithPort("2222"),
		WithPassword("alpine"),
	)
	res, err := sshClient.MustC().ExecuteCommand("ls -alh")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestExecuteCommandWithStream(t *testing.T) {
	sshClient := NewClient(
		WithPort("2222"),
		WithPassword("alpine"),
	)
	err := sshClient.MustC().ExecuteCommandWithStream("apt update", os.Stdout, os.Stderr)
	if err != nil {
		t.Error(err)
		return
	}
}
