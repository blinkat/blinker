package main

import (
	"fmt"
	"github.com/blinkat/blinks/strike"
	"io/ioutil"
)

var jquery string
var test string

func main() {
	jquery = "./jquery-1.11.2.js"
	test = "./test.js"

	buffer, err := ioutil.ReadFile(test)
	if err != nil {
		fmt.Println(err)
	}

	text := string(buffer)
	strike.Test(text)
}
