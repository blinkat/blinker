package html

type Tag struct {
	Name      string
	Parent    *Tag
	Children  []*Tag
	Attribute map[string]string
	IsString  bool
	IsDoctype bool
	Index     int
}

func (t *Tag) Add(tag *Tag) {
	tag.Index = len(t.Children)
	t.Children = append(t.Children, tag)
}

func NewTag(name string) *Tag {
	a := &Tag{
		Name:      name,
		Children:  nil,
		Attribute: nil,
		IsString:  false,
		IsDoctype: false,
	}
	return a
}
