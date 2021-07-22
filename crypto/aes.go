package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type AESCipher struct {
	secretKey string
}

func NewAESCipher(secretKey string) *AESCipher {
	return &AESCipher{
		secretKey: secretKey,
	}
}

func (t *AESCipher) AESEncrypt(srcstr string) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: r")
		}
	}()

	cipherBlock, err := getCipher(t.secretKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to get cipher")
	}

	src := []byte(srcstr)
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))

	for bs, be := 0, cipherBlock.BlockSize(); bs <= len(src); bs, be = bs+cipherBlock.BlockSize(), be+cipherBlock.BlockSize() {
		cipherBlock.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return strings.ToUpper(hex.EncodeToString(encrypted)), nil
}

func (t *AESCipher) AESDecrypt(hexstr string) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: r")
		}
	}()

	cipherBlock, err := getCipher(t.secretKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to get cipher")
	}

	encrypted, err := hex.DecodeString(hexstr)
	if err != nil {
		return "", err
	}

	decrypted := make([]byte, len(encrypted))
	for bs, be := 0, cipherBlock.BlockSize(); bs < len(encrypted); bs, be = bs+cipherBlock.BlockSize(), be+cipherBlock.BlockSize() {
		cipherBlock.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return string(decrypted[:trim]), nil
}

func getCipher(secretKey string) (cipher.Block, error) {
	pass := []byte(secretKey)
	secret := make([]byte, 16)
	copy(secret, pass)
	for i := 16; i < len(pass); {
		for j := 0; j < 16 && i < len(pass); j, i = j+1, i+1 {
			secret[j] ^= pass[i]
		}
	}
	return aes.NewCipher(secret)
}