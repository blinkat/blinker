package strike

import (
	"fmt"
	p "github.com/blinkat/blinks/strike/parser"
	w "github.com/blinkat/blinks/strike/process"
)

func Test(text string) {
	ast := p.ParseJs(text)
	wk := w.GeneratorWalker(nil)
	ret := w.AddScopeInfo(ast, wk)
	ret = w.MangleAst(ret, wk)
	fmt.Println(ret)
}
