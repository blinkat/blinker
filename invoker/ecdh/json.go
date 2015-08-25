package ecdh

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/blinkat/blinker/invoker"
)

type jsonPublic struct {
	X      invoker.JsonBytes `json:"x,omitempty"`
	Y      invoker.JsonBytes `json:"y,omitempty"`
	Curve  int               `json:"curve,omitempty"`
	KeyAlg string            `json:"key-alg,omitempty"`
	EncAlg string            `json:"enc-alg,omitempty"`
	IsComp bool              `json:"is-comp,omitempty"`
	//Type   string            `json:"type,omitempty"`
}

type jsonKey struct {
	Key    invoker.JsonBytes `json:"key,omitempty"`
	Header *jsonPublic       `json:"header,omitempty"`
}

func ParsePublic(b []byte) (invoker.PublicKey, error) {
	jpk := new(jsonPublic)
	err := json.Unmarshal(b, jpk)
	if err != nil {
		return nil, err
	}

	pub := &ecdsa.PublicKey{
		X:     jpk.X.BigInt(),
		Y:     jpk.Y.BigInt(),
		Curve: getCurve(jpk.Curve),
	}
	if pub.Curve == nil {
		return nil, fmt.Errorf(namespace+"unknow curve '%s'", jpk.Curve)
	}

	return &publicKey{
		key: pub,
		params: keyParams{
			key_algor: jpk.KeyAlg,
			enc_algor: jpk.EncAlg,
			curve:     jpk.Curve,
			is_comp:   jpk.IsComp,
		},
	}, nil
}

func parsePublic(b []byte) (invoker.PublicKey, []byte, error) {
	key := new(jsonKey)
	err := json.Unmarshal(b, key)
	if err != nil {
		return nil, nil, err
	}

	pub := &ecdsa.PublicKey{
		X:     key.Header.X.BigInt(),
		Y:     key.Header.Y.BigInt(),
		Curve: getCurve(key.Header.Curve),
	}
	if pub.Curve == nil {
		return nil, nil, fmt.Errorf(namespace+"unknow curve '%s'", key.Header.Curve)
	}

	return &publicKey{
		key: pub,
		params: keyParams{
			key_algor: key.Header.KeyAlg,
			enc_algor: key.Header.EncAlg,
			curve:     key.Header.Curve,
			is_comp:   key.Header.IsComp,
		},
	}, []byte(key.Key), nil
}

func ParseEncrypted(jpk *invoker.AsymmetricsJson) (*invoker.AsyEncrypted, error) {
	key, code, err := parsePublic(jpk.Key)
	if err != nil {
		return nil, err
	}

	return &invoker.AsyEncrypted{
		Public:       key,
		EncryptedKey: code,
		Part:         jpk.Part,
		Type:         jpk.Type,
	}, nil
}

func init() {
	err := invoker.RegisterAsymmetric("ecdh", GenerateKey, ParseEncrypted, ParsePublic)
	if err != nil {
		panic(err)
	}
}
