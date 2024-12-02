package db

import "testing"

func TestSqliteX_SearchByKeywords(t *testing.T) {
	err := NewSqlite().
		SetDbPath("../testing/db/sqlite/ChatStorage.sqlite").
		SearchByKeywords("whatsapp")
	if err != nil {
		t.Log(err)
		return
	}
}
