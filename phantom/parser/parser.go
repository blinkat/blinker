package parser

type Parser struct {
	scan        *scanner
	current_tag *Tag
	current_tok []byte
}

func NewParser(source []byte) *Parser {
	source = tag_comment.ReplaceAll(source, nil)
	p := &Parser{
		scan:        new_scanner(source),
		current_tag: NewTag([]byte("Blink::TopLevel")),
	}
	return p
}

func (p *Parser) next() []byte {
	p.current_tok = p.scan.next_tag()
	return p.current_tok
}

func (p *Parser) Parse() *Tag {
	for p.statement() != nil {
	}
	return p.current_tag
}

//-----------[ step ]----------
func (p *Parser) statement() *Tag {
	tok := p.scan.next_tag()
	if tok == nil {
		return nil
	}

	switch string(tok) {
	case "<":
		return p.read_begin_tag()
	case "</":
		p.current_tag = p.current_tag.Parent
		p.scan.next_tag()
		p.scan.next_tag()
		return p.statement()
	default:
		if p.current_tag == nil {
			return nil
		}
		tag := NewTag(tok)
		tag.Name = tag_white.ReplaceAll(tag.Name, nil)
		tag.IsString = true
		p.current_tag.Add(tag)
		tag.Parent = p.current_tag
		return tag
	}
}

func (p *Parser) read_begin_tag() *Tag {
	name := p.next()

	// comments
	if string(name) == "!--" {
		p.skip_comment()
		return p.statement()
	}

	ret := NewTag(name)

	attr := p.read_attributes()
	ret.Attribute = attr
	if string(p.current_tok) == "/>" {
		p.current_tag.Add(ret)
		return ret
	}

	if p.current_tag != nil {
		ret.Parent = p.current_tag
		p.current_tag.Add(ret)
	}
	p.current_tag = ret
	return ret
}

func (p *Parser) skip_comment() {
	for {
		tok := p.next()
		if tok == nil {
			return
		}
		if string(tok) == "--" {
			tok = p.next()
			if string(tok) == ">" {
				return
			}
		}
	}
}

func (p *Parser) read_attributes() []*Attribute {
	attrs := make([]*Attribute, 0)
	attr := p.next()
	for attr != nil {
		if string(attr) == ">" || string(attr) == "/>" {
			break
		}
		val := p.next()
		if string(val) != "=" {
			attrs = append(attrs, NewAttr(attr, nil))
			attr = val
			continue
		} else {
			attrs = append(attrs, NewAttr(attr, p.next()))
			attr = p.next()
		}
	}
	return attrs
}
