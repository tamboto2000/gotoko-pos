package aescrypt

import (
	"os"

	"github.com/tamboto2000/gotoko-pos/utils"
)

func Encrypt(data []byte) ([]byte, error) {
	key := os.Getenv("AES_CRYPT_KEY")
	return utils.AESCrypt([]byte(key), data)
}

func Decrypt(chiper []byte) ([]byte, error) {
	key := os.Getenv("AES_CRYPT_KEY")
	return utils.AESDecrypt([]byte(key), chiper)
}
