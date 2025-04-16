package constant

var (
	Debug = false
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
