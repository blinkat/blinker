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

const (
	KEY_NONE  = "NONE"
	KEY_A128K = "A128K"
	KEY_A192K = "A192K"
	KEY_A256K = "A256K"

	ENC_A128GCM = "A128GCM"
	ENC_A192GCM = "A192GCM"
	ENC_A256GCM = "A256GCM"
)

type publicKey struct {
	key    *ecdsa.PublicKey
	params keyParams
}

type privateKey struct {
	key    *ecdsa.PrivateKey
	puk    *publicKey
	params keyParams
}

type keyParams struct {
	key_algor string
	enc_algor string
	curve     int
	is_comp   bool
}

func (k *keyParams) keySize() int {
	switch k.enc_algor {
	case ENC_A128GCM:
		return 16
	case ENC_A192GCM:
		return 24
	case ENC_A256GCM:
		return 32
	}
	return -1
}

// generate private key
// opt 0 = enc 1 = curve 2 = key 3 = comp
func GenerateKey(opt ...int) (invoker.PrivateKey, error) {
	p := keyParams{
		key_algor: KEY_A256K,
		enc_algor: ENC_A256GCM,
		curve:     CURVE_P521,
		is_comp:   false,
	}
	leng := len(opt)
	if leng >= 1 {
		switch opt[0] {
		case 128:
			p.enc_algor = ENC_A128GCM
		case 192:
			p.enc_algor = ENC_A192GCM
		case 256:
			p.enc_algor = ENC_A256GCM
		default:
			return nil, fmt.Errorf(namespace+"unknow #%d encryption algorithm.", opt[0])
		}
	}
	if leng >= 2 {
		p.curve = opt[1]
	}
	if leng >= 3 {
		switch opt[2] {
		case 0:
			p.key_algor = KEY_NONE
		case 128:
			p.key_algor = KEY_A128K
		case 192:
			p.key_algor = KEY_A192K
		case 256:
			p.key_algor = KEY_A256K
		}
	}
	if leng >= 4 {
		p.is_comp = opt[3] != 0
	}

	curve := getCurve(p.curve)
	if curve == nil {
		return nil, fmt.Errorf(namespace+"unknow curve #%d.", p.curve)
	}
	ky, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &privateKey{
		params: p,
		key:    ky,
		puk: &publicKey{
			params: p,
			key:    &ky.PublicKey,
		},
	}, nil
}

// ======== public key ========
func (p *publicKey) genKey() ([]byte, *ecdsa.PublicKey, error) {
	pri, err := ecdsa.GenerateKey(p.key.Curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	out := deriveECDH(p.params.key_algor, []byte{}, []byte{}, pri, p.key, p.params.keySize())
	return out, &pri.PublicKey, nil
}

func (p *publicKey) encryptKey(cek []byte) (*jsonKey, []byte, error) {
	if p.params.key_algor == KEY_NONE {
		return &jsonKey{}, nil, nil
	}

	kek, pub, err := p.genKey()
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, nil, err
	}

	jek, err := key_wrap(block, cek)
	if err != nil {
		return nil, nil, err
	}

	return &jsonKey{
		Header: &jsonPublic{
			X:      invoker.JsonBytes(pub.X.Bytes()),
			Y:      invoker.JsonBytes(pub.Y.Bytes()),
			Curve:  p.params.curve,
			KeyAlg: p.params.key_algor,
			EncAlg: p.params.enc_algor,
			IsComp: p.params.is_comp,
		},
		Key: jek,
	}, authData(pub), nil
}

func authData(p *ecdsa.PublicKey) []byte {
	jpk := &jsonPublic{
		X: invoker.JsonBytes(p.X.Bytes()),
		Y: invoker.JsonBytes(p.Y.Bytes()),
	}
	j, _ := json.Marshal(jpk)
	return j
}

// gen rand key
func (p *publicKey) randGenKey() ([]byte, error) {
	key := make([]byte, p.params.keySize())
	_, err := io.ReadFull(rand.Reader, key)
	return key, err
}

func (p *publicKey) Encrypt(b []byte) ([]byte, error) {
	cek, err := p.randGenKey()
	if err != nil {
		return nil, err
	}
	key, auth, err := p.encryptKey(cek)
	if err != nil {
		return nil, err
	}
	return p.encryptWithKey(b, cek, key, auth)
}

func (p *publicKey) Bytes() ([]byte, error) {
	jpk := &jsonPublic{
		X:      invoker.JsonBytes(p.key.X.Bytes()),
		Y:      invoker.JsonBytes(p.key.Y.Bytes()),
		Curve:  p.params.curve,
		KeyAlg: p.params.key_algor,
		EncAlg: p.params.enc_algor,
		IsComp: p.params.is_comp,
	}

	js, err := json.Marshal(jpk)
	if err != nil {
		return nil, err
	}

	ret := &invoker.AsymmetricsPublic{
		Key:  invoker.JsonBytes(js),
		Type: "ecdh",
	}
	return json.Marshal(ret)
}

func (p *publicKey) encryptWithKey(plaintext, cek []byte, key *jsonKey, auth []byte) ([]byte, error) {
	var err error
	if p.params.is_comp {
		plaintext, err = invoker.Deflate(plaintext)
		if err != nil {
			return nil, err
		}
	}

	ad := getContCipher(p.params.enc_algor)
	if ad == nil {
		return nil, fmt.Errorf(namespace + "get aead failed")
	}

	//auth := pub.authData()
	part, err := ad.encrypt(cek, auth, plaintext)

	kj, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}

	j := &invoker.AsymmetricsJson{
		Key:  invoker.JsonBytes(kj),
		Part: part,
		Type: "ecdh",
	}
	return json.Marshal(j)
}

// ========= private key ==========
func (p *privateKey) Public() invoker.PublicKey {
	return p.puk
}

func (p *privateKey) decryptKey(encrypted []byte, pub *publicKey) ([]byte, error) {
	size := 0
	switch p.params.key_algor {
	case KEY_NONE:
		return deriveECDH(p.params.enc_algor, []byte{}, []byte{}, p.key, pub.key, p.params.keySize()), nil
	case KEY_A128K:
		size = 16
	case KEY_A192K:
		size = 24
	case KEY_A256K:
		size = 32
	default:
		return nil, fmt.Errorf(namespace+"unknow key. '%s'", p.params.key_algor)
	}

	key := deriveECDH(p.params.key_algor, []byte{}, []byte{}, p.key, pub.key, size)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return key_unwarp(block, encrypted)
}

func (p *privateKey) Decrypt(a *invoker.AsyEncrypted) ([]byte, error) {
	var puk *publicKey
	switch a.Public.(type) {
	case *publicKey:
		puk = a.Public.(*publicKey)
	default:
		return nil, fmt.Errorf(namespace + "unknow key.")
	}

	cek, err := p.decryptKey(a.EncryptedKey, a.Public.(*publicKey))
	if err != nil {
		return nil, err
	}
	auth := authData(puk.key)
	ad := getContCipher(p.params.enc_algor)
	if ad == nil {
		return nil, fmt.Errorf(namespace + "unknow enc.")
	}

	plaintext, err := ad.decrypt(cek, auth, a.Part)
	if err != nil {
		return nil, err
	}

	if p.params.is_comp {
		return invoker.Inflate(plaintext)
	}
	return plaintext, err
}
