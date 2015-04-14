package squeeze

import (
	"fmt"
	p "github.com/blinkat/blinks/strike/parser"
)

func Test() {
	ast := p.NewBinary("<", p.NewAtom(p.Type_String, "bbbb"), p.NewAtom(p.Type_String, "b"))
	val, err := evaluate(ast)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(val)
	}
}
