package pperf

import (
	"math/rand"
)

func randbytes(l int) []byte {
	buf := make([]byte, l)
	rand.Read(buf)
	return buf
}
