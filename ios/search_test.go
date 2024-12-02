package ios

import "testing"

func TestIOS_SearchByKeywords(t *testing.T) {
	err := NewIOS().SetPort("23").SearchByKeywords("whatsapp")
	if err != nil {
		t.Error(err)
		return
	}
}
