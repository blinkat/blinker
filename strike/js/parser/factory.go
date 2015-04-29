package parser

import (
	"github.com/blinkat/blinks/strike/js/parser/scanner"
)

func ParseJs(text string) IAst {
	p := generator_js(text)
	return p.Parse()
}

//--------------[ ctors ]----------------
func generator_js(text string) *jsparser {
	j := &jsparser{}
	j.in_directives = true
	j.input = scanner.GeneratorTokenizerJs(text)
	j.prev = nil
	j.peeked = nil
	j.in_func = 0
	j.in_loop = 0
	j.labels = make([]string, 0)
	j.token = j.next()
	return j
}

func init() {
	init_digits()
}
