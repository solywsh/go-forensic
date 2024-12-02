package logger

import (
	"github.com/charmbracelet/log"
	"os"
	"sync"
	"time"
)

var (
	logger             *log.Logger
	newLoggerOnce      sync.Once
	debugLogger        *log.Logger
	newDebugLoggerOnce sync.Once
)

func NewLogger() *log.Logger {
	newLoggerOnce.Do(func() {
		logger = log.New(os.Stderr)
		logger.SetReportTimestamp(true)
		logger.SetTimeFormat(time.Kitchen)
		logger.SetReportCaller(true)
	})
	return logger
}

func NewDebugLogger() *log.Logger {
	newDebugLoggerOnce.Do(func() {
		debugLogger = log.New(os.Stderr)
		debugLogger.SetReportTimestamp(true)
		debugLogger.SetTimeFormat(time.Kitchen)
		debugLogger.SetReportCaller(true)
		debugLogger.SetLevel(log.DebugLevel)
	})
	return debugLogger
}
