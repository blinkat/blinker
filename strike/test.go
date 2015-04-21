package strike

import (
	"fmt"
	p "github.com/blinkat/blinks/strike/parser"
	"github.com/blinkat/blinks/strike/process"
)

//w

func Test(text string) {
	ast := p.ParseJs(text)
	wk := process.GeneratorWalker(nil)
	ast = process.AddScopeInfo(ast, wk)
	ast = process.MangleAst(ast, wk)
	ast = process.Squeeze(ast, wk)
	fmt.Println(process.GenCode(ast, wk))
}
