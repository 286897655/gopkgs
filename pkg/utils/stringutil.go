package utils

import (
	"strings"
)

func StringIsNullOrEmpty(str string) bool {
	return 0 == len(strings.TrimSpace(str))
}

func StringsHasNullOrEmpty(args ...string) bool {
	for _, arg := range args {
		if StringIsNullOrEmpty(arg) {
			return true
		}
	}

	return false
}

func StringsEqual(str1, str2 string) bool {
	return str1 == str2
}
