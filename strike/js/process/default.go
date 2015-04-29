package process

import (
	p "github.com/blinkat/blinks/strike/js/parser"
)

type Default struct {
}

//----------[ walker step ]------------
func (d *Default) String(w *Walker, ast p.IAst) p.IAst {
	return p.NewString(ast.Name())
}

func (d *Default) Number(w *Walker, ast p.IAst) p.IAst {
	return p.NewAtom(p.Type_Number, ast.Name())
}

func (d *Default) Name(w *Walker, ast p.IAst) p.IAst {
	return p.NewAtom(ast.Type(), ast.Name())
}

func (d *Default) TopLevel(w *Walker, ast p.IAst) p.IAst {
	top := ast.(*p.Toplevel)
	return p.NewToplevel(map_(w, top.Statements)...)
}

func _vardefs(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Var)
	arr := make([]*p.VarDef, 0)
	for _, v := range a.Defs {
		var expr p.IAst
		if v.Expr != nil {
			expr = w.Walk(v.Expr)
		}
		arr = append(arr, p.NewDef(v.Name(), expr))
	}
	if ast.Type() == p.Type_Const {
		return p.NewConst(arr)
	}
	return p.NewVar(arr)
}

func (d *Default) Block(w *Walker, ast p.IAst) p.IAst {
	b := ast.(*p.Block)
	return p.NewBlock(map_(w, b.Statements))
}

func (d *Default) Var(w *Walker, ast p.IAst) p.IAst {
	return _vardefs(w, ast)
}

func (d *Default) Const(w *Walker, ast p.IAst) p.IAst {
	return _vardefs(w, ast)
}

func (d *Default) Try(w *Walker, ast p.IAst) p.IAst {
	t := ast.(*p.Try)
	body := map_(w, t.Body)
	return p.NewTry(body,
		map_(w, t.Finally),
		map_catch(w, t.Catchs),
	)
}

func (d *Default) Throw(w *Walker, ast p.IAst) p.IAst {
	return p.NewThrow(w.Walk(ast.(*p.Return).Expr))
}

func (d *Default) New(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.New)
	return p.NewNew(w.Walk(a.Expr), map_(w, a.Args))
}

func (d *Default) Switch(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Switch)
	return p.NewSwitch(w.Walk(a.Expr), _case(w, a.Cases))
}

func _case(w *Walker, ast []*p.Case) []*p.Case {
	ret := make([]*p.Case, 0)
	for _, v := range ast {
		if v.Expr != nil {
			ret = append(ret, p.NewCase(w.Walk(v.Expr), map_(w, v.Body)))
		} else {
			ret = append(ret, p.NewCase(nil, map_(w, v.Body)))
		}
	}
	return ret
}

func (d *Default) Break(w *Walker, ast p.IAst) p.IAst {
	return p.NewAtom(p.Type_Break, ast.Name())
}

func (d *Default) Continue(w *Walker, ast p.IAst) p.IAst {
	return p.NewAtom(p.Type_Coutinue, ast.Name())
}

func (d *Default) Conditional(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Conditional)
	return p.NewConditional(w.Walk(c.Expr), w.Walk(c.True), w.Walk(c.False))
}

func (d *Default) Assign(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Assign)
	return p.NewAssign(c.Name(), w.Walk(c.Left), w.Walk(c.Right))
}

func (d *Default) Dot(w *Walker, ast p.IAst) p.IAst {
	return p.NewDot(ast.Name(), w.Walk(ast.(*p.Dot).Expr))
}

func (d *Default) Call(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Call)
	return p.NewCall(w.Walk(c.Expr), map_(w, c.List))
}

func (d *Default) Function(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Function)
	return p.NewFunction(c.Type(), c.Name(), map_(w, c.Args), p.NewFuncBody(map_(w, c.Body.Exprs)))
}

func (d *Default) Debugger(w *Walker, ast p.IAst) p.IAst {
	return p.NewDebugger()
}

func (d *Default) Defun(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Function)
	return p.NewFunction(c.Type(), c.Name(), map_(w, c.Args), p.NewFuncBody(map_(w, c.Body.Exprs)))
}

func (d *Default) If(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.If)
	return p.NewIf(w.Walk(c.Cond), w.Walk(c.Body), w.Walk(c.Else))
}

func (d *Default) For(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.For)
	return p.NewFor(ast.Type(), w.Walk(c.Init), w.Walk(c.Cond), w.Walk(c.Step), w.Walk(c.Body))
}

func (d *Default) ForIn(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.For)
	return p.NewFor(ast.Type(), w.Walk(c.Init), w.Walk(c.Cond), w.Walk(c.Step), w.Walk(c.Body))
}

func (d *Default) While(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.While)
	return p.NewWhile(w.Walk(c.Expr), w.Walk(c.Body))
}

func (d *Default) Do(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Do)
	return p.NewDo(w.Walk(c.Cond), w.Walk(c.Body))
}

func (d *Default) Return(w *Walker, ast p.IAst) p.IAst {
	return p.NewReturn(w.Walk(ast.(*p.Return).Expr))
}

func (d *Default) Binary(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Binary)
	return p.NewBinary(c.Name(), w.Walk(c.Left), w.Walk(c.Right))
}

func (d *Default) UnaryPrefix(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Unary)
	return p.NewUnaryPrefix(c.Name(), w.Walk(c.Expr))
}

func (d *Default) UnaryPostfix(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Unary)
	return p.NewUnaryPostfix(c.Name(), w.Walk(c.Expr))
}

func (d *Default) Sub(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Sub)
	return p.NewSub(w.Walk(c.Expr), w.Walk(c.Ret))
}

func (d *Default) Object(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Object)
	arr := make([]*p.Property, 0)
	for _, v := range c.Propertys {
		arr = append(arr, _property(w, v))
	}
	return p.NewObject(arr)
}

func _property(w *Walker, ast *p.Property) *p.Property {
	if ast.Oper == "none" {
		return p.NewProperty(ast.Name(), w.Walk(ast.Expr))
	} else {
		return p.NewGetSet(ast.Name(), ast.Oper, w.Walk(ast.Expr))
	}
}

func (d *Default) Regexp(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Regexp)
	return p.NewRegexp(c.Name(), c.Mode)
}

func (d *Default) Array(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Array)
	return p.NewArray(map_(w, c.List))
}

func (d *Default) Stat(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Stat)
	return p.NewStat(w.Walk(c.Statement))
}

func (d *Default) Seq(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Seq)
	return p.NewSeq(w.Walk(c.Expr1), w.Walk(c.Expr2))
}

func (d *Default) Label(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Label)
	return p.NewLabel(c.Name(), w.Walk(c.Stat))
}

func (d *Default) With(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.While)
	return p.NewWith(w.Walk(c.Expr), w.Walk(c.Body))
}

func (d *Default) Atom(w *Walker, ast p.IAst) p.IAst {
	return p.NewAtom(ast.Type(), ast.Name())
}

func (d *Default) Directive(w *Walker, ast p.IAst) p.IAst {
	return p.NewDirective(ast.Name())
}
