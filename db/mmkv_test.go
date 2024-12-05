package db

import (
	"testing"
)

func TestVersionEmpty(t *testing.T) {
	mmkv, err := NewMMKV()
	if err != nil {
		t.Error(err)
		return
	}
	err = mmkv.Init(`../testing/db/mmkv`).WithId(`userinfo`)
	if err != nil {
		t.Error(err)
		return
	}
	for _, key := range mmkv.GetKeys() {
		t.Log(key)
		t.Log(mmkv.GetValue(key))
	}
}
