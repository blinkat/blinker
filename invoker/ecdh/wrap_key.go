package ecdh

import (
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"errors"
)

var default_iv = []byte{0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6}

func key_wrap(block cipher.Block, cek []byte) ([]byte, error) {
	if len(cek)%8 != 0 {
		return nil, errors.New(namespace + ": key wrap input must be 8 byte blocks")
	}

	n := len(cek) / 8
	r := make([][]byte, n)

	for i := range r {
		r[i] = make([]byte, 8)
		copy(r[i], cek[i*8:])
	}

	buf := make([]byte, 16)
	tbs := make([]byte, 8)
	copy(buf, default_iv)

	for t := 0; t < 6*n; t++ {
		copy(buf[8:], r[t%n])

		block.Encrypt(buf, buf)
		binary.BigEndian.PutUint64(tbs, uint64(t+1))

		for i := 0; i < 8; i++ {
			buf[i] = buf[i] ^ tbs[i]
		}
		copy(r[t%n], buf[8:])
	}

	out := make([]byte, (n+1)*8)
	copy(out, buf[:8])
	for i := range r {
		copy(out[(i+1)*8:], r[i])
	}

	return out, nil
}

func key_unwarp(block cipher.Block, ciphertext []byte) ([]byte, error) {
	if len(ciphertext)%8 != 0 {
		return nil, errors.New(namespace + ": key wrap input must be 8 byte blocks")
	}

	n := (len(ciphertext) / 8) - 1
	r := make([][]byte, n)

	for i := range r {
		r[i] = make([]byte, 8)
		copy(r[i], ciphertext[(i+1)*8:])
	}

	buffer := make([]byte, 16)
	tBytes := make([]byte, 8)
	copy(buffer[:8], ciphertext[:8])

	for t := 6*n - 1; t >= 0; t-- {
		binary.BigEndian.PutUint64(tBytes, uint64(t+1))

		for i := 0; i < 8; i++ {
			buffer[i] = buffer[i] ^ tBytes[i]
		}
		copy(buffer[8:], r[t%n])

		block.Decrypt(buffer, buffer)

		copy(r[t%n], buffer[8:])
	}

	if subtle.ConstantTimeCompare(buffer[:8], default_iv) == 0 {
		return nil, errors.New(namespace + ": failed to unwrap key")
	}

	out := make([]byte, n*8)
	for i := range r {
		copy(out[i*8:], r[i])
	}

	return out, nil
}
