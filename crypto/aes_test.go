package crypto

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCrypto(t *testing.T) {
	crypto := NewAESCipher("abcdef123456")
	encrypted := crypto.AESEncrypt("helloWorld")
	decrypted := crypto.AESDecrypt(encrypted)
	require.Equal(t, "helloWorld", decrypted)
}
