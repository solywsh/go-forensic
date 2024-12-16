package utils

import (
	"golang.org/x/exp/constraints"
	"strings"
)

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

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
