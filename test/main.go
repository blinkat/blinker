package main

import (
	"fmt"
	//"github.com/blinkat/blinker/invoker"
	"github.com/blinkat/blinker/invoker/ecdh"
)

func main() {
	key, err := ecdh.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}

	plaintext := []byte("what fuck.")

	puk := ecdh.NewPublicKey(ecdh.CURVE_P521, ecdh.KEY_A256K, ecdh.ENC_A256GCM)
	cip, err := puk.Encrypt(plaintext)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(cip))

	enc, err := ecdh.ParseEncrypted(cip)
	if err != nil {
		fmt.Println(err)
	}

	pl, err := enc.Decrypt(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(pl))
	// cip, err := puk.Encrypt(plaintext)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(cip))

	// enc, err := ecdh.ParseEncrypted(cip)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// pla, err := enc.Decrypt(key)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(pla))

	// text, err := puk.GenKey()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// puk, err = ecdh.ParsePublic(text)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// cipher, err := puk.Encrypt(plaintext)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(string(cipher))

	// enc, err := ecdh.ParseEncrypted(cipher)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// plain, err := enc.Decrypt(key)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(plain)
}
