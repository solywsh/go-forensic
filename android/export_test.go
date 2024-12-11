package android

import "testing"

func TestExportByKeywords(t *testing.T) {
	ad := NewAndroid().SetOutput("../temp/android")
	err := ad.ExportAppDataByKeywords("facebook")
	if err != nil {
		t.Error(err)
	}
}
