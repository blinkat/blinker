package process

//add scope info
import "github.com/blinkat/blinker/strike/js/parser"

type scopeWalker struct {
	current     *parser.AstScope
	having_eval []*parser.AstScope
	Default
}

//--------------[ walkers ]-------------
type only_func func() parser.IAst

func (s *scopeWalker) with_only_scope(cont only_func) parser.IAst {
	s.current = parser.NewScope(s.current)
	s.current.Labels = parser.NewScope(nil)

	ret := cont()
	s.current.Body = append(s.current.Body, ret)
	ret.SetScope(s.current)
	s.current = s.current.Parent
	return ret
}

func (s *scopeWalker) define(name string, t int) string {
	return s.current.Define(name, t)
}

func (w *scopeWalker) Reference(name string) {
	w.current.Refs[name] = parser.TrueScope()
}

//---------------[ ctor ]---------------
func newScopeWalker() *scopeWalker {
	s := &scopeWalker{}
	s.current = nil
	s.having_eval = make([]*parser.AstScope, 0)
	return s
}

func (s *scopeWalker) _lambda(w *Walker, ast parser.IAst) parser.IAst {
	f := ast.(*parser.Function)
	is_defun := f.Type() == parser.Type_Defunc
	name := f.Name()
	if is_defun {
		name = s.define(name, parser.Type_Defunc)
	}
	return parser.NewFunction(ast.Type(), name, f.Args, s.with_only_scope(func() parser.IAst {
		if !is_defun {
			s.define(name, parser.Type_Lambda)
		}
		for _, v := range f.Args {
			s.define(v.Name(), parser.Type_Arg)
		}
		return parser.NewFuncBody(map_(w, f.Body.Exprs))
	}).(*parser.FuncBody))
}

func (s *scopeWalker) _break(w *Walker, ast parser.IAst) parser.IAst {
	n := ast.Name()
	if n != "" {
		s.current.Labels.Refs[n] = parser.TrueScope()
	}
	return nil
}

//-------------[ foots ]-------------
func (s *scopeWalker) Function(w *Walker, ast parser.IAst) parser.IAst {
	return s._lambda(w, ast)
}

func (s *scopeWalker) Defun(w *Walker, ast parser.IAst) parser.IAst {
	return s._lambda(w, ast)
}

func (s *scopeWalker) Label(w *Walker, ast parser.IAst) parser.IAst {
	s.current.Labels.Define(ast.Name(), parser.Type_Var)
	return nil
}

func (s *scopeWalker) Break(w *Walker, ast parser.IAst) parser.IAst {
	return s._break(w, ast)
}

func (s *scopeWalker) Continue(w *Walker, ast parser.IAst) parser.IAst {
	return s._break(w, ast)
}

func (s *scopeWalker) With(w *Walker, ast parser.IAst) parser.IAst {
	for sp := s.current; sp != nil; sp = sp.Parent {
		sp.UsesWith = true
	}
	return nil
}

func (s *scopeWalker) vardefs(w *Walker, ast parser.IAst, t int) parser.IAst {
	a := ast.(*parser.Var)

	for _, v := range a.Defs {
		s.define(v.Name(), t)
		if v.Expr != nil {
			s.Reference(v.Name())
		}
	}
	return nil
}

func (s *scopeWalker) Var(w *Walker, ast parser.IAst) parser.IAst {
	return s.vardefs(w, ast, parser.Type_Var)
}

func (s *scopeWalker) Const(w *Walker, ast parser.IAst) parser.IAst {
	return s.vardefs(w, ast, parser.Type_Const)
}
func (s *scopeWalker) Try(w *Walker, ast parser.IAst) parser.IAst {
	a := ast.(*parser.Try)
	if a.Catchs != nil {
		body := map_(w, a.Body)

		return parser.NewTry(
			body,
			map_(w, a.Finally),
			parser.NewCatch(s.define(a.Catchs.Name(), parser.Type_Catch), map_(w, a.Catchs.Body)),
		)
	}
	return nil
}

func (s *scopeWalker) Name(w *Walker, ast parser.IAst) parser.IAst {
	if ast.Name() == "eval" {
		s.having_eval = append(s.having_eval, s.current)
	}
	s.Reference(ast.Name())
	return nil
}

func AddScopeInfo(ast parser.IAst, w *Walker) parser.IAst {
	scope := newScopeWalker()
	ret := scope.with_only_scope(func() parser.IAst {
		w.foots = scope
		ret := w.Walk(ast)
		w.foots = nil

		temp_scope := scope.current
		for _, v := range scope.having_eval {
			if v.Has("eval") == nil {
				for temp_scope != nil {
					temp_scope.UsesEval = true
					temp_scope = temp_scope.Parent
				}
			}
		}

		fixrefs(scope.current)
		return ret
	})

	return ret
}

func fixrefs(sc *parser.AstScope) {

	for i := len(sc.Children) - 1; i >= 0; i-- {
		fixrefs(sc.Children[i])
	}

	for k, _ := range sc.Refs {
		origin := sc.Has(k)
		for s := sc; s != nil; s = s.Parent {
			s.Refs[k] = origin
			if s == origin {
				break
			}
		}
	}
}
