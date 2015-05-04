package js

import (
	p "github.com/blinkat/blinker/strike/js/parser"
	"github.com/blinkat/blinker/strike/js/process"
)

func Strike(text string) string {
	ast := p.ParseJs(text)
	wk := process.GeneratorWalker(nil)
	ast = process.AddScopeInfo(ast, wk)
	ast = process.MangleAst(ast, wk)
	ast = process.Squeeze(ast, wk)
	return process.GenCode(ast, wk)
}
