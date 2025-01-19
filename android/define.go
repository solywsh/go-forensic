package android

import (
	"github.com/electricbubble/gadb"
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/pathx"
	"github.com/solywsh/go-forensic/utils/printer"
	"os"
	"path/filepath"
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
		"/sdcard/Android/":  2,
		"/sdcard/Download/": 1,
		"/sdcard/Pictures/": 1,
		"/sdcard/Movies/":   1,
		"/sdcard/Music/":    1,
		"/sdcard/":          1,
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
	if !pathx.PathExists(localPath) {
		err := os.MkdirAll(filepath.Dir(localPath), 0755)
		if err != nil {
			return err
		}
	}
	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()
	return t.device.Pull(remotePath, localFile)
}

func (t *Android) removeOutput() {
	t.cleanPathOnce.Do(func() {
		t.spinner.Msg("cleaning " + t.output)
		osx.RemoveAll(t.output)
	})
}
