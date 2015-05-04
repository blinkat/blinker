package html

import (
	"bytes"
	"github.com/blinkat/blinker/storm"
	"github.com/blinkat/blinker/strike/css"
	"github.com/blinkat/blinker/strike/js"
	"io/ioutil"
)

func Strike(text string) string {
	res := Parse([]byte(text))
	tag := get_html(res)
	str := gen_code(tag)
	return str
}

// generator codes
func gen_code(tag *Tag) string {
	//fmt.Println(tag.Name)
	if tag.IsString {
		return tag.Name
	}
	if IsBtml(tag.Name) {
		return gen_btml(tag)
	} else {
		return gen_html(tag)
	}
}

func gen_html(tag *Tag) string {
	var buffer bytes.Buffer
	if tag.IsDoctype {
		return "<" + tag.Name + ">"
	} else {
		//write head
		buffer.WriteString("<" + tag.Name)
		for k, v := range tag.Attribute {
			if k == "style" {
				buffer.WriteString(" style=\"" + css.Strike(v) + "\"")
			} else {
				buffer.WriteString(" " + k + "=\"" + v + "\"")
			}
		}
	}

	if len(tag.Children) == 0 {
		if CanSingle(tag.Name) {
			buffer.WriteString("/>")
		} else {
			buffer.WriteString("></" + tag.Name + ">")
		}
	} else {
		buffer.WriteRune('>')
		children := gen_children(tag)
		if tag.Name == "script" {
			children = js.Strike(children)
		}
		buffer.WriteString(children)
		buffer.WriteString("</" + tag.Name + ">")
	}
	return buffer.String()
}

func gen_btml(tag *Tag) string {
	if IsNonHeadBtml(tag.Name) {
		return gen_children(tag)
	}
	return ""
}

func gen_children(tag *Tag) string {
	var buf bytes.Buffer
	for _, v := range tag.Children {
		buf.WriteString(gen_code(v))
	}
	return buf.String()
}

//--------- handler ------------

func get_html(html *BlinkHtml) *Tag {
	if html.Master != nil {
		if path, ok := html.Master.Attribute["path"]; ok {
			buffer, err := ioutil.ReadFile(path)
			if err != nil {
				storm.Error("\"" + path + "\" can not find master!")
			}

			master := Parse(buffer)
			master = insert_content(html, master)
			return get_html(master)
		} else {
			storm.Warring("blink:master not have \"path\"")
		}
	}
	return html.TopLevel
}

func insert_content(html, master *BlinkHtml) *BlinkHtml {
	for k, content := range master.Contents {
		if page, ok := html.PageArea[k]; ok {
			content.Children = page.Children
		}
	}

	return master
}
