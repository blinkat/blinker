package invoker

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
)

type Crypter interface {
	AuthData() []byte
	EncryptKey(plaintext []byte) ([]byte, error)
	DecryptKey(ciphertext []byte) ([]byte, error)
	GenKey() ([]byte, error)
	Algorithm() EncAlgorithm
	//Identify() string
}

type json_ciphertext struct {
	Key    []byte `json:"key,omitempty"`
	IV     []byte `json:"iv,omitempty"`
	Cipher []byte `json:"ciphertext,omitempty"`
	Tag    []byte `json:"tag,omitempty"`
}

type Invoker struct {
	crypter Crypter
}

func (c *Invoker) Encrypt(plaintext []byte) ([]byte, error) {
	var err error
	switch CompressAlgorithm {
	case COMPRESS_DEF:
		plaintext, err = deflate(plaintext)
		if err != nil {
			return nil, err
		}
	}

	cek, err := c.crypter.GenKey()
	if err != nil {
		return nil, err
	}

	key, err := c.crypter.EncryptKey(cek)
	if err != nil {
		return nil, err
	}

	ad := get_cont_cipher(c.crypter.Algorithm())
	if ad == nil {
		return nil, fmt.Errorf(namespace + ": encrypt failed.")
	}
	auth := c.crypter.AuthData()
	part, err := ad.Encrypt(cek, auth, plaintext)
	j := &json_ciphertext{
		Key:    key,
		Cipher: part.cip,
		IV:     part.iv,
		Tag:    part.tag,
	}

	ciphertext, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func (c *Invoker) Decrypt(ciphertext []byte) ([]byte, error) {
	var err error

	j := &json_ciphertext{}
	err = json.Unmarshal(ciphertext, j)
	if err != nil {
		return nil, err
	}

	cek, err := c.crypter.DecryptKey(j.Key)
	auth := c.crypter.AuthData()
	ad := get_cont_cipher(c.crypter.Algorithm())
	if ad == nil {
		return nil, fmt.Errorf(namespace + ": encrypt failed.")
	}

	part := &AeadPart{
		iv:  j.IV,
		cip: j.Cipher,
		tag: j.Tag,
	}

	plaintext, err := ad.Decrypt(cek, auth, part)
	if err != nil {
		return nil, err
	}

	switch CompressAlgorithm {
	case COMPRESS_DEF:
		plaintext, err = inflate(plaintext)
	}
	return plaintext, err
}

func NewInvoker(key Crypter) *Invoker {
	c := &Invoker{
		crypter: key,
	}

	return c
}

func deflate(src []byte) ([]byte, error) {
	out := new(bytes.Buffer)
	writer, _ := flate.NewWriter(out, 1)
	io.Copy(writer, bytes.NewBuffer(src))
	err := writer.Close()
	return out.Bytes(), err
}

func inflate(src []byte) ([]byte, error) {
	out := new(bytes.Buffer)
	reader := flate.NewReader(bytes.NewBuffer(src))
	_, err := io.Copy(out, reader)
	if err != nil {
		return nil, err
	}
	err = reader.Close()
	return out.Bytes(), err
}
