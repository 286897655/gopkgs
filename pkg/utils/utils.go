package utils

func BoolFunc2Bool(fc func() bool) bool {
	return fc()
}

func BoolFunc2Str(fc func() bool, trueStr, falseStr string) string {
	if fc() {
		return trueStr
	}

	return falseStr
}

func BoolFunc2Int(fc func() bool, trueInt, falseInt int) int {
	if fc() {
		return trueInt
	}

	return falseInt
}

func BoolFunc2OneZero(fc func() bool) int {
	if !fc() {
		return 0
	}

	return 1
}
