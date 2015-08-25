package main

import (
	"fmt"
	"github.com/blinkat/blinker/invoker"
	_ "github.com/blinkat/blinker/invoker/ecdh"
)

func main() {
	key, err := invoker.GenAsymKey("ecdh")
	if err != nil {
		fmt.Println(err)
	}

	plaintext := []byte("what fuck.")

	pukBits, err := key.Public().Bytes()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(pukBits))

	puk, err := invoker.ParsePublic(pukBits)
	if err != nil {
		fmt.Println(err)
	}

	cip, err := puk.Encrypt(plaintext)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("ciphertext:", string(cip))

	enc, err := invoker.ParseEncrypted(cip)
	if err != nil {
		fmt.Println(err)
	}

	pl, err := enc.Decrypt(key)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("plaintext:", string(pl))
}
