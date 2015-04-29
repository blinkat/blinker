package process

import (
	par "github.com/blinkat/blinks/strike/js/parser"
)

type prepare_ifs struct {
	Default
}

func (p *prepare_ifs) redo_if(statements []par.IAst, w *Walker) []par.IAst {
	statements = map_(w, statements)

	for k, v := range statements {
		if v.Type() == par.Type_If {
			ifs := v.(*par.If)
			if ifs.Else == nil && aborts(ifs.Body) {
				cond := w.Walk(ifs.Cond)
				body := p.redo_if(statements[k+1:], w)

				var e par.IAst
				if len(body) == 1 {
					e = body[0]
				} else {
					e = par.NewBlock(body)
				}

				statements = append(statements[0:k], par.NewIf(cond, ifs.Body, e))
				return statements
			}
		}
	}
	return statements
}

func (p *prepare_ifs) redo_if_lambda(w *Walker, ast par.IAst) par.IAst {
	f := ast.(*par.Function)
	bodys := p.redo_if(f.Body.Exprs, w)
	return par.NewFunction(f.Type(), f.Name(), f.Args, par.NewFuncBody(bodys))
}

func (p *prepare_ifs) redo_if_block(w *Walker, ast par.IAst) par.IAst {
	f := ast.(*par.Block)
	if f.Statements != nil {
		return par.NewBlock(p.redo_if(f.Statements, w))
	} else {
		return par.NewBlock(nil)
	}
}

func (p *prepare_ifs) Defun(w *Walker, ast par.IAst) par.IAst {
	return p.redo_if_lambda(w, ast)
}

func (p *prepare_ifs) Function(w *Walker, ast par.IAst) par.IAst {
	return p.redo_if_lambda(w, ast)
}

func (p *prepare_ifs) Block(w *Walker, ast par.IAst) par.IAst {
	return p.redo_if_block(w, ast)
}

func (p *prepare_ifs) TopLevel(w *Walker, ast par.IAst) par.IAst {
	return par.NewToplevel(p.redo_if(ast.(*par.Toplevel).Statements, w)...)
}

func (p *prepare_ifs) Try(w *Walker, ast par.IAst) par.IAst {
	t := ast.(*par.Try)
	body := p.redo_if(t.Body, w)

	var c *par.Catch
	var f []par.IAst

	if t.Catchs != nil {
		c = par.NewCatch(t.Catchs.Name(), p.redo_if(t.Catchs.Body, w))
	}

	if t.Finally != nil {
		f = p.redo_if(t.Finally, w)
	}

	return par.NewTry(body, f, c)
}

func PrepareIfs(ast par.IAst, w *Walker) par.IAst {
	pre := &prepare_ifs{}
	w.foots = pre
	ret := w.Walk(ast)
	w.foots = nil
	return ret
}
