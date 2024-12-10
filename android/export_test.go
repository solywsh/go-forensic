package android

import "testing"

func TestExportByKeywords(t *testing.T) {
	ad := NewAndroid().SetOutput("../temp/android")
	err := ad.ExportByKeywords("facebook", "tencent")
	if err != nil {
		t.Error(err)
	}
}
