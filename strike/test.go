package strike

import (
	p "github.com/blinkat/blinks/strike/parser"
	"github.com/blinkat/blinks/strike/process"
	//"github.com/blinkat/blinks/strike/squeeze"
	"fmt"
)

//w

func Test(text string) {
	ast := p.ParseJs(text)
	wk := process.GeneratorWalker(nil)
	ast = process.AddScopeInfo(ast, wk)
	ast = process.MangleAst(ast, wk)
	fmt.Println(process.GenCode(ast, wk))

	//fmt.Println(ast.(*p.Toplevel).Statements[0].(*p.Var).Defs[0].Expr.(*p.Binary).Left.(*p.Binary).Left.Name())

	//wk := w.GeneratorWalker(nil)
	//ret := w.AddScopeInfo(ast, wk)
	//ret = w.MangleAst(ret, wk)
	//fmt.Println(ret)
	//squeeze.Test()
}
