package strike

import (
	"bytes"
	//"github.com/blinkat/blinker/strike/css"
	//"github.com/blinkat/blinker/strike/js"
	j "github.com/blinkat/blinker/strike/json"
	"io/ioutil"
)

// return js css
func CombineForJson(path string) (string, string, error) {
	json, err := j.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	prefix := get_base_path(path)

	var comment string
	if json.Get("comment") != nil {
		comment = prefix + string(json.Get("comment").(j.JsonString))
		//comment = prefix + comment
	}

	var js_ret bytes.Buffer
	// add comment
	com, err := ioutil.ReadFile(comment)
	if err == nil {
		js_ret.WriteString("/*")
		js_ret.Write(com)
		js_ret.WriteString("*/\n")
	}

	js := json.Get("js")
	if js != nil {
		write(&js_ret, js.(j.JsonArray), prefix+"js/", ".js")
	}

	var css_ret bytes.Buffer
	css := json.Get("css")
	if css != nil {
		write(&css_ret, css.(j.JsonArray), prefix+"css/", ".css")
	}

	return string(js_ret.Bytes()), string(css_ret.Bytes()), nil
}

func get_base_path(path string) string {
	bs := []byte(path)
	for i := len(bs) - 1; i >= 0; i-- {
		if bs[i] == byte('/') || bs[i] == byte('\\') {
			return string(bs[:i+1])
		}
	}
	return path
}

func write(writer *bytes.Buffer, json j.JsonArray, prefix, postfix string) {
	for _, v := range json {
		p := prefix + string(v.(j.JsonString)) + postfix
		bs, err := ioutil.ReadFile(p)
		if err == nil {
			writer.Write(bs)
			writer.WriteRune('\n')
		}
	}
}

// combine some js file
func Combine(paths []string) string {
	var ret bytes.Buffer
	for _, v := range paths {
		text, err := ioutil.ReadFile(v)
		if err == nil {
			ret.Write(text)
			ret.WriteRune('\n')
		}
	}
	return string(ret.Bytes())
}
