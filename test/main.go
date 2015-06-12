package main

import (
	"fmt"
	"github.com/blinkat/blinker/strike"
)

func test_strike_js() {
	path := "../phatom/test/bootstrap.js"
	//path := "../phatom/test/test.js"
	js := strike.StrikeJsPath(path)
	fmt.Println(js)
}

func test_combine() {
	path := "../phatom/setting.json"
	js, css, err := strike.CombineForJson(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(js)
	fmt.Println(css)
}

func main() {
	test_strike_js()
}
