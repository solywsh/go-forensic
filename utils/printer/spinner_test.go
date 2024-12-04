package printer

import (
	"testing"
	"time"
)

func TestSpinnerX_Run(t *testing.T) {
	// tty in command
	s := NewSpinnerX().Msg("Loading...")
	s.Run()
	time.Sleep(2 * time.Second)
	s.Msg("Done")
	time.Sleep(2 * time.Second)
	s.Quit()
}
