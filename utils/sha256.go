package utils

import (
	"bytes"
	"crypto/sha256"
)

func SHA256Hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func SHA256HashEqual(hash []byte, data []byte) bool {
	hashed := SHA256Hash(data)
	return bytes.Equal(hash, hashed)
}
