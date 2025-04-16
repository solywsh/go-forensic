package ios

import (
	"testing"
)

func TestIOS_ExportAppDataByKeywords(t *testing.T) {
	err := NewSSHelper(
		WithPort("2222"),
		WithOutput("../temp"),
	).ExportAppDataByKeywords("passbook", "passes")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestIOS_ExportBySpecify(t *testing.T) {
	err := NewSSHelper(
		WithOutput("../temp")).
		ExportBySpecify("/var/mobile/Library/Passes")
	if err != nil {
		t.Error(err)
		return
	}
}
