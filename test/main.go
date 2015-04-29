package main

import (
	"fmt"
	"github.com/blinkat/blinks/strike"
)

type test_s struct {
	M map[int]float64
	V string
}

func main() {
	t := &test_s{
		M: make(map[int]float64),
		V: "nothing",
	}
	t.M[1] = 123.66654897
	t.M[0x55] = 6547.5541631
	fmt.Println(strike.ConverJsonFormat(t, 2))
}
