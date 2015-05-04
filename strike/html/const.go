package html

import "github.com/blinkat/blinker/strike/js/parser/adapter"
import "regexp"

var name_read_end []byte
var tag_comment *regexp.Regexp
var tag_white *regexp.Regexp
var can_single_tag []string
var non_head_btml []string

func init() {
	name_read_end = []byte{'/', '>', '='}
	tag_comment = regexp.MustCompile("<!--(\\s|.)*?--[ ]*>")
	tag_white = regexp.MustCompile("[ \u00a0\n\r\t\f\u000b\u200b\u180e\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u202f\u205f\u3000\uFEFF]")
	can_single_tag = []string{
		"input",
		"img",
	}

	non_head_btml = []string{
		"blink:content",
		"blink:page",
		"blink:master",
		"blink:toplevel",
	}
}

func CanSingle(s string) bool {
	return member(can_single_tag, s)
}

// is blink html
func IsBtml(s string) bool {
	return member(non_head_btml, s)
}

func IsNonHeadBtml(s string) bool {
	return member(non_head_btml, s)
}

func member(arr []string, s string) bool {
	for _, v := range arr {
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
