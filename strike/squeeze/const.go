package squeeze

import "regexp"

var unary_prefix_symbol []string
var binary_symbol_1 []string
var binary_symbol_2 []string

func init() {
	unary_prefix_symbol = []string{"!", "delete"}
	binary_symbol_1 = []string{"in", "instanceof", "==", "!=", "===", "!==", "<", "<=", ">=", ">"}
	binary_symbol_2 = []string{"&&", "||"}
}

func member(arr []string, s string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}
