package utils

import "testing"

func TestIsPortAvailable(t *testing.T) {
	t.Log(IsPortAvailable("tcp", 2222))
}
