package constant

import "time"

var (
	Debug = false
)

const (
	timeout  = 5 * time.Second
	interval = 50 * time.Millisecond
)

func GetVersion() string {
	return "0.0.1_beta"
}

func GetToolName() string {
	return "go-forensic"
}

func GetDebug() bool {
	return Debug
}

func SetDebug(_debug bool) {
	Debug = _debug
}

func GetTimeout() time.Duration {
	return timeout
}

func GetInterval() time.Duration {
	return interval
}
