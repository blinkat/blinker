package scanner

import (
	"fmt"
)

func Test(text string) {
	t := generatorTokenizer(text)
	t.find_str("*/")
	ch := t.peek()
	fmt.Println("number:", ch, "\nstring:", string(ch), "\nposition:", t.pos)
}
