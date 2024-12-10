package android

import (
	"github.com/solywsh/go-forensic/utils/logger"
	"github.com/solywsh/go-forensic/utils/printer"
)

var (
	androidFocusMap = map[string]int{
		/*
			"/data/user/0/",
			"/data/user/../",
		*/
		"/data/user/": 2,
		"/data/data":  1,
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
