package strike

import (
	p "github.com/blinkat/blinks/strike/parser"
	"github.com/blinkat/blinks/strike/process"
)

func StrikeJs(text string) string {
	ast := p.ParseJs(text)
	wk := process.GeneratorWalker(nil)
	ast = process.AddScopeInfo(ast, wk)
	ast = process.MangleAst(ast, wk)
	ast = process.Squeeze(ast, wk)
	return process.GenCode(ast, wk)
}

func StrikeCss(text string) string {
	return string(strike_css([]byte(text)))
}
