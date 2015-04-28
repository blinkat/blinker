package html

import (
	"bytes"
	"github.com/blinkat/blinks/phantom/parser"
	"github.com/blinkat/blinks/storm"
	"github.com/blinkat/blinks/strike"
	"io/ioutil"
)

func HandleHtml(path string) string {
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		storm.Error("\"" + path + "\" can not find file!")
	}

	res := parser.Parse(buffer)
	tag := get_html(res)
	str := gen_code(tag)
	return str
}

// generator codes
func gen_code(tag *parser.Tag) string {
	//fmt.Println(tag.Name)
	if tag.IsString {
		return tag.Name
	}
	if parser.IsBtml(tag.Name) {
		return gen_btml(tag)
	} else {
		return gen_html(tag)
	}
}

func gen_html(tag *parser.Tag) string {
	var buffer bytes.Buffer
	if tag.IsDoctype {
		return "<" + tag.Name + ">"
	} else {
		//write head
		buffer.WriteString("<" + tag.Name)
		for k, v := range tag.Attribute {
			if k == "style" {
				buffer.WriteString(" style=\"" + strike.StrikeCss(v) + "\"")
			} else {
				buffer.WriteString(" " + k + "=\"" + v + "\"")
			}
		}
	}

	if len(tag.Children) == 0 {
		if parser.CanSingle(tag.Name) {
			buffer.WriteString("/>")
		} else {
			buffer.WriteString("></" + tag.Name + ">")
		}
	} else {
		buffer.WriteRune('>')
		children := gen_children(tag)
		if tag.Name == "script" {
			children = strike.StrikeJs(children)
		}
		buffer.WriteString(children)
		buffer.WriteString("</" + tag.Name + ">")
	}
	return buffer.String()
}

func gen_btml(tag *parser.Tag) string {
	if parser.IsNonHeadBtml(tag.Name) {
		return gen_children(tag)
	}
	return ""
}

func gen_children(tag *parser.Tag) string {
	var buf bytes.Buffer
	for _, v := range tag.Children {
		buf.WriteString(gen_code(v))
	}
	return buf.String()
}

//--------- handler ------------

func get_html(html *parser.BlinkHtml) *parser.Tag {
	if html.Master != nil {
		if path, ok := html.Master.Attribute["path"]; ok {
			buffer, err := ioutil.ReadFile(path)
			if err != nil {
				storm.Error("\"" + path + "\" can not find master!")
			}

			master := parser.Parse(buffer)
			master = insert_content(html, master)
			return get_html(master)
		} else {
			storm.Warring("blink:master not have \"path\"")
		}
	}
	return html.TopLevel
}

func insert_content(html, master *parser.BlinkHtml) *parser.BlinkHtml {
	for k, content := range master.Contents {
		if page, ok := html.PageArea[k]; ok {
			content.Children = page.Children
		}
	}

	return master
}
