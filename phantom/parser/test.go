package parser

import "fmt"

func Test(text []byte) {
	p := NewParser(text)
	ret := p.Parse()
	//show_tag(ret, 0)
	//fmt.Println(string(ret.Children[0].Name))
	fmt.Println(string(ret.Name))
}

func show_tag(tag *Tag, i int) {
	for ind := 0; ind < i; ind++ {
		fmt.Print("	")
	}
	fmt.Print("tag: ", string(tag.Name), " attr: ")
	show_attr(tag)
	fmt.Print("\n")
	if tag.Children != nil {
		for _, v := range tag.Children {
			show_tag(v, i+1)
		}
	}
	for ind := 0; ind < i; ind++ {
		fmt.Print("	")
	}
	if !tag.IsString {
		fmt.Println("end: ", string(tag.Name))
	}
}

func show_attr(tag *Tag) {
	for _, v := range tag.Attribute {
		fmt.Print(string(v.Name), " = ", string(v.Value), ", ")
	}
}
