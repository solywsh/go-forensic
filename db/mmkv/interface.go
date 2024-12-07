package mmkv

type MMKV interface {
	WithId(id, decrypt string) error
	Keys() []string
	Count() uint64
	Get(key string) []byte
}

func FactoryNewMMKV(dir string) MMKV {
	return NewMMKV(dir)
}
