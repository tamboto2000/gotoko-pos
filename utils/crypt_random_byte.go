package utils

import (
	"crypto/rand"
	"math/big"
)

// CryptRandString generate cryptographically secure
// random ASCII characters with n length
func CryptRandString(n int) string {
	ret := make([]byte, n)
	min := 33
	max := 126
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
		ret[i] = byte(num.Int64() + int64(min))
	}

	return string(ret)
}
