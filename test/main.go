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
	test = "./test.css"

	buffer, err := ioutil.ReadFile(jquery)
	if err != nil {
		fmt.Println(err)
	}

	//text := string(buffer)
	strike.Test(buffer)
}
