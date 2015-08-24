package ecdh

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/blinkat/blinker/invoker"
	"io"
)

// type aeadPart struct {
// 	iv  []byte
// 	cip []byte
// 	tag []byte
// }

type aeadCont struct {
	key_bytes  int
	auth_bytes int
	get_aead   func(key []byte) (cipher.AEAD, error)
}

func (a *aeadCont) encrypt(key, aad, pt []byte) (*invoker.CipherPart, error) {
	aead, err := a.get_aead(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aead.NonceSize())
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	cip := aead.Seal(nil, iv, pt, aad)
	oft := len(cip) - a.auth_bytes

	return &invoker.CipherPart{
		Iv:         iv,
		Ciphertext: cip[:oft],
		Tag:        cip[oft:],
	}, nil
}

func (a *aeadCont) decrypt(key, aad []byte, part *invoker.CipherPart) ([]byte, error) {
	aead, err := a.get_aead(key)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, part.Iv, append(part.Ciphertext, part.Tag...), aad)
}

func new_aes_gcm(size int) *aeadCont {
	return &aeadCont{
		key_bytes:  size,
		auth_bytes: 16,
		get_aead: func(key []byte) (cipher.AEAD, error) {
			a, err := aes.NewCipher(key)
			if err != nil {
				return nil, err
			}
			return cipher.NewGCM(a)
		},
	}
}

func getContCipher(alg string) *aeadCont {
	switch alg {
	case ENC_A128GCM:
		return new_aes_gcm(16)
	case ENC_A192GCM:
		return new_aes_gcm(24)
	case ENC_A256GCM:
		return new_aes_gcm(32)
	}
	return nil
}
