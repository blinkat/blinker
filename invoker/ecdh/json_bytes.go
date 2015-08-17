package ecdh

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"strings"
)

type json_bytes []byte

func (j json_bytes) ToBase64() string {
	ret := base64.URLEncoding.EncodeToString([]byte(j))
	return strings.TrimRight(ret, "=")
}

func (j *json_bytes) FromBase64(s string) error {
	missing := (4 - len(s)%4) % 4
	s += strings.Repeat("=", missing)
	ret, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	*j = ret
	return nil
}

func (j json_bytes) BigInt() *big.Int {
	return new(big.Int).SetBytes([]byte(j))
}

func (j *json_bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.ToBase64())
}

func (j *json_bytes) UnmarshalJSON(data []byte) error {
	var code string
	err := json.Unmarshal(data, &code)
	if err != nil {
		return err
	}

	if code == "" {
		return nil
	}

	return j.FromBase64(code)
}
