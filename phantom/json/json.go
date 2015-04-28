package json

import "strconv"
import "bytes"

type Json []JsonValue

func (j Json) Get(name string) JsonValue {
	for _, v := range j {
		switch v.(type) {
		case *JsonBlock:
			if v.(*JsonBlock).Name == name {
				return v.(*JsonBlock).Value
			}
		}
	}
	return nil
}

func (j Json) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	leng := len(j) - 1
	for k, v := range j {
		buf.WriteString(v.String())
		if k < leng && !is_comment(v) {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune('}')
	return buf.String()
}

// json value type
type JsonValue interface {
	String() string
}

type JsonArray []JsonValue

func (j JsonArray) String() string {
	var buf bytes.Buffer
	buf.WriteRune('[')
	leng := len(j) - 1
	for k, v := range j {
		buf.WriteString(v.String())
		if k < leng {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune(']')
	return buf.String()
}

type JsonString string

func (j JsonString) String() string {
	return "\"" + string(j) + "\""
}

type JsonBool bool

func (j JsonBool) String() string {
	if j {
		return "true"
	} else {
		return "false"
	}
}

type JsonNumber float64

func (j JsonNumber) String() string {
	return strconv.FormatFloat(float64(j), 'f', -1, 64)
}

type JsonBlock struct {
	Name  string
	Value JsonValue
}

func (j *JsonBlock) String() string {
	return "\"" + j.Name + "\":" + j.Value.String()
}

type JsonMultComment string

func (j JsonMultComment) String() string {
	return "/*" + string(j) + "*/"
}

type JsonSingleComment string

func (j JsonSingleComment) String() string {
	return "//" + string(j) + "\n"
}

// ------------ [ init ] ----------------
func init() {
	json_punc_chars = []byte{
		'[',
		']',
		'{',
		'}',
		':',
		',',
	}
}

func is_punc_chars(b byte) bool {
	for _, v := range json_punc_chars {
		if v == b {
			return true
		}
	}
	return false
}

func is_comment(val JsonValue) bool {
	switch val.(type) {
	case JsonMultComment, JsonSingleComment:
		return true
	}
	return false
}

func is_array(val JsonValue) bool {
	switch val.(type) {
	case JsonArray:
		return true
	}
	return false
}

func is_json(val JsonValue) bool {
	switch val.(type) {
	case Json:
		return true
	}
	return false
}
