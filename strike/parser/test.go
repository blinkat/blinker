package parser

import (
	"fmt"
)

func Test(t string) {
	js := generator_js(t)
	p := js.Parse()
	fmt.Println(p)
}
