package main

import (
	"fmt"
	"github.com/blinkat/blinks/strike/parser"
	"io/ioutil"
)

var jquery string

func main() {
	jquery = "./jquery-1.11.2.js"

	buffer, err := ioutil.ReadFile(jquery)
	if err != nil {
		fmt.Println(err)
	}

	text := string(buffer)
	parser.Test(text)
}
