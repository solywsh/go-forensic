package logger

import "testing"

func TestNewLogger(t *testing.T) {
	log := NewLogger()
	log.Info("test", "k", "114514")
	log.Debug("test", "k", "114514")
	dlog := NewDebugLogger()
	dlog.Debug("test", "k", "114514")
}
