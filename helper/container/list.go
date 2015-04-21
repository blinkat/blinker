package container

import "bytes"

type StringNode struct {
	Data string
	Next *StringNode
	Prev *StringNode
}

type StringList struct {
	Head  *StringNode
	Last  *StringNode
	Count int
}

func (s *StringList) Push(val string) *StringNode {
	node := &StringNode{
		Data: val,
	}
	if s.Head == nil {
		s.Head = node
		s.Last = node
		node.Prev = nil
	} else {
		s.Last.Next = node
		node.Prev = s.Last
		s.Last = node
	}
	s.Count += 1
	return node
}

func (s *StringList) Get(index int) *StringNode {
	i := 0
	for node := s.Head; node != nil; node = node.Next {
		if i == index {
			return node
		}
		i += 1
	}
	return nil
}

func (s *StringList) Insert(index int, val string) *StringNode {
	new_node := &StringNode{
		Data: val,
	}
	node := s.Get(index)
	if node != nil {
		if node.Prev == nil {
			new_node.Next = s.Head
			s.Head.Prev = new_node
			s.Head = new_node
		} else {
			node.Prev.Next = new_node
			new_node.Prev = node.Prev
			new_node.Next = node
			node.Prev = new_node
		}
		s.Count += 1
		return new_node
	}
	panic("can't not find element")
}

func (s *StringList) Remove(index int) *StringNode {
	node := s.Get(index)
	if node != nil {
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
		s.Count -= 1
	}
	return node
}

func (s *StringList) Empty() bool {
	return s.Head == nil && s.Last == nil
}

func (s *StringList) ToString() string {
	var bufs bytes.Buffer
	for node := s.Head; node != nil; node = node.Next {
		bufs.WriteString(node.Data)
	}
	return bufs.String()
}

func (s *StringList) StringForSpace(space string) string {
	var bufs bytes.Buffer
	for node := s.Head; node != nil; node = node.Next {
		bufs.WriteString(node.Data + space)
	}
	return bufs.String()
}

func NewStringList() *StringList {
	s := &StringList{}
	s.Count = 0
	return s
}
