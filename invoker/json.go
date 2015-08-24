package invoker

// type AsymmetricsCipher struct {
// 	IV         JsonBytes `json:"iv,omitempty"`
// 	Ciphertext JsonBytes `json:"ciphertext,omitempty"`
// 	Tag        JsonBytes `json:"tag,omitempty"`
// }

type AsymmetricsJson struct {
	//Key        JsonBytes `json:"key,omitempty"`
	Key  JsonBytes   `json:"key,omitempty"`
	Part *CipherPart `json:"cipher,omitempty"`
	Type string      `json:"type,omitempty"`
}
