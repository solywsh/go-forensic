package logger

import "testing"

func TestNewLogger(t *testing.T) {
	log := NewLogger()
	log.Info("test", "k", "114514")
}
