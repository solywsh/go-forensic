package utils

import "strings"

func Contain[T comparable](s T, arr []T) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func ContainString(s string, arr []string) bool {
	for _, v := range arr {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}
