package invoker

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type AeadPart struct {
	iv  []byte
	cip []byte
	tag []byte
}

type aead_cont struct {
	key_bytes  int
	auth_bytes int
	get_aead   func(key []byte) (cipher.AEAD, error)
}

func (a *aead_cont) Encrypt(key, aad, pt []byte) (*AeadPart, error) {
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

	return &AeadPart{
		iv:  iv,
		cip: cip[:oft],
		tag: cip[oft:],
	}, nil
}

func (a *aead_cont) Decrypt(key, aad []byte, part *AeadPart) ([]byte, error) {
	aead, err := a.get_aead(key)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, part.iv, append(part.cip, part.tag...), aad)
}

func new_aes_gcm(size int) *aead_cont {
	return &aead_cont{
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

func get_cont_cipher(alg EncAlgorithm) *aead_cont {
	switch alg {
	case A128GCM:
		return new_aes_gcm(16)
	case A192GCM:
		return new_aes_gcm(24)
	case A256GCM:
		return new_aes_gcm(32)
	}
	return nil
}
