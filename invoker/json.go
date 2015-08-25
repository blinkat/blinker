package invoker

type AsymmetricsJson struct {
	Key  JsonBytes   `json:"key,omitempty"`
	Part *CipherPart `json:"cipher,omitempty"`
	Type string      `json:"type,omitempty"`
}

type AsymmetricsPublic struct {
	Key  JsonBytes `json:"key"`
	Type string    `json:"type"`
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
