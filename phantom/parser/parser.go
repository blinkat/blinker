package parser

import (
	"bytes"
	"github.com/blinkat/blinks/storm"
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
		if p.current_tag.Name != "script" {
			tag.Name = tag_white.ReplaceAllString(tag.Name, "")
		}
		tag.IsString = true
		p.current_tag.Add(tag)
		tag.Parent = p.current_tag
		return tag
	}
}

func (p *Parser) read_begin_tag() *Tag {
	name := strings.ToLower(p.next())

	// comments
	if name == "!--" {
		p.skip_comment()
		return p.statement()
	}

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

	/*if p.current_tok == "/>" {
		p.current_tag.Add(ret)
		return ret
	}*/

	if p.current_tag != nil {
		ret.Parent = p.current_tag
		p.current_tag.Add(ret)
	}

	if p.current_tok != "/>" && !ret.IsDoctype {
		p.current_tag = ret
	}
	name = strings.ToLower(ret.Name)

	if name == "blink:master" {
		p.res.Master = ret
	} else if name == "blink:page" {
		if id, ok := ret.Attribute["id"]; ok && id != "" {
			p.res.PageArea[id] = ret
		} else {
			storm.Error("blink page or content need id")
		}
	} else if name == "blink:content" {
		if id, ok := ret.Attribute["id"]; ok && id != "" {
			p.res.Contents[id] = ret
		} else {
			storm.Error("blink page or content need id")
		}
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
