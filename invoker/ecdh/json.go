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
	Curve  string            `json:"curve,omitempty"`
	KeyAlg string            `json:"key-alg,omitempty"`
	EncAlg string            `json:"enc-alg,omitempty"`
	IsComp bool              `json:"is-comp,omitempty"`
}

type jsonKey struct {
	Key    invoker.JsonBytes `json:"key,omitempty"`
	Header *jsonPublic       `json:"header,omitempty"`
}

func ParsePublic(b []byte) (invoker.PublicKey, error) {
	k, _, err := parsePublic(b)
	return k, err
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

func ParseEncrypted(b []byte) (*invoker.AsyEncrypted, error) {
	jpk := new(invoker.AsymmetricsJson)
	err := json.Unmarshal(b, jpk)
	if err != nil {
		return nil, err
	}

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
