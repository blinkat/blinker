package scanner

import (
	"fmt"
)

func Test(text string) {
	t := generatorTokenizer(text)
	for !t.eof() {
		tok := t.next_token("")
		fmt.Println(tok)
	}
}
