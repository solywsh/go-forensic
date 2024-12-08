package ios

import "testing"

func TestIOS_SearchByKeywords(t *testing.T) {
	err := NewIOS().SetPort("2222").
		SetOutput("../temp").
		ExportAppDataByKeywords("whatsapp")
	if err != nil {
		t.Error(err)
		return
	}
}
