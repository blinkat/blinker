package process

import (
	p "github.com/blinkat/blinks/strike/parser"
)

type Walker struct {
	def_foot *Default
	foots    Foots
	stack    []p.IAst
}

type WalkFunc func(w *Walker, ast p.IAst) p.IAst

func (w *Walker) Walk(ast p.IAst) p.IAst {
	if ast == nil {
		return nil
	}
	w.stack = append(w.stack, ast)
	defer w.pop()
	fn := w.call(ast, w.foots)
	if fn != nil {
		ret := fn(w, ast)
		if ret != nil {
			return ret
		}
	}
	fn = w.call(ast, w.def_foot)
	if fn != nil {
		return fn(w, ast)
	}
	return nil
}

func (w *Walker) Stack() []p.IAst {
	return w.stack
}

func (w *Walker) Parent() p.IAst {
	leng := len(w.stack)
	if leng >= 2 {
		return w.stack[leng-2]
	}
	return nil
}

func (w *Walker) call(ast p.IAst, foots Foots) WalkFunc {
	if foots == nil {
		return nil
	}

	switch ast.Type() {
	case p.Type_String:
		return foots.String
	case p.Type_Number:
		return foots.Number
	case p.Type_Name:
		return foots.Name
	case p.Type_TopLevel:
		return foots.TopLevel
	case p.Type_Block:
		return foots.Block
	case p.Type_Var:
		return foots.Var
	case p.Type_Const:
		return foots.Const
	case p.Type_Try:
		return foots.Try
	case p.Type_Thorw:
		return foots.Throw
	case p.Type_New:
		return foots.New
	case p.Type_Switch:
		return foots.Switch
	case p.Type_Break:
		return foots.Break
	case p.Type_Coutinue:
		return foots.Continue
	case p.Type_Conditional:
		return foots.Conditional
	case p.Type_Assign:
		return foots.Assign
	case p.Type_Dot:
		return foots.Dot
	case p.Type_Call:
		return foots.Call
	case p.Type_Func:
		return foots.Function
	case p.Type_Debugger:
		return foots.Debugger
	case p.Type_Defunc:
		return foots.Defun
	case p.Type_If:
		return foots.If
	case p.Type_For:
		return foots.For
	case p.Type_For_In:
		return foots.ForIn
	case p.Type_While:
		return foots.While
	case p.Type_Do:
		return foots.Do
	case p.Type_Return:
		return foots.Return
	case p.Type_Binnary:
		return foots.Binary
	case p.Type_Unary_Prefix:
		return foots.UnaryPrefix
	case p.Type_Unary_Postfix:
		return foots.UnaryPostfix
	case p.Type_Sub:
		return foots.Sub
	case p.Type_Object:
		return foots.Object
	case p.Type_Regexp:
		return foots.Regexp
	case p.Type_Array:
		return foots.Array
	case p.Type_Stat:
		return foots.Stat
	case p.Type_Seq:
		return foots.Seq
	case p.Type_Label:
		return foots.Label
	case p.Type_With:
		return foots.With
	case p.Type_Atom:
		return foots.Atom
	case p.Type_Directive:
		return foots.Directive
	}
	return nil
}

func (w *Walker) pop() {
	w.stack = w.stack[:len(w.stack)-1]
}

func GeneratorWalker(foots Foots) *Walker {
	w := &Walker{}
	w.foots = foots
	w.def_foot = &Default{}
	return w
}
