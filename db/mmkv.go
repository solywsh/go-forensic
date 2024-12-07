package db

import "github.com/solywsh/go-forensic/db/mmkv"

type MMKV struct {
	mmkv.MMKV
}

func NewMMKV(dir string) *MMKV {
	return &MMKV{mmkv.FactoryNewMMKV(dir)}
}
