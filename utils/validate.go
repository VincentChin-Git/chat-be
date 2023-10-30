package utils

import "unicode"

func IsAllNumber(num string) bool {
	for _, char := range num {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
