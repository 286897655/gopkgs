package utils

import "strings"

func StringIsNullOrEmpty(str string) bool {
	return 0 == len(strings.TrimSpace(str))
}

func StringsEqual(str1, str2 string) bool {
	return str1 == str2
}
