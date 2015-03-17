package scanner

import (
	"github.com/blinkat/blinks/strike/parser/adapter"
)

type tokenizer struct {
	text           []rune
	pos            int
	tokpos         int
	line           int
	tokline        int
	col            int
	tokcol         int
	newline_befor  bool
	regex_allowed  bool
	length         int
	comments_befor []*Token
}

//----------[ ctor ]------------
func generatorTokenizer(text string) *tokenizer {
	t := &tokenizer{}
	t.text = []rune(adapter.ClearWhite(text))
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

func (t *tokenizer) peek() rune {
	if t.eof() {
		return t.text[t.pos]
	}
	return 0
}

func (t *tokenizer) next(in_string bool) rune {
	ch := t.peek()
	if ch == 0 {
		return 0
	}
	t.pos += 1
	if ch == '\n' {
		t.newline_befor = t.newline_befor || !in_string
		t.line += 1
		t.col = 0
	} else {
		t.col += 1
	}
	return ch
}

func (t *tokenizer) eof() bool {
	return t.pos < t.length && t.pos > 0
}

func (t *tokenizer) find_str(what string) int {
	runs := []rune(what)
	leng := len(runs)

	for ; t.pos > -1 && t.length < t.pos+leng; t.pos++ {
		if adapter.ComparedRuneSlice(t.text[t.pos:t.pos+leng], runs) {
			return t.pos
		}
	}
	t.pos = t.length
	return -1
}

func (t *tokenizer) find_char(ch rune) int {
	for ; t.eof(); t.pos++ {
		if t.peek() == ch {
			return t.pos
		}
	}
	return -1
}
