package mmkv

type MMKV interface {
	WithId(id, decrypt string) error
	Keys() []string
	Count() uint64
	Get(key string) []byte
}

func FactoryNewMMKV(dir string) MMKV {
	// Due to some confidentiality reasons,
	// I cannot publicly disclose the specific parsing logic of my existing MMKV
	// until alternative methods are found.
	// If you have a corresponding solution,
	// please implement the MMKV interface and submit it to me,
	// and I will merge it.
	return NewMMKV(dir)
}
