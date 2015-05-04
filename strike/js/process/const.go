package process

import (
	p "github.com/blinkat/blinker/strike/js/parser"
	"regexp"
)

var unary_prefix_symbol []string
var binary_symbol_1 []string
var binary_symbol_2 []string
var is_break []string
var is_number *regexp.Regexp
var ok_ops []string
var dot_call_no_parens []int
var squeeze_binary_arr []int
var squeeze_binary_arr2 []string

var make_num_e *regexp.Regexp
var make_num_match1 *regexp.Regexp
var make_num_match2 *regexp.Regexp

var make_block_code *regexp.Regexp

var squeeze_dot *regexp.Regexp
var squeeze_for *regexp.Regexp
var squeeze_for_in *regexp.Regexp

func init() {
	unary_prefix_symbol = []string{"!", "delete"}
	binary_symbol_1 = []string{"in", "instanceof", "==", "!=", "===", "!==", "<", "<=", ">=", ">"}
	binary_symbol_2 = []string{"&&", "||"}
	is_break = []string{"return", "throw", "break", "continue"}
	is_number = regexp.MustCompile("^[1-9][0-9]*$")
	ok_ops = []string{"+", "-", "/", "*", "%", ">>", "<<", ">>>", "|", "^", "&"}

	make_num_e = regexp.MustCompile("^0\\.")
	make_num_match1 = regexp.MustCompile("^(.*?)(0+)$")
	make_num_match2 = regexp.MustCompile("^0?\\.(0+)(.*)$")
	make_block_code = regexp.MustCompile(";+\\s*$")

	squeeze_dot = regexp.MustCompile("(?i)[a-f.]")
	squeeze_for = regexp.MustCompile(";*\\s*$")
	squeeze_for_in = regexp.MustCompile(";+$")
	squeeze_binary_arr2 = []string{
		"&&", "||", "*",
	}

	dot_call_no_parens = []int{
		p.Type_Name,
		p.Type_Array,
		p.Type_Object,
		p.Type_String,
		p.Type_Dot,
		p.Type_Sub,
		p.Type_Call,
		p.Type_Regexp,
		p.Type_Defunc,
	}

	squeeze_binary_arr = []int{
		p.Type_Assign,
		p.Type_Conditional,
		p.Type_Seq,
	}
}

func member(arr []string, s string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}

func member_int(arr []int, s int) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}
