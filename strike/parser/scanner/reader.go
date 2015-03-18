package scanner

import (
	"github.com/blinkat/blinks/strike/parser/adapter"
	"strconv"
	"strings"
)

func hex_bytes(t *tokenizer, n int) int {
	num := 0
	for ; n > 0; n-- {
		digit, err := strconv.ParseInt(string(t.next(false)), 16, 32)
		if err != nil {
			t.throw("Invalid hex-character pattern in string")
		}
		num = (num << 4) | int(digit)
	}
	return num
}

func skip_whitespace(t *tokenizer) {
	for ch := t.peek(); ch != 0 && adapter.WhitespaceChars(ch); ch = t.peek() {
		t.next(false)
	}
}

type read_while_func func(ch rune, i int) bool

func read_while(t *tokenizer, fn read_while_func) string {
	ret := make([]rune, 0)
	ch := t.peek()
	i := 0
	for ch != 0 && fn(ch, i) {
		i += 1
		ret = append(ret, t.next(false))
		ch = t.peek()
	}
	return string(ret)
}

func read_num(t *tokenizer, prefix rune) *Token {
	has_e := false
	after_e := false
	has_x := false
	has_dot := prefix == '.'
	num := read_while(t, func(ch rune, i int) bool {
		if ch == 'x' || ch == 'X' {
			if has_x {
				return false
			}
			has_x = true
			return true
		}
		if !has_x && (ch == 'E' || ch == 'e') {
			if has_e {
				return false
			}
			has_e = true
			after_e = true
			return true
		}
		if ch == '-' {
			if after_e || i == 0 && prefix == 0 {
				return true
			}
			return false
		}
		if ch == '+' {
			return after_e
		}
		after_e = false
		if ch == '.' {
			if !has_dot && !has_x && !has_e {
				has_dot = true
				return true
			}
			return false
		}
		return adapter.IsAlphanumericChar(ch)
	})

	if prefix != 0 {
		num = string(prefix) + num
	}
	return t.token(TokenNumber, num, false)
}

func read_escaped_char(t *tokenizer, instr bool) rune {
	ch := t.next(instr)
	switch ch {
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	case 'b':
		return '\b'
	case 'v':
		return '\u000b'
	case 'f':
		return '\f'
	case 'x':
		return rune(hex_bytes(t, 2))
	case 'u':
		return rune(hex_bytes(t, 4))
	case '\n':
		return '\u0020'
	case '0':
		return '0'
	default:
		return ch
	}
}

func read_string(t *tokenizer) *Token {
	quote := t.next(false)
	ret := make([]rune, 0)
	for {
		ch := t.next(false)
		if ch == '\\' {
			otcal_len := 0
			first := rune(0)
			retstr := read_while(t, func(ch rune, i int) bool {
				if ch >= '0' && ch <= '7' {
					if first == 0 {
						first = ch
						otcal_len += 1
						return otcal_len != 0
					} else if first <= '3' && otcal_len <= 2 {
						otcal_len += 1
						return otcal_len != 0
					} else if first >= '4' && otcal_len <= 1 {
						otcal_len += 1
						return otcal_len != 0
					}
				}
				return false
			})
			if otcal_len > 0 {
				code, err := strconv.ParseInt(retstr, 8, 32)
				if err != nil {
					t.throw("Unterminated string constant")
				}
				ch = rune(code)
			} else {
				ch = read_escaped_char(t, true)
			}
		} else if ch == quote {
			break
		} else if ch == '\n' {
			t.throw("Unterminated string constant")
		}
		ret = append(ret, ch)
	}
	return t.token(TokenString, string(ret), false)
}

func read_line_comment(t *tokenizer) *Token {
	t.next(false)
	i := t.find_char('\n')
	var ret string
	if i == -1 {
		ret = string(t.text[t.pos:])
		t.pos = -1
	} else {
		ret = string(t.text[t.pos:i])
		t.pos = i
	}
	return t.token(TokenLineComment, ret, true)
}

