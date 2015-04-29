package json

import (
	"bytes"
	"fmt"
	"github.com/blinkat/blinks/strike/js/parser/adapter"
	"io/ioutil"
	"strconv"
)

// read json file

func ReadFile(path string) (Json, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ReadBytes(file)
}

func ReadText(text string) (Json, error) {
	return ReadBytes([]byte(text))
}

func ReadBytes(buffer []byte) (Json, error) {
	scanner := &json_scanner{
		source: buffer,
		length: len(buffer),
		pos:    0,
	}

	val, ty := scanner.next()
	if val[0] == '{' && ty == json_punc {
		return read_body(scanner)
	} else {
		return nil, fmt.Errorf("symbol error")
	}
}

// remove comment
func RemoveComment(json JsonValue) JsonValue {
	switch json.(type) {
	case JsonArray:
		ret := make(JsonArray, 0)
		for _, v := range json.(JsonArray) {
			if !is_comment(v) {
				if is_array(v) || is_json(v) {
					RemoveComment(v)
				}
				ret = append(ret, v)
			}
		}
		return ret
	case Json:
		{
			ret := make(Json, 0)
			for _, v := range json.(Json) {
				if !is_comment(v) {
					if is_array(v) || is_json(v) {
						RemoveComment(v)
					}
					ret = append(ret, v)
				}
			}
			return ret
		}
	}
	return json
}

func read_body(scanner *json_scanner) (Json, error) {
	ret := make(Json, 0)
	for val, ty := scanner.next(); ty != json_eof; val, ty = scanner.next() {
		if ty == json_punc && val == "}" {
			break
		} else if ty == json_punc && val == "," {
			continue
		} else if ty == json_mult_comment || ty == json_single_comment {
			value, err := gen_json_value(val, ty)
			if err != nil {
				return nil, err
			}
			ret = append(ret, value)
			continue
		} else if ty == json_name {
			name := val
			val, ty := scanner.next()
			if val[0] != ':' {
				return nil, fmt.Errorf("lost ':'")
			}
			val, ty = scanner.next()
			var value JsonValue
			var err error
			if ty == json_punc {
				switch val[0] {
				case '{':
					value, err = read_body(scanner)
				case '[':
					value, err = read_array(scanner)
				default:
					return nil, fmt.Errorf("symbol error:" + val)
				}
			} else {
				value, err = gen_json_value(val, ty)
			}
			if err != nil {
				return nil, err
			}
			ret = append(ret, &JsonBlock{
				Name:  name,
				Value: value,
			})
		} else {
			return nil, fmt.Errorf("symbol error: json loss name")
		}
	}
	return ret, nil
}

func read_array(scanner *json_scanner) (JsonValue, error) {
	ret := make(JsonArray, 0)
	for val, ty := scanner.next(); ty != json_eof; val, ty = scanner.next() {
		if ty == json_punc && val == "]" {
			break
		} else if ty == json_punc && val == "," {
			continue
		} else if ty == json_error {
			return nil, fmt.Errorf(val)
		} else {
			if ty == json_punc {
				switch val[0] {
				case '{':
					jv, err := read_body(scanner)
					if err != nil {
						return nil, err
					}
					ret = append(ret, jv)
				case ',':
				}
			} else {
				jv, err := gen_json_value(val, ty)
				if err != nil {
					return nil, err
				}
				ret = append(ret, jv)
			}
		}
	}
	return ret, nil
}

func gen_json_value(val string, t int) (JsonValue, error) {
	switch t {
	case json_value, json_name:
		return JsonString(val), nil
	case json_number:
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		return JsonNumber(num), nil
	case json_bool:
		b := val == "true"
		return JsonBool(b), nil
	case json_single_comment:
		return JsonSingleComment(val), nil
	case json_mult_comment:
		return JsonMultComment(val), nil
	}
	return nil, fmt.Errorf("can not gen json value, the type code:" + fmt.Sprint(t))
}

// -------- [ scanner ] ---------
const (
	json_name = iota
	json_single_comment
	json_mult_comment
	json_punc
	json_value
	json_bool
	json_number
	json_eof
	json_error
)

var json_punc_chars []byte

type json_scanner struct {
	source []byte
	pos    int
	length int
	in_val bool
}

func (j *json_scanner) next_byte() byte {
	if j.pos >= j.length {
		return 0
	}
	ch := j.source[j.pos]
	j.pos += 1
	return ch
}

func (j *json_scanner) next() (string, int) {
	if j.pos >= j.length {
		return "", json_eof
	}
	j.jump_white()

	if j.in_val {
		j.in_val = false
		return j.read_value()
	}
	ch := j.next_byte()

	if is_punc_chars(ch) {
		if ch == ':' {
			j.in_val = true
		} else {
			j.in_val = false
		}
		return string(ch), json_punc
	}
	if ch == '/' {
		return j.read_comment()
	}
	if ch == '"' {
		return j.read_name(), json_name
	}
	return "symbol error", json_error
}

func (j *json_scanner) read_value() (string, int) {
	ch := j.next_byte()
	if ch == '"' {
		return j.read_name(), json_value
	} else if is_punc_chars(ch) {
		return string(ch), json_punc
	}

	// else is num or bool
	var buf bytes.Buffer
	for ch != 0 && !adapter.WhitespaceChars(rune(ch)) && !is_punc_chars(ch) {
		buf.WriteByte(ch)
		ch = j.next_byte()
	}

	if ch == ']' || ch == '}' {
		j.pos -= 1
	}

	val := buf.String()
	if val == "true" || val == "false" {
		return val, json_bool
	} else if adapter.IsFloat(val) || adapter.IsDecNumber(val) {
		return val, json_number
	} else {
		return "\"" + val + "\" is not string , number or bool", json_error
	}
}

func (j *json_scanner) read_name() string {
	var buf bytes.Buffer
	prev := byte(0)
	for ch := j.next_byte(); ch != 0; ch = j.next_byte() {
		if ch == '"' {
			if prev != '\\' {
				break
			}
		}
		buf.WriteByte(ch)
		prev = ch
	}
	return buf.String()
}

func (j *json_scanner) read_comment() (string, int) {
	ch := j.next_byte()
	var buffer bytes.Buffer
	t := 0
	if ch == '*' {
		t = json_mult_comment
		for {
			ch = j.next_byte()
			if ch == 0 {
				break
			}
			if ch == '*' {
				ch = j.next_byte()
				if ch == '/' {
					break
				} else {
					buffer.WriteByte('*')
				}
			}
			buffer.WriteByte(ch)
		}
	} else {
		t = json_single_comment
		for {
			ch = j.next_byte()
			if ch == '\n' || ch == 0 {
				break
			}
			buffer.WriteByte(ch)
		}
	}
	return buffer.String(), t
}

func (j *json_scanner) jump_white() {
	ch := j.source[j.pos]
	for ch != 0 && adapter.WhitespaceChars(rune(ch)) {
		j.next_byte()
		ch = j.source[j.pos]
	}
}
