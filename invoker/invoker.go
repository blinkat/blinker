package invoker

import (
	"fmt"
)

type PublicKey interface {
	Encrypt(b []byte) ([]byte, error)
	GenKey() ([]byte, error)
	//EncryptWithKey(b, k []byte) ([]byte, error)
	//Bytes() ([]byte, error)
	//Encode() ([]byte, error)
}

type PrivateKey interface {
	Public() PublicKey
	Decrypt(a *AsyEncrypted) ([]byte, error)
	//Encode() ([]byte, error)
}

// asymmetric encrypted
type CipherPart struct {
	Iv         JsonBytes `json:"iv,omitempty"`
	Ciphertext JsonBytes `json:"ciphertext,omitempty"`
	Tag        JsonBytes `json:"tag,omitempty"`
}

type AsyEncrypted struct {
	Public       PublicKey
	EncryptedKey []byte
	Part         *CipherPart
	Type         string
}

func (a *AsyEncrypted) Decrypt(prk PrivateKey) ([]byte, error) {
	return prk.Decrypt(a)
}

type GenAsymmetric func(opt ...string) (PrivateKey, error)

var (
	Asymmetrics = make(map[string]GenAsymmetric)
)

const (
	namespace = "blinker/invoker: "
)

// register crypto key
func RegisterAsymmetric(name string, fn GenAsymmetric) error {
	if fn != nil && name != "" {
		Asymmetrics[name] = fn
	}
	return fmt.Errorf(namespace + "params failed.")
}

func AsymmetricKey(name string, opt ...string) (PrivateKey, error) {
	if fn, ok := Asymmetrics[name]; ok {
		return fn(opt...)
	}
	return nil, fmt.Errorf(namespace+"can not found '%s' asymmetric.", name)
}
