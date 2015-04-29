package main

import (
	"fmt"
	"github.com/blinkat/blinks/strike"
)

func main() {
	ret := strike.StrikeHtmlPath("./page.html")
	fmt.Println(ret)
}
