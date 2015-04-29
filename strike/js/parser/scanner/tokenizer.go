package scanner

import (
	"fmt"
	"github.com/blinkat/blinks/strike/js/parser/adapter"
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

//------------[ interface ]--------------
func (t *tokenizer) Next(regexp string) *Token {
	return t.next_token(regexp)
}
func (t *tokenizer) Eof() bool {
	return t.eof()
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
	return t.pos < t.length && t.pos >= 0
}

func (t *tokenizer) find_str(what string) int {
	runs := []rune(what)
	leng := len(runs)

	for i := t.pos; i > -1 && i < t.length-leng; i++ {
		j := 0
		for j < leng {
			if t.text[i+j] != runs[j] {
				break
			}
			j += 1
		}
		if j == leng {
			return i
		}
	}
	t.pos = t.length
	return -1
}

func (t *tokenizer) find_char(ch rune) int {
	for i := t.pos; t.eof(); i++ {
		if t.text[i] == ch {
			return i
		}
	}
	return -1
}

func (t *tokenizer) start_token() {
	t.tokline = t.line
	t.tokcol = t.col
	t.tokpos = t.pos
}

func (t *tokenizer) token(ty int, value string, is_comment bool) *Token {
	t.regexAllowed(ty, value)
	ret := &Token{}
	ret.Type = ty
	ret.Value = value
	ret.Line = t.tokline
	ret.Col = t.tokcol
	ret.Pos = t.tokpos
	ret.Endpos = t.pos
	ret.Nlb = t.newline_befor

	if !is_comment {
		ret.CommentsBefor = t.comments_befor
		t.comments_befor = make([]*Token, 0)

		for _, v := range ret.CommentsBefor {
			ret.Nlb = ret.Nlb || v.Nlb
		}
	}
	t.newline_befor = false
	return ret
}

func (t *tokenizer) regexAllowed(ty int, value string) {
	b := ty == TokenOperator && !adapter.UnaryPostfix(value)
	b = b || (ty == TokenKeyword && adapter.KeywordsBeforExpression(value))

	rs := []rune(value)
	if len(rs) == 1 {
		b = b || (ty == TokenPunc && adapter.PuncBeforExpression(rs[0]))
	}
	t.regex_allowed = b
}

func (t *tokenizer) throw(msg string) {
	panic(fmt.Sprint("sanner error:", msg, "\n\tline:", t.tokline, "colum:", t.tokcol))
}

func (t *tokenizer) next_token(regexp string) *Token {
	if regexp != "" {
		return read_regexp(t, regexp)
	}
	skip_whitespace(t)
	t.start_token()
	ch := t.peek()
	if ch == 0 {
		return nil
	}
	if adapter.IsDigit(ch) {
		return read_num(t, 0)
	}
	if ch == '"' || ch == '\'' {
		return read_string(t)
	}
	if adapter.PuncChars(ch) {
		return t.token(TokenPunc, string(t.next(false)), false)
	}
	if ch == '.' {
		return handle_dot(t)
	}
	if ch == '/' {
		return handle_slash(t)
	}
	if adapter.OperatorChars(ch) {
		return read_operator(t, 0)
	}
	if ch == '\\' || adapter.IsIdentifierStart(ch) {
		return read_word(t)
	}

	t.throw("Unexpected character '" + string(ch) + "'")
	return nil
}
