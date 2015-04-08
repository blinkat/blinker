package process

import (
	p "github.com/blinkat/blinks/strike/parser"
)

type Foots interface {
	String(w *Walker, ast p.IAst) p.IAst
	Number(w *Walker, ast p.IAst) p.IAst
	Name(w *Walker, ast p.IAst) p.IAst
	TopLevel(w *Walker, ast p.IAst) p.IAst
	Block(w *Walker, ast p.IAst) p.IAst
	Var(w *Walker, ast p.IAst) p.IAst
	Const(w *Walker, ast p.IAst) p.IAst
	Try(w *Walker, ast p.IAst) p.IAst
	Throw(w *Walker, ast p.IAst) p.IAst
	New(w *Walker, ast p.IAst) p.IAst
	Switch(w *Walker, ast p.IAst) p.IAst
	Break(w *Walker, ast p.IAst) p.IAst
	Continue(w *Walker, ast p.IAst) p.IAst
	Conditional(w *Walker, ast p.IAst) p.IAst
	Assign(w *Walker, ast p.IAst) p.IAst
	Dot(w *Walker, ast p.IAst) p.IAst
	Call(w *Walker, ast p.IAst) p.IAst
	Function(w *Walker, ast p.IAst) p.IAst
	Debugger(w *Walker, ast p.IAst) p.IAst
	Defun(w *Walker, ast p.IAst) p.IAst
	If(w *Walker, ast p.IAst) p.IAst
	For(w *Walker, ast p.IAst) p.IAst
	ForIn(w *Walker, ast p.IAst) p.IAst
	While(w *Walker, ast p.IAst) p.IAst
	Do(w *Walker, ast p.IAst) p.IAst
	Return(w *Walker, ast p.IAst) p.IAst
	Binary(w *Walker, ast p.IAst) p.IAst
	UnaryPrefix(w *Walker, ast p.IAst) p.IAst
	UnaryPostfix(w *Walker, ast p.IAst) p.IAst
	Sub(w *Walker, ast p.IAst) p.IAst
	Object(w *Walker, ast p.IAst) p.IAst
	Regexp(w *Walker, ast p.IAst) p.IAst
	Array(w *Walker, ast p.IAst) p.IAst
	Stat(w *Walker, ast p.IAst) p.IAst
	Seq(w *Walker, ast p.IAst) p.IAst
	Label(w *Walker, ast p.IAst) p.IAst
	With(w *Walker, ast p.IAst) p.IAst
	Atom(w *Walker, ast p.IAst) p.IAst
	Directive(w *Walker, ast p.IAst) p.IAst
}

//-------------[ catch ]--------------
func map_(w *Walker, arr []p.IAst) []p.IAst {
	if arr == nil {
		return nil
	}

	top := make([]p.IAst, 0)
	ret := make([]p.IAst, 0)
	for _, v := range arr {
		val := w.Walk(v)
		if val.AtTop() {
			top = append(top, val)
		} else {
			ret = append(ret, w.Walk(v))
		}

	}
	return append(top, ret...)
}

func map_catch(w *Walker, c *p.Catch) *p.Catch {
	if c == nil {
		return nil
	}

	return p.NewCatch(c.Name(), map_(w, c.Body))
}
