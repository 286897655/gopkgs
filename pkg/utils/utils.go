package utils

func BoolFunc(fc func() bool) bool {
	return fc()
}

func OneZeroFunc(fc func() bool) int {
	if !fc() {
		return 0
	}

	return 1
}
