package html

import (
	"bytes"
	"fmt"
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
	fmt.Println(str)
	return str
}

func gen_code(tag *parser.Tag) string {
	//fmt.Println(tag.Name)
	if tag.IsString {
		return tag.Name
	}

	var buffer bytes.Buffer
	is_toplevel := tag.Name == "blink:toplevel"

	if !is_toplevel {
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
	}

	if len(tag.Children) == 0 {
		if parser.CanSingle(tag.Name) {
			buffer.WriteString("/>")
		} else {
			buffer.WriteString("></" + tag.Name + ">")
		}

	} else {
		if !is_toplevel {
			buffer.WriteRune('>')
		}
		if tag.Name == "script" {
			var code bytes.Buffer
			for _, v := range tag.Children {
				code.WriteString(gen_code(v))
			}
			str_code := strike.StrikeJs(code.String())
			buffer.WriteString(str_code)
		} else {
			for _, v := range tag.Children {
				buffer.WriteString(gen_code(v))
			}
		}
		if !tag.IsDoctype && !is_toplevel {
			buffer.WriteString("</" + tag.Name + ">")
		}
	}
	return buffer.String()
}

func gen_code_doctype(tag *parser.Tag) string {
	var buffer bytes.Buffer
	buffer.WriteString("<" + tag.Name)
	for k, _ := range tag.Attribute {
		buffer.WriteString(" " + k)
	}
	//buffer.WriteRune('>')
	//fmt.Println(tag.Attribute)
	return buffer.String()
}

//--------- handler ------------

func get_html(html *parser.BlinkHtml) *parser.Tag {
	if html.Master != nil {
		buffer, err := ioutil.ReadFile(html.Master.Attribute["path"])
		if err != nil {
			storm.Error("\"" + html.Master.Attribute["path"] + "\" can not find master!")
		}

		master := parser.Parse(buffer)
		master = insert_content(html, master)
		return get_html(master)
	}
	return html.TopLevel
}

func insert_content(html, master *parser.BlinkHtml) *parser.BlinkHtml {
	for k, content := range master.Contents {
		ret := make([]*parser.Tag, 0)
		ret = append(ret, content.Parent.Children[:content.Index]...)
		if page, ok := html.PageArea[k]; ok {
			ret = append(ret, page.Children...)
		} else {
			ret = append(ret, content.Children...)
		}
		ret = append(ret, content.Parent.Children[content.Index+1:]...)
		content.Parent.Children = ret
	}

	return master
}
