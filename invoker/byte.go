package invoker

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"strings"
)

type JsonBytes []byte

func (j JsonBytes) ToBase64() string {
	ret := base64.URLEncoding.EncodeToString([]byte(j))
	return strings.TrimRight(ret, "=")
}

func (j *JsonBytes) FromBase64(s string) error {
	missing := (4 - len(s)%4) % 4
	s += strings.Repeat("=", missing)
	ret, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	*j = ret
	return nil
}

func (j JsonBytes) BigInt() *big.Int {
	return new(big.Int).SetBytes([]byte(j))
}

func (j *JsonBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.ToBase64())
}

func (j *JsonBytes) UnmarshalJSON(data []byte) error {
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

func Deflate(src []byte) ([]byte, error) {
	out := new(bytes.Buffer)
	writer, _ := flate.NewWriter(out, 1)
	io.Copy(writer, bytes.NewBuffer(src))
	err := writer.Close()
	return out.Bytes(), err
}

func Inflate(src []byte) ([]byte, error) {
	out := new(bytes.Buffer)
	reader := flate.NewReader(bytes.NewBuffer(src))
	_, err := io.Copy(out, reader)
	if err != nil {
		return nil, err
	}
	err = reader.Close()
	return out.Bytes(), err
}
