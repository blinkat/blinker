package scanner

import (
	"github.com/blinkat/blinks/strike"
	"github.com/blinkat/blinks/strike/parser/adapter"
)

type Tokenizer interface {
	Next(regexp string) *Token
	Eof() bool
}

func GeneratorTokenizer(parser_type int, text string) Tokenizer {
	switch parser_type {
	case strike.JS_PARSER:
		return newJsTokenizer(text)
	}
	return nil
}

func newJsTokenizer(text string) Tokenizer {
	t := &tokenizer{}
	t.text = []rune(adapter.ClearWhite(text))
	//t.text = []rune(text)
	t.pos = 0
	t.col = 0
	t.line = 0
	t.tokcol = 0
	t.tokline = 0
	t.tokpos = 0
	t.newline_befor = false
	t.regex_allowed = false
	t.comments_befor = make([]*Token, 0)
	t.length = len(t.text)
	return t
}
