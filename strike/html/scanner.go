package html

import (
	"bytes"
	"github.com/blinkat/blinks/strike/js/parser/adapter"
)

type scanner struct {
	source     []byte
	pos        int
	length     int
	in_string  bool
	in_doctype bool
}

func new_scanner(content []byte) *scanner {
	s := &scanner{
		pos:        0,
		source:     content,
		length:     len(content),
		in_string:  false,
		in_doctype: false,
	}
	return s
}

//---------- [ out ] -----------
func (s *scanner) next_tag() []byte {
	s.jump_white()
	ch := s.peek()
	if ch == 0 {
		return nil
	}

	if s.in_string {
		if ch == '<' {
			s.in_string = false
			return s.next_tag()
		}
		//s.next()
		return s.read_content()
	}
	if ch == '/' {
		s.next()
		n := s.peek()
		if n == '>' {
			s.next()
			s.in_string = true
			return []byte{'/', '>'}
		}
	}
	if ch == '>' {
		if s.in_doctype {
			s.in_doctype = false
		}
		s.next()
		s.in_string = true
		return []byte{'>'}
	}
	if ch == '=' {
		s.next()
		return []byte{'='}
	}
	if ch == '"' {
		s.next()
		str := s.read_value()
		if s.in_doctype {
			return []byte("\"" + string(str) + "\"")
		}
		return str
	}
	if ch == '<' {
		s.next()
		ch = s.peek()
		if ch == '/' {
			s.next()
			return []byte{'<', '/'}
		}

		return []byte{'<'}
	}

	return s.read_string()
}

func (s *scanner) next() byte {
	if s.eof() {
		return 0
	}

	ch := s.source[s.pos]
	s.pos += 1
	return ch
}

func (s *scanner) eof() bool {
	return s.pos < 0 || s.pos >= s.length
}

func (s *scanner) peek() byte {
	if s.eof() {
		return 0
	}
	return s.source[s.pos]
}

func (s *scanner) jump_white() {
	ch := s.peek()
	for ch != 0 && adapter.WhitespaceChars(rune(ch)) {
		s.next()
		ch = s.peek()
	}
}

func (s *scanner) read_string() []byte {
	var buffer bytes.Buffer
	ch := s.peek()
	for ch != 0 && !is_end_name(ch) {
		//fmt.Println(string(ch), ":", !is_end_name(ch))
		buffer.WriteByte(ch)
		s.next()
		ch = s.peek()
	}
	bs := buffer.Bytes()
	if bytes.Contains(bytes.ToLower(bs), []byte("!doctype")) {
		s.in_doctype = true
	}
	return bs
}

func (s *scanner) read_value() []byte {
	var buffer bytes.Buffer
	ch := s.peek()
	for ch != 0 && ch != '"' {
		if ch == '\n' || ch == '\r' {
			s.next()
			ch = s.peek()
			continue
		}
		buffer.WriteByte(ch)
		s.next()
		ch = s.peek()
	}
	s.next()
	return buffer.Bytes()
}

func (s *scanner) read_end_to(end byte) []byte {
	var buffer bytes.Buffer
	ch := s.peek()
	for ch != 0 && ch != end {
		buffer.WriteByte(ch)
		s.next()
		ch = s.peek()
	}
	//s.next()
	return buffer.Bytes()
}

func (s *scanner) read_content() []byte {
	var buffer bytes.Buffer
	ch := s.peek()
	for ch != 0 && ch != '<' {
		buffer.WriteByte(ch)
		s.next()
		ch = s.peek()
	}
	s.in_string = false
	return buffer.Bytes()
}
