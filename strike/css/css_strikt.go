package css

import (
	"bytes"
	"regexp"
	"strconv"
)

func strike_css(codes []byte) []byte {
	//codes := []byte(text)
	null := []byte{}
	//line := []byte{'-', '-', '-'}
	preserved_tokens := make([][]byte, 0)

	// ---- remove comment ----
	codes = comment_.ReplaceAll(codes, null)
	codes = protect_string.ReplaceAllFunc(codes, func(str []byte) []byte {
		quote := str[0]
		token := str[1 : len(str)-1]
		ret := comment_like.ReplaceAll(token, []byte("alpha(opacity="))
		preserved_tokens = append(preserved_tokens, ret)
		preserver := append([]byte{quote}, gen_flag(preserved_token, len(preserved_tokens)-1, 10)...)
		preserver = append(preserver, quote)
		return preserver
	})

	// preserve any calc(...)
	codes = calc_.ReplaceAllFunc(codes, func(val []byte) []byte {
		preserved_tokens = append(preserved_tokens, val)
		//preserver := preserved_token + fmt.Sprint(len(preserved_tokens)-1) + "___"
		return gen_flag(preserved_token, len(preserved_tokens)-1, 10)
	})

	codes = regexp_replace("\\s+", codes, []byte{' '})

	codes = space_remover.ReplaceAllFunc(codes, func(val []byte) []byte {
		val = bytes.Replace(val, []byte{':'}, pseudo_class_colon, -1)
		val = bytes.Replace(val, []byte{'\\', '\\'}, []byte{'\\', '\\', '\\', '\\'}, -1)
		val = bytes.Replace(val, []byte{'\\', '$'}, []byte{'\\', '\\', '\\', '$'}, -1)
		return val
	})

	codes = regexp_replace("\\s+([!{};:>+\\(\\)\\],])", codes, []byte("$1"))
	codes = regexp_replace(string(pseudo_class_colon), codes, []byte{':'})
	codes = regexp_replace(":first\\-(line|letter)(\\{|,)", codes, []byte(":first-$1 $2"))
	codes = regexp_replace("\\*/", codes, []byte("*/"))
	codes = regexp_replace("^(.*)(@charset \"[^\"]*\";)", codes, []byte("$2$1"))
	codes = regexp_replace("^(\\s*@charset [^;]+;\\s*)+", codes, []byte("$1"))
	codes = regexp_replace("\\band\\(", codes, []byte("and ("))
	codes = regexp_replace("([!{}:;>+\\(\\[,])\\s+", codes, []byte("$1"))
	codes = regexp_replace(";+}", codes, []byte("}"))
	codes = regexp_replace("([\\s:])(0)(px|em|%|in|cm|mm|pc|pt|ex)", codes, []byte("$1$2"))
	codes = regexp_replace(":0 0 0 0(;|})", codes, []byte(":0$1"))
	codes = regexp_replace(":0 0 0(;|})", codes, []byte(":0$1"))
	codes = regexp_replace(":0 0(;|})", codes, []byte(":0$1"))

	codes = background_position.ReplaceAllFunc(codes, func(val []byte) []byte {
		arr := bytes.SplitN(bytes.ToLower(val), []byte{':'}, 1)
		if len(arr) >= 2 {
			var buffer bytes.Buffer
			buffer.Write(arr[0])
			buffer.WriteString(":0 0")
			buffer.Write(arr[1])
			return buffer.Bytes()
		}
		return val
	})

	codes = regexp_replace("(:|\\s)0+\\.(\\d+)", codes, []byte("$1.$2"))

	// rgb conver to hex
	codes = color_.ReplaceAllFunc(codes, func(val []byte) []byte {

		start := bytes.IndexByte(val, '(')
		end := bytes.IndexByte(val, ')')
		if start != -1 && end != -1 {
			rgb := bytes.Split(val, []byte{','})
			var buffer bytes.Buffer
			buffer.WriteRune('#')
			for _, v := range rgb {
				num, err := strconv.ParseInt(string(v), 10, 16)
				if err != nil {
					return val
				}
				if num < 16 {
					buffer.WriteRune('0')
				}
				buffer.WriteString(strconv.FormatInt(num, 16))
			}
			return buffer.Bytes()
		}
		return val
	})

	// #ffffff to #fff
	codes = color_shortener.ReplaceAllFunc(codes, func(val []byte) []byte {
		strs := bytes.Split(val, []byte{'#'})
		if len(strs) > 1 {
			hex := strs[1]

			if len(hex) == 6 && hex[0] == hex[1] && hex[2] == hex[3] && hex[4] == hex[5] {
				var buffer bytes.Buffer
				buffer.Write(strs[0])
				buffer.WriteByte('#')
				buffer.WriteByte(hex[0])
				buffer.WriteByte(hex[2])
				buffer.WriteByte(hex[4])
				return buffer.Bytes()
			}
			return val
		}
		return val
	})

	// border:none -> border:0
	codes = border_.ReplaceAllFunc(codes, func(val []byte) []byte {
		strs := bytes.SplitN(val, []byte{':'}, 1)
		if len(strs) > 1 {
			return append(strs[0], ':', '0')
		}
		return val
	})

	codes = comment_like.ReplaceAll(codes, []byte("alpha(opacity="))

	if bytes.Contains(codes, []byte{'{', '}'}) {
		codes = empty_rule.ReplaceAll(codes, null)
	}

	for i, v := range preserved_tokens {
		//codes = bytes.Replace(codes, preserved_token+fmt.Sprint(i)+"___", v, 1)
		codes = bytes.Replace(codes, gen_flag(preserved_token, i, 10), v, 1)
	}

	codes = bytes.Trim(codes, "")

	return codes
}

func runes_indexof_rune(what rune, start int, runes []rune) int {
	for i := start; i < len(runes); i++ {
		if runes[i] == what {
			return i
		}
	}
	return -1
}

func runes_indexof(what string, start int, runes []rune) int {
	what_rune := []rune(what)
	what_leng := len(what_rune)
	max_leng := len(runes)

	for i := start; i < max_leng-what_leng; i++ {
		j := 0
		for j < what_leng {
			if runes[i+j] != what_rune[j] {
				break
			}
			j += 1
		}
		if j == what_leng {
			return i
		}
	}
	return -1
}

func regexp_replace(regex string, source, news []byte) []byte {
	re := regexp.MustCompile(regex)
	return re.ReplaceAll(source, news)
}

func gen_flag(arr []byte, i int, base int) []byte {
	ret := strconv.AppendInt(arr, int64(i), base)
	ret = append(ret, '-', '-', '-')
	return ret
}
