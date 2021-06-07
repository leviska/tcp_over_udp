package util

import "math/rand"

func Truncate(s string, i int) string {
	if len(s) < i {
		return s
	}
	return s[:i] + "...(truncated)"
}

func RandomString(size int) string {
	res := make([]byte, size)
	for i := range res {
		res[i] = byte(rand.Intn(26)) + 'a'
	}
	return string(res)
}
