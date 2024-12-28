package android

import (
	"bytes"
	"github.com/electricbubble/gadb"
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/printer"
	"os"
	"sync"
)

var (
	androidFocusMap = map[string]int{
		/*
			"/data/user/0/",
			"/data/user/../",
		*/
		"/data/user/": 2,
		//"/data/data":  1,
		/*
			"/sdcard/Android/data",
			"/sdcard/Android/media",
			"/sdcard/Android/obb",
		*/
		"/sdcard/Android/": 2,
	}

	log = logger.NewLogger()
)

type Android struct {
	output         string
	cleanDirBefore bool
	spinner        *printer.SpinnerX
	keywords       []string
	device         gadb.Device
	cleanPathOnce  sync.Once
}

func NewAndroid() *Android {
	return &Android{
		output:         "./temp",
		cleanDirBefore: true,
	}
}

func (t *Android) SetOutput(output string) *Android {
	t.output = output
	return t
}

func (t *Android) SetCleanDirBefore(cleanDirBefore bool) *Android {
	t.cleanDirBefore = cleanDirBefore
	return t
}

func (t *Android) AdbRunShellCommand(args ...string) (string, error) {
	return t.device.RunShellCommand("su", append([]string{"-c"}, args...)...)
}

func (t *Android) AdbPullFile(remotePath, localPath string) error {
	buffer := bytes.NewBuffer(nil)
	err := t.device.Pull(remotePath, buffer)
	if err != nil {
		return err
	}
	return os.WriteFile(localPath, buffer.Bytes(), os.ModePerm)
}

func (t *Android) removeOutput() {
	t.cleanPathOnce.Do(func() {
		t.spinner.Msg("cleaning " + t.output)
		osx.RemoveAll(t.output)
	})
}
