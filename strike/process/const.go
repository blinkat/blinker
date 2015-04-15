package process

import (
	"regexp"
)

var unary_prefix_symbol []string
var binary_symbol_1 []string
var binary_symbol_2 []string
var is_break []string
var is_number *regexp.Regexp
var ok_ops []string

func init() {
	unary_prefix_symbol = []string{"!", "delete"}
	binary_symbol_1 = []string{"in", "instanceof", "==", "!=", "===", "!==", "<", "<=", ">=", ">"}
	binary_symbol_2 = []string{"&&", "||"}
	is_break = []string{"return", "throw", "break", "continue"}
	is_number = regexp.MustCompile("^[1-9][0-9]*$")
	ok_ops = []string{"+", "-", "/", "*", "%", ">>", "<<", ">>>", "|", "^", "&"}
}

func member(arr []string, s string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}
