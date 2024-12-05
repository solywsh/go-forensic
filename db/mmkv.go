package db

import "C"
import (
	"encoding/json"
	"errors"
	"github.com/ebitengine/purego"
)

type MMKV struct {
	init     func(string)
	withId   func(string) error
	getKeys  func() []string
	getValue func(string) string
	close    func()
}

func NewMMKV() (*MMKV, error) {
	libc, err := purego.Dlopen("./mmkv/mmkv.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, err
	}
	var mk = &MMKV{}
	var InitializeMMKV func(string)
	purego.RegisterLibFunc(&InitializeMMKV, libc, "InitializeMMKV")
	mk.init = InitializeMMKV
	var MMKVWithID func(string) string
	purego.RegisterLibFunc(&MMKVWithID, libc, "MMKVWithID")
	mk.withId = func(id string) error {
		errStr := MMKVWithID(id)
		if errStr != "" {
			return errors.New(errStr)
		}
		return nil
	}
	var GetAllKeys func() string
	purego.RegisterLibFunc(&GetAllKeys, libc, "GetAllKeys")
	mk.getKeys = func() []string {
		keysStr := GetAllKeys()
		var keys []string
		if keysStr != "" {
			err := json.Unmarshal([]byte(keysStr), &keys)
			if err != nil {
				log.Error(err)
				return nil
			}
		}
		return keys
	}
	var GetValue func(string) string
	purego.RegisterLibFunc(&GetValue, libc, "GetValue")
	mk.getValue = GetValue

	var Close func()
	purego.RegisterLibFunc(&Close, libc, "Close")
	mk.close = Close
	return mk, nil
}

func (m *MMKV) Init(dir string) *MMKV {
	m.init(dir)
	return m
}

func (m *MMKV) WithId(filename string) error {
	return m.withId(filename)
}

func (m *MMKV) GetValue(key string) string {
	return m.getValue(key)
}

func (m *MMKV) GetKeys() []string {
	return m.getKeys()
}

func (m *MMKV) Close() {
	m.close()
}
