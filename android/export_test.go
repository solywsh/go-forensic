package android

import "testing"

func TestExportByKeywords(t *testing.T) {
	ad := NewAndroid().SetOutput("../temp/android")
	err := ad.ExportAppDataByKeywords("facebook")
	if err != nil {
		t.Error(err)
	}
}

func TestAdbRunShellCommand(t *testing.T) {
	ad := NewAndroid()
	command, err := ad.AdbRunShellCommand("ls", "/data/user/0")
	if err != nil {
		return
	}
	t.Log(command)
}
