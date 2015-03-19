package parser

import (
	"github.com/blinkat/blinks/strike"
	"github.com/blinkat/blinks/strike/parser/scanner"
)

func Parser(text string, parser_type int) *Ast {
	switch parser_type {
	case strike.JS_PARSER:

	}
	return nil
}

//--------------[ ctors ]----------------
func generator_js(text string) *jsparser {
	j := &jsparser{}
	j.in_directives = true
	j.input = scanner.GeneratorTokenizer(strike.JS_PARSER, text)
	j.prev = nil
	j.peeked = nil
	j.in_func = 0
	j.in_loop = 0
	j.labels = make([]string, 0)
	j.token = j.next()
	return j
}
