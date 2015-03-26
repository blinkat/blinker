package process

import (
	"github.com/blinkat/blinks/strike/parser"
)

func init() {
	default_footstep = &def_step{}
}

// -----------[ default walker footstep ] -------------
type def_step struct {
}

func (d *def_step) Step() WalkerStep {

}

// -----------[ foots ] -----------------

func (d *def_step) vardefs(wk *Walker, ast *parser.Ast) *parser.Ast {

}

// -----------[ ast maps ] -------------
type MapFunc func(ast *parser.Ast, i int) *parser.Ast

func WalkerMap(ast *parser.Ast, fn MapFunc) {
	ret := make([]*parser.Ast, 0)
	top := make([]*parser.Ast, 0)
	for _, v := range ast.Attributes {
		val := fn(v, 0)

	}
}
