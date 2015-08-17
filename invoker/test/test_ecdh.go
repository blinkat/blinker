package test

import (
	"fmt"
	"github.com/blinkat/blinker/invoker"
	"github.com/blinkat/blinker/invoker/ecdh"
	"time"
)

func TestECDH() {
	key, err := ecdh.GenerateKey(ecdh.KEY_A256K, invoker.A256GCM, ecdh.Curve_P521)
	if err != nil {
		fmt.Println(err)
	}
	crypter := invoker.NewInvoker(key)
	plaintext := "holy shit."

	now := time.Now()
	ciphertext, err := crypter.Encrypt([]byte(plaintext))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("encrypt use time:", (time.Now().UnixNano()-now.UnixNano())/int64(time.Millisecond))
	//fmt.Println(string(ciphertext))

	now = time.Now()
	plt, err := crypter.Decrypt(ciphertext)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("decrypt use time:", (time.Now().UnixNano()-now.UnixNano())/int64(time.Millisecond))

	fmt.Println(string(plt))
}
