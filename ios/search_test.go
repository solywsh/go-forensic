package ios

import "testing"

func TestSearch(t *testing.T) {
	err := NewIOS().SetPort("23").SearchKeywords("yippi")
	if err != nil {
		t.Error(err)
		return
	}
}
