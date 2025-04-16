package constant

var (
	debug = false
)

func GetVersion() string {
	return "0.0.1_beta"
}

func GetToolName() string {
	return "go-forensic"
}

func GetDebug() bool {
	return debug
}

func SetDebug(_debug bool) {
	debug = _debug
}
