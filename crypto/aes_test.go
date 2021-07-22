package crypto

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCrypto(t *testing.T) {
	crypto := NewAESCipher("abcdef123456")

	encrypted, err := crypto.AESEncrypt("helloWorld")
	require.NoError(t, err)

	decrypted, err := crypto.AESDecrypt(encrypted)
	require.NoError(t, err)

	require.Equal(t, "helloWorld", decrypted)
}

func TestWrongEncrypted(t *testing.T) {
	wrongEncrypted := "ca8e81c53958038846775d6f00c46df086e601d33bd0af78e1b4d396e90e42afe330e53787ea0832e3562e7e718feb71"
	_, err := NewAESCipher("hseRTo5bUFhdeI9W").AESDecrypt(wrongEncrypted)
	require.Error(t, err)
}
