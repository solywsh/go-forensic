package db

import (
	"testing"
)

func TestNewMMKV(t *testing.T) {
	m := NewMMKV(`../testing/db/mmkv`)
	err := m.WithId(`userinfo`, "")
	if err != nil {
		t.Error(err)
	}
	for _, k := range m.Keys() {
		t.Log("key:", k)
		t.Log("value:", string(m.Get(k)))
	}
}

func TestDecryptFile(t *testing.T) {
	m := NewMMKV(`../testing/db/mmkv`)
	err := m.WithId(`pref_id_search_history`, "fd9f5bef68c54a1ecf70757a6d6f565b")
	if err != nil {
		t.Error(err)
	}
	t.Log("count:", m.Count())
	for _, k := range m.Keys() {
		t.Log("key:", k)
		t.Log("value:", string(m.Get(k)))
	}
}
