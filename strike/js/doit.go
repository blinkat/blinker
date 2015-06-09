package js

import (
	"bytes"
	p "github.com/blinkat/blinker/strike/js/parser"
	"github.com/blinkat/blinker/strike/js/process"
	"io/ioutil"
)

func Strike(text string) string {
	ast := p.ParseJs(text)
	wk := process.GeneratorWalker(nil)
	ast = process.AddScopeInfo(ast, wk)
	ast = process.MangleAst(ast, wk)
	ast = process.Squeeze(ast, wk)
	return process.GenCode(ast, wk)
}

// combine some js file
func Combine(paths []string) string {
	var ret bytes.Buffer
	for v, _ := range paths {
		text, err := ioutil.ReadFile(v)
		if err == nil {
			ret.Write(text)
			ret.WriteRune('\n')
		}
	}
	return string(ret.Bytes())
}
