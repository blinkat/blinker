package process

import "github.com/blinkat/blinker/strike/js/parser"

type astMangle struct {
	Default
	scope *parser.AstScope
}

func newMangle() *astMangle {
	a := &astMangle{}
	a.scope = nil

	return a
}

func (a *astMangle) get_mangled(name string) string {
	if a.scope.Parent == nil {
		return name
	}

	if v, ok := a.scope.Names[name]; ok && (v == parser.Type_Defunc || v == parser.Type_Lambda) {
		return name
	}

	return a.scope.GetMangled(name)
}

func (a *astMangle) _lambda(w *Walker, ast parser.IAst) parser.IAst {
	f := ast.(*parser.Function)
	is_defun := f.Type() == parser.Type_Defunc
	name := f.Name()
	extra := make(map[string]string)
	if name != "" {
		if is_defun {
			name = a.get_mangled(name)
		} else if f.Body.Scope().References(name) {
			if !(a.scope.UsesEval || a.scope.UsesWith) {
				extra[name] = a.scope.NextMangled()
				name = extra[name]
			} else {
				extra[name] = name
			}
		} else {
			name = ""
		}
	}
	args := make([]parser.IAst, 0)
	body := a.with_scope(f.Body.Scope(), func() parser.IAst {
		for _, v := range f.Args {
			n := a.get_mangled(v.Name())
			args = append(args, parser.NewString(n))
		}
		return parser.NewFuncBody(map_(w, f.Body.Exprs))
	}, extra).(*parser.FuncBody)
	return parser.NewFunction(f.Type(), name, args, body)
}

func (a *astMangle) with_scope(s *parser.AstScope, cont only_func, extra map[string]string) parser.IAst {
	_scope := a.scope
	a.scope = s

	if extra != nil {
		for k, v := range extra {
			s.SetMangled(k, v)
		}
	}

	for k, _ := range s.Names {
		a.get_mangled(k)
	}

	ret := cont()
	ret.SetScope(s)
	a.scope = _scope
	return ret
}

func (a *astMangle) vardefs(w *Walker, ast parser.IAst) parser.IAst {
	v := ast.(*parser.Var)
	defs := make([]*parser.VarDef, 0)
	for _, v := range v.Defs {
		defs = append(defs, parser.NewDef(a.get_mangled(v.Name()), w.Walk(v.Expr)))
	}
	return parser.NewVar(defs)
}

func (a *astMangle) _break(w *Walker, ast parser.IAst) parser.IAst {
	n := ast.Name()
	if n != "" {
		return parser.NewAtom(parser.Type_Break, a.scope.Labels.GetMangled(n))
	}
	return nil
}

//--------------[ walkers ]-----------------
func (a *astMangle) Function(w *Walker, ast parser.IAst) parser.IAst {
	return a._lambda(w, ast)
}

func (a *astMangle) Defun(w *Walker, ast parser.IAst) parser.IAst {
	ret := a._lambda(w, ast)
	switch w.Parent().Type() {
	case parser.Type_TopLevel, parser.Type_Func, parser.Type_Defunc:
		ret.SetAtTop(true)
	}
	return ret
}

func (a *astMangle) Label(w *Walker, ast parser.IAst) parser.IAst {
	r := ast.(*parser.Label)
	if a.scope.Labels.Refs[r.Name()] != nil {
		return parser.NewLabel(a.scope.Labels.GetMangled(r.Name()), w.Walk(r.Stat))
	}
	return w.Walk(r.Stat)
}

func (a *astMangle) Break(w *Walker, ast parser.IAst) parser.IAst {
	return a._break(w, ast)
}

func (a *astMangle) Continue(w *Walker, ast parser.IAst) parser.IAst {
	return a._break(w, ast)
}

func (a *astMangle) Var(w *Walker, ast parser.IAst) parser.IAst {
	return a.vardefs(w, ast)
}

func (a *astMangle) Const(w *Walker, ast parser.IAst) parser.IAst {
	return a.vardefs(w, ast)
}

func (a *astMangle) Name(w *Walker, ast parser.IAst) parser.IAst {
	return parser.NewAtom(ast.Type(), a.get_mangled(ast.Name()))
}

func (a *astMangle) Try(w *Walker, ast parser.IAst) parser.IAst {
	t := ast.(*parser.Try)
	ret := parser.NewTry(map_(w, t.Body), nil, nil)
	if t.Catchs != nil {
		ret.Catchs = parser.NewCatch(a.get_mangled(t.Catchs.Name()), map_(w, t.Catchs.Body))
	}
	if t.Finally != nil {
		ret.Finally = map_(w, t.Finally)
	}
	return ret
}

func (a *astMangle) TopLevel(w *Walker, ast parser.IAst) parser.IAst {
	t := ast.(*parser.Toplevel)

	return a.with_scope(t.Scope(), func() parser.IAst {
		return parser.NewToplevel(map_(w, t.Statements)...)
	}, nil)
}

func (a *astMangle) Directive(w *Walker, ast parser.IAst) parser.IAst {
	ast.SetAtTop(true)
	return ast
}

//-------------[ out interface ]------------
func MangleAst(ast parser.IAst, w *Walker) parser.IAst {
	m := newMangle()
	w.foots = m
	ret := w.Walk(ast)
	w.foots = nil
	return ret
}
