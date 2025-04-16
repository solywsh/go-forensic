package utils

import (
	"fmt"
	"github.com/spf13/cast"
	"net"
	"strconv"
)

func FormatNetSpeed(speed int) string {
	cast.ToString(speed / 1000000)
	if speed < 1000 {
		return fmt.Sprintf("%d b/s", speed)
	} else if speed < 1000000 {
		return fmt.Sprintf("%d Kb/s", speed/1000)
	} else if speed < 1000000000 {
		return fmt.Sprintf("%d Mb/s", speed/1000000)
	} else {
		return fmt.Sprintf("%d Gb/s", speed/1000000000)
	}
}

func IsPortAvailable(protocol string, port int) bool {
	addr := "127.0.0.1:" + strconv.Itoa(port)
	listener, err := net.Listen(protocol, addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