func read_multiline_comment(t *tokenizer) *Token {
	t.next(false)
	i := t.find_str("*/")
	if i == -1 {
		t.throw("Unterminated multiline comment")
	}
	text := string(t.text[t.pos:i])
	t.pos = i + 2
	nlen := len(strings.Split(text, "\n"))
	t.line += nlen - 1
	t.newline_befor = t.newline_befor || nlen > 0
	//if adapter.TestComment(text) {
	//}
	return t.token(TokenMultComment, text, true)
}

func read_name(t *tokenizer) string {
	backslash := false
	name := make([]rune, 0)
	ch := rune(0)
	escaped := false

	for ch = t.peek(); ch != 0; ch = t.peek() {
		if !backslash {
			if ch == '\\' {
				escaped = true
				backslash = true
				t.next(false)
			} else if adapter.IsIdentifierChar(ch) {
				name = append(name, ch)
			} else {
				break
			}
		} else {
			if ch != 'u' {
				t.throw("Expecting UnicodeEscapeSequence -- uXXXX")
			}
			ch = read_escaped_char(t, false)
			if !adapter.IsIdentifierChar(ch) {
				t.throw("Unicode char: " + string(ch) + " is not valid in identifier")
			}
			name = append(name, ch)
			backslash = false
		}
	}

	ret := string(name)
	if adapter.Keywords(ret) && escaped {
		hex := strconv.FormatInt(int64(name[0]), 16)
		hex = strings.ToUpper(hex)
		rtemp := []rune("0000")
		ret = "\\u" + string(rtemp[len(hex):]) + hex + string(name[1:])
	}
	return ret
}

func read_regexp(t *tokenizer, regexp string) *Token {
	prev := false
	ch := rune(0)
	in_class := false
	ret := []rune(regexp)
	for ch = t.next(false); t.eof(); ch = t.next(false) {
		if prev {
			ret = append(ret, '\\', ch)
			prev = false
		} else if ch == '[' {
			in_class = true
			ret = append(ret, ch)
		} else if ch == ']' && in_class {
			in_class = false
			ret = append(ret, ch)
		} else if ch == '/' && !in_class {
			break
		} else if ch == '\\' {
			prev = true
		} else {
			ret = append(ret, ch)
		}
	}
	mods := read_name(t)
	tok := t.token(TokenRegexp, string(ret), false)
	tok.Attributes = []string{mods}
	return tok
}

func read_operator(t *tokenizer, prefix rune) *Token {
	ret := make([]rune, 0)
	if prefix == 0 {
		ret = append(ret, t.next(false))
	} else {
		ret = append(ret, prefix)
	}
	return t.token(TokenOperator, string(read_operator_grow(t, ret)), false)

}

func read_operator_grow(t *tokenizer, op []rune) []rune {
	if t.peek() == 0 {
		return op
	}
	bigger := append(op, t.peek())
	if adapter.Operator(string(bigger)) {
		t.next(false)
		return read_operator_grow(t, bigger)
	} else {
		return op
	}
}

func handle_slash(t *tokenizer) *Token {
	t.next(false)
	regex_allowed := t.regex_allowed
	switch t.peek() {
	case '/':
		t.comments_befor = append(t.comments_befor, read_line_comment(t))
		t.regex_allowed = regex_allowed
		return t.next_token("")
	case '*':
		t.comments_befor = append(t.comments_befor, read_multiline_comment(t))
		t.regex_allowed = regex_allowed
		return t.next_token("")
	}
	if t.regex_allowed {
		return read_regexp(t, "")
	} else {
		return read_operator(t, '/')
	}
}

func handle_dot(t *tokenizer) *Token {
	t.next(false)
	if adapter.IsDigit(t.peek()) {
		return read_num(t, '.')
	} else {
		return t.token(TokenPunc, ".", false)
	}
}

func read_word(t *tokenizer) *Token {
	word := read_name(t)
	if !adapter.Keywords(word) {
		return t.token(TokenName, word, false)
	} else if adapter.Operator(word) {
		return t.token(TokenKeyword, word, false)
	} else if adapter.KeywordsAtom(word) {
		return t.token(TokenAtom, word, false)
	} else {
		return t.token(TokenKeyword, word, false)
	}
}
