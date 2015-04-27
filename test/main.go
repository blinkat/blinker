package main

import (
	"github.com/blinkat/blinks/phantom/html"
)

var jquery string
var test string

func main() {
	jquery = "./jquery-1.11.2.js"
	test = "./page.html"
	html.HandleHtml(test)
}
