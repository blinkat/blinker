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

func test_combine() *strike.Combined {
	path := "../phatom/setting.json"
	ret, err := strike.CombineForJson(path)
	if err != nil {
		fmt.Println(err)
	}
	return ret
}

func main() {
	ret := test_combine()
	s := strike.StrikeJsText(ret.Javascript)
	fmt.Println(s)
}
