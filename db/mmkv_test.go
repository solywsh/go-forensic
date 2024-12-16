package db

import "testing"

func TestMMKV_SearchByKeywords(t *testing.T) {
	err := NewMMKV().
		SetFilePath("../testing/db/mmkv/mmkv.default").
		SearchByKeywords("color_choose_BACKGROUND")
	if err != nil {
		t.Log(err)
		return
	}
}
