/********************************************************************
				---------adapter----------
	adapt string for key word
	@2015-3-17

********************************************************************/
package adapter

func init() {
	init_array()
	init_regexp()
}

//------------------[ private method ]---------------------
func in_strings(arr []string, what string) bool {
	for _, v := range arr {
		if v == what {
			return true
		}
	}
	return false
}

func in_chars(arr []rune, what rune) bool {
	for _, v := range arr {
		if v == what {
			return true
		}
	}
	return false
}

//-----------------[ public method ]--------------------

func UnaryPrefix(what string) bool {
	return in_strings(sUNARY_PREFIX, what)
}

func UnaryPostfix(what string) bool {
	return in_strings(sUNARY_POSTFIX, what)
}

func StatementsWithLabels(what string) bool {
	return in_strings(sSTATEMENTS_WITH_LABELS, what)
}

func AtomicStartToken(what string) bool {
	return in_strings(sATOMIC_START_TOKEN, what)
}

func Keywords(what string) bool {
	return in_strings(sKEYWORDS, what)
}

func ReservedWords(what string) bool {
	return in_strings(sRESERVED_WORDS, what)
}

func KeywordsBeforExpression(what string) bool {
	return in_strings(sKEYWORDS_BEFOR_EXPRESSION, what)
}

func KeywordsAtom(what string) bool {
	return in_strings(sKEYWORDS_ATOM, what)
}

func Operator(what string) bool {
	return in_strings(sOPERATOR, what)
}

func OperatorChars(ch rune) bool {
	return in_chars(cOPERATOR_CHARS, ch)
}

func WhitespaceChars(ch rune) bool {
	return in_chars(cWHITESPACE_CHARS, ch)
}

func PuncBeforExpression(ch rune) bool {
	return in_chars(cPUNC_BEFORE_EXPRESSION, ch)
}

func PuncChars(ch rune) bool {
	return in_chars(cPUNC_CHARS, ch)
}
func RegexpModifiers(ch rune) bool {
	return in_chars(cREGEXP_MODIFIERS, ch)
}

func Assignment(what string) string {
	if v, ok := mASSIGNMENT[what]; ok {
		return v
	} else {
		return ""
	}
}

func Precedence(what string) int {
	if v, ok := mPRECEDENCE[what]; ok {
		return v
	}
	return -1
}

func ComparedRuneSlice(s1, s2 []rune) bool {
	len1 := len(s1)
	len2 := len(s2)
	if len1 != len2 {
		return false
	} else {
		for i := 0; i < len1; i++ {
			if s1[i] != s2[i] {
				return false
			}
		}
	}
	return true
}
