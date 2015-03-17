package adapter

//----------------[ string array ]-------------------
var sUNARY_PREFIX []string
var sUNARY_POSTFIX []string
var sSTATEMENTS_WITH_LABELS []string
var sATOMIC_START_TOKEN []string
var sKEYWORDS []string
var sRESERVED_WORDS []string
var sKEYWORDS_BEFOR_EXPRESSION []string
var sKEYWORDS_ATOM []string
var sOPERATOR []string

//----------------[ chars ]-----------------
var cOPERATOR_CHARS []rune
var cWHITESPACE_CHARS []rune
var cPUNC_BEFORE_EXPRESSION []rune
var cPUNC_CHARS []rune
var cREGEXP_MODIFIERS []rune

//----------------[ map ]---------------
var mASSIGNMENT map[string]string
var mPRECEDENCE map[string]int

func init_array() {
	sUNARY_PREFIX = []string{
		"typeof",
		"void",
		"delete",
		"--",
		"++",
		"!",
		"~",
		"-",
		"+",
	}
	sUNARY_POSTFIX = []string{"--", "++"}

	mASSIGNMENT = make(map[string]string)
	mASSIGNMENT["+="] = "+"
	mASSIGNMENT["-="] = "-"
	mASSIGNMENT["/="] = "/"
	mASSIGNMENT["*="] = "*"
	mASSIGNMENT["%="] = "%"
	mASSIGNMENT[">>="] = ">>"
	mASSIGNMENT["<<="] = "<<"
	mASSIGNMENT[">>>="] = ">>>"
	mASSIGNMENT["|="] = "|"
	mASSIGNMENT["^="] = "^"
	mASSIGNMENT["&="] = "&"
	mASSIGNMENT["="] = "true"

	mPRECEDENCE = make(map[string]int)
	mPRECEDENCE["||"] = 1
	mPRECEDENCE["&&"] = 2
	mPRECEDENCE["|"] = 3
	mPRECEDENCE["^"] = 4
	mPRECEDENCE["&"] = 5
	mPRECEDENCE["=="] = 6
	mPRECEDENCE["==="] = 6
	mPRECEDENCE["!="] = 6
	mPRECEDENCE["!=="] = 6
	mPRECEDENCE["<"] = 7
	mPRECEDENCE[">"] = 7
	mPRECEDENCE["<="] = 7
	mPRECEDENCE[">="] = 7
	mPRECEDENCE["in"] = 7
	mPRECEDENCE["instanceof"] = 7
	mPRECEDENCE[">>"] = 8
	mPRECEDENCE["<<"] = 8
	mPRECEDENCE[">>>"] = 8
	mPRECEDENCE["+"] = 9
	mPRECEDENCE["-"] = 9
	mPRECEDENCE["*"] = 10
	mPRECEDENCE["/"] = 10
	mPRECEDENCE["%"] = 10

	sSTATEMENTS_WITH_LABELS = []string{
		"for", "do", "while", "switch",
	}

	sATOMIC_START_TOKEN = []string{"atom", "num", "string", "regexp", "name"}
	sKEYWORDS = []string{
		"break",
		"case",
		"catch",
		"const",
		"continue",
		"debugger",
		"default",
		"delete",
		"do",
		"else",
		"finally",
		"for",
		"function",
		"if",
		"in",
		"instanceof",
		"new",
		"return",
		"switch",
		"throw",
		"try",
		"typeof",
		"var",
		"void",
		"while",
		"with",
	}
	sRESERVED_WORDS = []string{
		"abstract",
		"boolean",
		"byte",
		"char",
		"class",
		"double",
		"enum",
		"export",
		"extends",
		"final",
		"float",
		"goto",
		"implements",
		"import",
		"int",
		"interface",
		"long",
		"native",
		"package",
		"private",
		"protected",
		"public",
		"short",
		"static",
		"super",
		"synchronized",
		"throws",
		"transient",
		"volatile",
	}
	sKEYWORDS_BEFOR_EXPRESSION = []string{
		"return",
		"new",
		"delete",
		"throw",
		"else",
		"case",
	}
	sKEYWORDS_ATOM = []string{
		"false",
		"null",
		"true",
		"undefined",
	}
	sOPERATOR = []string{
		"in",
		"instanceof",
		"typeof",
		"new",
		"void",
		"delete",
		"++",
		"--",
		"+",
		"-",
		"!",
		"~",
		"&",
		"|",
		"^",
		"*",
		"/",
		"%",
		">>",
		"<<",
		">>>",
		"<",
		">",
		"<=",
		">=",
		"==",
		"===",
		"!=",
		"!==",
		"?",
		"=",
		"+=",
		"-=",
		"/=",
		"*=",
		"%=",
		">>=",
		"<<=",
		">>>=",
		"|=",
		"^=",
		"&=",
		"&&",
		"||",
	}

	cOPERATOR_CHARS = []rune{'+', '-', '*', '&', '%', '=', '<', '>', '!', '?', '|', '~', '^'}
	cWHITESPACE_CHARS = []rune{' ', '\u00a0', '\n', '\r', '\t', '\f', '\u000b', '\u200b', '\u180e', '\u2000', '\u2001', '\u2002', '\u2003', '\u2004', '\u2005', '\u2006', '\u2007', '\u2008', '\u2009', '\u200a', '\u202f', '\u205f', '\u3000', '\uFEFF'}
	cPUNC_BEFORE_EXPRESSION = []rune{'[', '{', '(', ',', '.', ';', ':'}
	cPUNC_CHARS = []rune{'[', ']', '{', '}', '(', ')', ',', ';', ':'}
	cREGEXP_MODIFIERS = []rune{'g', 'm', 's', 'i', 'y'}
}
