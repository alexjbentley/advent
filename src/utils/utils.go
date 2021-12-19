package utils

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func GetMask(length int) (mask int) {
	mask = 0b1
	for i := 0; i < length; i++ {
		mask = mask | 0b1<<i
	}
	return
}
