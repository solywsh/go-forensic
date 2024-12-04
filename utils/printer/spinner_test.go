package printer

import (
	"testing"
	"time"
)

func TestSpinnerX_Run(t *testing.T) {

	s := NewSpinnerX().Msg("Loading...")
	s.Run()
	time.Sleep(2 * time.Second)
	s.Msg("Done")
	time.Sleep(2 * time.Second)
	s.Quit()
}

/*
// tty in command

func main() {
	s := NewSpinnerX().Msg("Loading...")
	s.Run()
	time.Sleep(2 * time.Second)
	s.SetSpinner(spinner.Globe).Msg("Done")
	time.Sleep(2 * time.Second)
	s.Quit()
}
*/
