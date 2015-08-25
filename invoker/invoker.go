package invoker

import (
	"encoding/json"
	"fmt"
)

type PublicKey interface {
	Encrypt(b []byte) ([]byte, error)
	Bytes() ([]byte, error)
}

type PrivateKey interface {
	Public() PublicKey
	Decrypt(a *AsyEncrypted) ([]byte, error)
}

// ========== asymmetric encryption maps ============
type GenAsymmetric func(opt ...int) (PrivateKey, error)               //gen asymmetric key
type EncryptedParse func(enc *AsymmetricsJson) (*AsyEncrypted, error) //encrypted parse func
type PublicParse func(enc []byte) (PublicKey, error)                  //asymmetric public key parser

type algorithm struct {
	gen    GenAsymmetric
	enc    EncryptedParse
	pub    PublicParse
	isAsym bool // is asymmetric algorithm
}

var (
	asymmetrics = make(map[string]*algorithm)
)

const (
	namespace = "blinker/invoker: "
)

// register crypto key
func RegisterAsymmetric(name string, gen GenAsymmetric, enc EncryptedParse, pub PublicParse) error {
	if gen != nil && name != "" {
		asymmetrics[name] = &algorithm{
			gen:    gen,
			enc:    enc,
			pub:    pub,
			isAsym: true,
		}
		return nil
	}
	return fmt.Errorf(namespace + "name and gen can't be null.")
}

// generate asymmetric key
func GenAsymKey(name string, opt ...int) (PrivateKey, error) {
	if alg, ok := asymmetrics[name]; ok {
		if alg.isAsym {
			return alg.gen(opt...)
		}
	}
	return nil, fmt.Errorf(namespace+"unknow '%s' algorithm.", name)
}

func ParseEncrypted(b []byte) (*AsyEncrypted, error) {
	jpk := new(AsymmetricsJson)
	err := json.Unmarshal(b, jpk)
	if err != nil {
		return nil, err
	}

	if alg, ok := asymmetrics[jpk.Type]; ok && alg.isAsym {
		return alg.enc(jpk)
	}
	return nil, fmt.Errorf(namespace + "parse key failed.")
}

func ParsePublic(b []byte) (PublicKey, error) {
	jpk := new(AsymmetricsPublic)
	err := json.Unmarshal(b, jpk)
	if err != nil {
		return nil, err
	}

	if alg, ok := asymmetrics[jpk.Type]; ok && alg.isAsym {
		return alg.pub([]byte(jpk.Key))
	}
	return nil, fmt.Errorf(namespace + "parse key failed.")
}
