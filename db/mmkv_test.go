package db

import "testing"

func TestVersionEmpty(t *testing.T) {
	t.Log(NewMmkv().
		InitMMKV(`../testing/db/mmkv/userinfo`).
		WithID(`userinfo`).
		AllKeys())
}
