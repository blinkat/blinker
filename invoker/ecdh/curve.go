package ecdh

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	_ "crypto/sha256"
	"encoding/binary"
	"github.com/blinkat/blinker/invoker/concat"
)

func GetCurve(c CurveAlgorithm) elliptic.Curve {
	switch c {
	case Curve_P224:
		return elliptic.P224()
	case Curve_P256:
		return elliptic.P256()
	case Curve_P384:
		return elliptic.P384()
	case Curve_P521:
		return elliptic.P521()
	}
	return nil
}

func DeriveECDH(alg string, apu, apv []byte, pri *ecdsa.PrivateKey, pub *ecdsa.PublicKey, size int) []byte {
	alg_ := length_prefixed([]byte(alg))
	pt_u := length_prefixed(apu)
	pt_v := length_prefixed(apv)

	sup := make([]byte, 4)
	binary.BigEndian.PutUint32(sup, uint32(size)*8)

	z, _ := pri.PublicKey.Curve.ScalarMult(pub.X, pub.Y, pri.D.Bytes())
	reader := concat.NewConcat(crypto.SHA256, z.Bytes(), alg_, pt_u, pt_v, sup, []byte{})

	key := make([]byte, size)
	reader.Read(key)

	return key
}

func length_prefixed(data []byte) []byte {
	out := make([]byte, len(data)+4)
	binary.BigEndian.PutUint32(out, uint32(len(data)))
	copy(out[4:], data)
	return out
}
