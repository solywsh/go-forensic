package db

/*
#cgo LDFLAGS: -L. -lmmkv
#include "mmkv.h"
*/
import "C"

type MmkvCli struct {
	mmkv C.MMkv
}

func NewMmkv() *MmkvCli {
	return &MmkvCli{}
}

func (m *MmkvCli) InitMMKV(dir string) *MmkvCli {
	C.InitializeMMKV(dir)
	return m
}

func (m *MmkvCli) WithID(filename string) *MmkvCli {
	m.mmkv = C.MMKVWithID(filename)
	return m
}

func (m *MmkvCli) GetString(key string) string {
	return m.mmkv.GetString(key)
}

func (m *MmkvCli) AllKeys() []string {
	return m.mmkv.AllKeys()
}
