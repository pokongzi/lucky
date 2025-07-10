package aes_crypto

import (
	"testing"
)

func TestAes(t *testing.T) {
	p := []byte("plaintextdddffdfdfddhghbvbnvnvnvnmyiiuooscvopcovkkrtrvnvbbmkmqscxbsfddfd ff ")
	key := []byte("12345678abcdefgh")
	t.Log(len(p))
	ciphertext, _ := Base64AESCBCEncrypt(p, key) // A67NhD3RBiNaMgG6HTm8LQ==
	t.Log(ciphertext, len(ciphertext))
	plaintext, _ := Base64AESCBCDecrypt(ciphertext, key) // plaintext
	t.Log(string(plaintext))

}
