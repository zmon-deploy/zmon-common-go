package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"strings"
)

type AESCipher struct {
	cipher cipher.Block
}

func NewAESCipher(password string) *AESCipher {
	pass := []byte(password)
	secret := make([]byte, 16)
	copy(secret, pass)
	for i := 16; i < len(pass); {
		for j := 0; j < 16 && i < len(pass); j, i = j+1, i+1 {
			secret[j] ^= pass[i]
		}
	}
	cipher, _ := aes.NewCipher(secret)
	return &AESCipher{
		cipher: cipher,
	}
}

func (t *AESCipher) AESEncrypt(srcstr string) string {
	src := []byte(srcstr)
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))

	for bs, be := 0, t.cipher.BlockSize(); bs <= len(src); bs, be = bs+t.cipher.BlockSize(), be+t.cipher.BlockSize() {
		t.cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return strings.ToUpper(hex.EncodeToString(encrypted))
}

func (t *AESCipher) AESDecrypt(hexstr string) string {
	encrypted, _ := hex.DecodeString(hexstr)
	decrypted := make([]byte, len(encrypted))
	for bs, be := 0, t.cipher.BlockSize(); bs < len(encrypted); bs, be = bs+t.cipher.BlockSize(), be+t.cipher.BlockSize() {
		t.cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return string(decrypted[:trim])
}
