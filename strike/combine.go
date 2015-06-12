package strike

import (
	"bytes"
	//"github.com/blinkat/blinker/strike/css"
	//"github.com/blinkat/blinker/strike/js"
	j "github.com/blinkat/blinker/strike/json"
	"io/ioutil"
)

type Combined struct {
	Javascript string
	Css        string
	Name       string
}

// return js css
func CombineForJson(path string) (*Combined, error) {
	json, err := j.ReadFile(path)
	if err != nil {
		return nil, err
	}

	prefix := get_base_path(path)
	var js_ret bytes.Buffer
	var css_ret bytes.Buffer

	if json.Get("comment") != nil {
		//comment = prefix + string(json.Get("comment").(j.JsonString))
		cmt := json.Get("comment")
		if cmt.(j.Json).Get("js") != nil {
			js := cmt.(j.Json).Get("js")
			t, e := ioutil.ReadFile(prefix + string(js.(j.JsonString)))
			if e == nil {
				js_ret.Write(t)
			}
		}

		if cmt.(j.Json).Get("css") != nil {
			css := cmt.(j.Json).Get("css")
			t, e := ioutil.ReadFile(prefix + string(css.(j.JsonString)))
			if e == nil {
				css_ret.Write(t)
			}
		}
	}

	js := json.Get("js")
	if js != nil {
		write(&js_ret, js.(j.JsonArray), prefix+"js/", ".js")
	}

	css := json.Get("css")
	if css != nil {
		write(&css_ret, css.(j.JsonArray), prefix+"css/", ".css")
	}

	return &Combined{
		Javascript: string(js_ret.Bytes()),
		Css:        string(css_ret.Bytes()),
		Name:       string(json.Get("name").(j.JsonString)),
	}, nil
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
