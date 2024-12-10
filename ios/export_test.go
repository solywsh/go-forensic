package ios

import "testing"

func TestIOS_ExportAppDataByKeywords(t *testing.T) {
	err := NewIOS().SetPort("2222").
		SetOutput("../temp").
		ExportAppDataByKeywords("whatsapp")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestIOS_ExportBySpecify(t *testing.T) {
	err := NewIOS().
		SetOutput("../temp").
		ExportBySpecify("/var/mobile/Library/Passes")
	if err != nil {
		t.Error(err)
		return
	}
}
