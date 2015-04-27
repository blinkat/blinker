package parser

import "github.com/blinkat/blinks/strike/parser/adapter"
import "regexp"

var name_read_end []byte
var tag_comment *regexp.Regexp
var tag_white *regexp.Regexp
var can_single_tag []string

func init() {
	name_read_end = []byte{'/', '>', '='}
	tag_comment = regexp.MustCompile("<!--(\\s|.)*?--[ ]*>")
	tag_white = regexp.MustCompile("[ \u00a0\n\r\t\f\u000b\u200b\u180e\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u202f\u205f\u3000\uFEFF]")
	can_single_tag = []string{
		"input",
		"img",
	}
}

func CanSingle(s string) bool {
	for _, v := range can_single_tag {
		if v == s {
			return true
		}
	}
	return false
}

func is_end_name(ch byte) bool {
	if adapter.WhitespaceChars(rune(ch)) {
		return true
	}
	for _, v := range name_read_end {
		if v == ch {
			return true
		}
	}
	return false
}
