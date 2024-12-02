package ios

import "testing"

func TestIOS_SearchByKeywords(t *testing.T) {
	err := NewIOS().SetPort("2223").
		SetOutput("../temp").
		SearchByKeywords("whatsapp")
	if err != nil {
		t.Error(err)
		return
	}
}
