package html

import (
	"bytes"
	"github.com/blinkat/blinker/storm"
	"strings"
)

//---------[ resault ]----------
type BlinkHtml struct {
	TopLevel *Tag
	Master   *Tag
	Contents map[string]*Tag
	PageArea map[string]*Tag
}

type Parser struct {
	scan        *scanner
	current_tag *Tag
	current_tok string
	res         *BlinkHtml
}

func Parse(source []byte) *BlinkHtml {
	p := NewParser(source)
	p.Parse()
	return p.res
}

func NewParser(source []byte) *Parser {
	source = tag_comment.ReplaceAll(source, nil)
	p := &Parser{
		scan:        new_scanner(source),
		current_tag: NewTag("blink:toplevel"),
		res:         &BlinkHtml{},
	}
	p.res.Contents = make(map[string]*Tag)
	p.res.PageArea = make(map[string]*Tag)
	return p
}

func (p *Parser) next() string {
	p.current_tok = string(p.scan.next_tag())
	return p.current_tok
}

func (p *Parser) Parse() *Tag {
	for p.statement() != nil {
	}
	p.res.TopLevel = p.current_tag
	return p.current_tag
}

//-----------[ step ]----------
func (p *Parser) statement() *Tag {
	tok := p.next()
	if tok == "" {
		return nil
	}

	switch string(tok) {
	case "<":
		return p.read_begin_tag()
	case "</":
		p.current_tag = p.current_tag.Parent
		p.next()
		p.next()
		return p.statement()
	default:
		if p.current_tag == nil {
			return nil
		}
		tag := NewTag(tok)
		/*if p.current_tag.Name != "script" {
			tag.Name = tag_white.ReplaceAllString(tag.Name, "")
		}*/
		tag.IsString = true
		p.current_tag.Add(tag)
		tag.Parent = p.current_tag
		return tag
	}
}

func (p *Parser) read_begin_tag() *Tag {
	name := strings.ToLower(p.next())
	ret := NewTag(name)

	if name == "!doctype" {
		var buffer bytes.Buffer
		buffer.WriteString(name)
		next := p.next()
		for next != ">" && next != "" {
			buffer.WriteString(" " + next)
			next = p.next()
		}
		ret.IsDoctype = true
		ret.Name = buffer.String()
	} else {
		attr := p.read_attributes()
		ret.Attribute = attr
	}

	if p.current_tag != nil {
		ret.Parent = p.current_tag
		p.current_tag.Add(ret)
	}

	if p.current_tok != "/>" && !ret.IsDoctype {
		p.current_tag = ret
	}

	if IsBtml(ret.Name) {
		return p.handle_custom_tag(ret)
	}
	return ret
}

func (p *Parser) skip_comment() {
	for {
		tok := p.next()
		if tok == "" {
			return
		}
		if tok == "--" {
			tok = p.next()
			if tok == ">" {
				return
			}
		}
	}
}

func (p *Parser) read_attributes() map[string]string {
	attrs := make(map[string]string)
	attr := strings.ToLower(string(p.next()))
	for attr != "" {
		if attr == ">" || attr == "/>" {
			break
		}
		val := p.next()
		if val != "=" {
			attrs[attr] = ""
			attr = val
			continue
		} else {
			attrs[attr] = p.next()
			attr = p.next()
		}
	}
	return attrs
}

func (p *Parser) handle_custom_tag(tag *Tag) *Tag {
	switch tag.Name {
	case "blink:master":
		p.res.Master = tag
	case "blink:page", "blink:content":
		if id, ok := tag.Attribute["id"]; ok && id != "" {
			if tag.Name == "blink:page" {
				p.res.PageArea[id] = tag
			} else {
				p.res.Contents[id] = tag
			}
		} else {
			storm.Warring(tag.Name + " need \"id\"")
		}
	}
	return tag
}
