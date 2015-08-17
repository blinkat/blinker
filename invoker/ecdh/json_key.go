package ecdh

type json_puk struct {
	Key   json_bytes `json:"k,omitempty"`
	X     json_bytes `json:"x,omitempty"`
	Y     json_bytes `json:"y,omitempty"`
	Curve int        `json:"c,omitempty"`
}
