package parser

type Tag struct {
	Name      []byte
	Parent    *Tag
	Children  []*Tag
	Attribute []*Attribute
	IsString  bool
}

func (t *Tag) Add(tag *Tag) {
	t.Children = append(t.Children, tag)
}

type Attribute struct {
	Name  []byte
	Value []byte
}

func NewTag(name []byte) *Tag {
	a := &Tag{
		Name:      name,
		Children:  nil,
		Attribute: nil,
		IsString:  false,
	}
	return a
}

func NewAttr(name, val []byte) *Attribute {
	a := &Attribute{
		Name:  name,
		Value: val,
	}
	return a
}
