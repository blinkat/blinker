package strike

import (
	"regexp"
)

//------------[ css ]-------------
var protect_string *regexp.Regexp
var calc_ *regexp.Regexp
var space_remover *regexp.Regexp
var background_position *regexp.Regexp
var color_ *regexp.Regexp
var color_shortener *regexp.Regexp
var border_ *regexp.Regexp
var empty_rule *regexp.Regexp
var comment_like *regexp.Regexp
var comment_ *regexp.Regexp

var preserve_candidate_comment []byte
var preserved_token []byte
var pseudo_class_colon []byte

func init_css_regexp() {
	protect_string = regexp.MustCompile("(\"([^\\\\\"]|\\\\.|\\\\)*\")|('([^\\\\']|\\\\.|\\\\)*')")
	calc_ = regexp.MustCompile("(calc\\s*\\([^};]*\\))")
	space_remover = regexp.MustCompile("(^|\\})(([^\\{:])+:)+([^\\{]*\\{)")
	background_position = regexp.MustCompile("(?i)(background-position|transform-origin|webkit-transform-origin|moz-transform-origin|o-transform-origin|ms-transform-origin):0(;|})")
	color_ = regexp.MustCompile("rgb\\s*\\(\\s*([0-9,\\s]+)\\s*\\)")
	color_shortener = regexp.MustCompile("([^\"'=\\s])(\\s*)#([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])")
	border_ = regexp.MustCompile("(?i)(border|border-top|border-right|border-bottom|border-right|outline|background):none(;|})")
	empty_rule = regexp.MustCompile("[^\\}\\{/;]+\\{\\}")
	comment_like = regexp.MustCompile("(?i)progid:DXImageTransform.Microsoft.Alpha\\(Opacity=")
	comment_ = regexp.MustCompile("\\/\\*(\\s|.)*?\\*\\/")

	preserved_token = []byte("___YUICSSMIN_PRESERVED_TOKEN_")
	preserve_candidate_comment = []byte("___YUICSSMIN_PRESERVE_CANDIDATE_COMMENT_")
	pseudo_class_colon = []byte("___YUICSSMIN_PSEUDOCLASSCOLON___")
}

//------------[ init ]------------
func init() {
	init_css_regexp()
}
