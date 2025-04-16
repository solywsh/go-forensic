package logger

import (
	"github.com/charmbracelet/log"
	"os"
	"sync"
	"time"
)

var (
	logger     *log.Logger
	loggerOnce sync.Once
)

func NewLogger() *log.Logger {
	loggerOnce.Do(func() {
		logger = log.New(os.Stderr)
		logger.SetReportTimestamp(true)
		logger.SetTimeFormat(time.Kitchen)
		logger.SetReportCaller(true)
	})
	return logger
}

func NewDebugLogger() *log.Logger {
	loggerOnce.Do(func() {
		logger = log.New(os.Stderr)
		logger.SetReportTimestamp(true)
		logger.SetTimeFormat(time.Kitchen)
		logger.SetReportCaller(true)
		logger.SetLevel(log.DebugLevel)
	})
	return logger
}
