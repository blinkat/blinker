package ecdh

import (
	"crypto/aes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/blinkat/blinker/invoker"
	"io"
)

const namespace = "blinker/invoker/ecdh"

// encryption key algorithms
type KeyAlgorithm int
type CurveAlgorithm int

const (
	KEY_NONE = KeyAlgorithm(iota)
	KEY_A128K
	KEY_A192K
	KEY_A256K
)

const (
	Curve_P224 = CurveAlgorithm(iota)
	Curve_P256
	Curve_P384
	Curve_P521
)

type PublicKey struct {
	key    *ecdsa.PublicKey
	params *key_params
}

type PrivateKey struct {
	key    *ecdsa.PrivateKey
	pub    *PublicKey
	params *key_params
}

type key_params struct {
	key_algor KeyAlgorithm
	enc_algor invoker.EncAlgorithm
	curve     CurveAlgorithm
}

// public
func (p *PublicKey) GenKey() ([]byte, *ecdsa.PublicKey, error) {
	pri, err := ecdsa.GenerateKey(p.key.Curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	out := DeriveECDH(string(p.params.key_algor), []byte{}, []byte{}, pri, p.key, p.params.KeySize())
	return out, &pri.PublicKey, nil
}

func (p *PublicKey) EncryptKey(cek []byte) ([]byte, error) {
	kek, pub, err := p.GenKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}

	jek, err := key_wrap(block, cek)
	if err != nil {
		return nil, err
	}

	jpk := &json_puk{
		X:     json_bytes(pub.X.Bytes()),
		Y:     json_bytes(pub.Y.Bytes()),
		Curve: int(p.params.curve),
		Key:   jek,
	}

	bys, err := json.Marshal(jpk)
	if err != nil {
		return nil, err
	}
	return bys, nil
}

// private
func (p *PrivateKey) GenKey() ([]byte, error) {
	// ret, _, err := p.Public().GenKey()
	// return ret, err
	key := make([]byte, p.params.KeySize())
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (p *PrivateKey) Algorithm() invoker.EncAlgorithm {
	return p.params.enc_algor
}

func (p *PrivateKey) AuthData() []byte {
	jpk := &json_puk{
		X:     json_bytes(p.pub.key.X.Bytes()),
		Y:     json_bytes(p.pub.key.Y.Bytes()),
		Curve: int(p.params.curve),
		Key:   json_bytes{},
	}
	j, _ := json.Marshal(jpk)
	return j
}

func (p *PrivateKey) Public() *PublicKey {
	return p.pub
}

func (p *PrivateKey) EncryptKey(cek []byte) ([]byte, error) {
	return p.pub.EncryptKey(cek)
}

func (p *PrivateKey) DecryptKey(cek []byte) ([]byte, error) {
	jpk := &json_puk{}
	err := json.Unmarshal(cek, jpk)
	if err != nil {
		return nil, err
	}

	pub := &ecdsa.PublicKey{
		X:     jpk.X.BigInt(),
		Y:     jpk.Y.BigInt(),
		Curve: GetCurve(CurveAlgorithm(jpk.Curve)),
	}
	if pub.Curve == nil {
		return nil, err
	}

	apu := []byte{}
	apv := []byte{}

	var ksize int
	switch p.params.key_algor {
	case KEY_NONE:
		return DeriveECDH(string(p.params.enc_algor), apu, apv, p.key, pub, p.params.KeySize()), nil
	case KEY_A128K:
		ksize = 16
	case KEY_A192K:
		ksize = 24
	case KEY_A256K:
		ksize = 32
	default:
		return nil, fmt.Errorf(namespace + ": key size error")
	}

	key := DeriveECDH(string(p.params.key_algor), apu, apv, p.key, pub, ksize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return key_unwarp(block, []byte(jpk.Key))
}

// params
func (p *key_params) KeySize() int {
	switch p.enc_algor {
	case invoker.A128GCM:
		return 16
	case invoker.A192GCM:
		return 24
	case invoker.A256GCM:
		return 32
	}
	return -1
}

// generate key
// k = KeyAlgorithm
// e = EncAlgorithm
// c = CurveAlgorithm
func GenerateKey(k KeyAlgorithm, e invoker.EncAlgorithm, c CurveAlgorithm) (*PrivateKey, error) {
	cur := GetCurve(c)
	if cur == nil {
		return nil, fmt.Errorf(namespace+": unknow curve #%d.", int(c))
	}
	key, err := ecdsa.GenerateKey(cur, rand.Reader)
	if err != nil {
		return nil, err
	}

	params := &key_params{
		key_algor: k,
		enc_algor: e,
		curve:     c,
	}

	ret := &PrivateKey{
		key: key,
		pub: &PublicKey{
			key:    &key.PublicKey,
			params: params,
		},
		params: params,
	}

	return ret, nil
}
