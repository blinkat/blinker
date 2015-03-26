package process

import (
	"github.com/blinkat/blinks/strike/parser"
)

var default_footstep FootStep

//------------- [ walker footstep interface ] ----------------
type FootStep interface {
	Step(ast *parser.Ast) WalkerStep
}

//------------- [ walker ] ------------
type Walker struct {
	steps FootStep
	stack []*parser.Ast
}

type WalkerStep func(walk *Walker, ast *parser.Ast) *parser.Ast

func (w *Walker) Walk(ast *parser.Ast, i int) *parser.Ast {
	if ast == nil {
		return nil
	}
	defer w.pop()
	w.stack = append(w.stack, ast)
	step := w.steps.Step(ast)
	if step != nil {
		ret := step(w, ast)
		if ret != nil {
			return ret
		}
	}
	step = default_footstep.Step(ast)
	return step(w, ast)
}

func (w *Walker) pop() {
	w.stack = w.stack[:len(w.stack)-1]
}

func (w *Walker) Dive(ast *parser.Ast) *parser.Ast {
	defer w.pop()
	w.stack = append(w.stack, ast)
	step = default_footstep.Step(ast)
	return step(w, ast)
}

func (w *Walker) Parent() *parser.Ast {
	leng := len(w.stack)
	if leng >= 2 {
		return w.stack[leng-2]
	}
	return nil
}

func (w *Walker) Stack() []*parser.Ast {
	return w.stack
}
