package process

import (
	"fmt"
	p "github.com/blinkat/blinker/strike/js/parser"
	"github.com/blinkat/blinker/strike/js/parser/adapter"
	"strconv"
)

type squeeze struct {
	Default
	scope *p.AstScope
}

func (s *squeeze) negate(c p.IAst, w *Walker) p.IAst {
	not_c := p.NewUnaryPrefix("!", c)

	switch c.Type() {
	case p.Type_Unary_Prefix:
		op := c.(*p.Unary)
		if c.Name() == "!" && boolean_expr(op.Expr) {
			return op.Expr
		} else {
			return not_c
		}
	case p.Type_Seq:
		ret := c.(*p.Seq)
		ret.Expr2 = s.negate(ret.Expr2, w)
		return ret
	case p.Type_Conditional:
		ret := c.(*p.Conditional)
		return best_of(not_c, p.NewConditional(ret.Expr, s.negate(ret.True, w), s.negate(ret.False, w)), GeneratorWalker(nil))

	case p.Type_Binnary:
		op := c.(*p.Binary)
		switch op.Name() {
		case "==":
			return p.NewBinary("!=", op.Left, op.Right)
		case "!=":
			return p.NewBinary("==", op.Left, op.Right)
		case "===":
			return p.NewBinary("!==", op.Left, op.Right)
		case "!==":
			return p.NewBinary("===", op.Left, op.Right)
		case "&&":
			return best_of(not_c, p.NewBinary("||", s.negate(op.Left, w), s.negate(op.Right, w)), GeneratorWalker(nil))
		case "||":
			return best_of(not_c, p.NewBinary("&&", s.negate(op.Left, w), s.negate(op.Right, w)), GeneratorWalker(nil))
		}
		break
	}
	return not_c
}

func (s *squeeze) make_conditional(c, t, e p.IAst, w *Walker) p.IAst {
	return WhenConstant(c, func(ast p.IAst, val interface{}) p.IAst {
		if val != nil {
			return t
		} else {
			return e
		}
	}, func(nothing p.IAst) p.IAst {
		if c.Type() == p.Type_Unary_Prefix && c.Name() == "!" {
			u := c.(*p.Unary)
			if e != nil {
				return p.NewConditional(u.Expr, e, t)
			} else {
				return p.NewBinary("||", u.Expr, t)
			}
		} else {
			if e != nil {
				return best_of(p.NewConditional(c, t, e), p.NewConditional(s.negate(c, w), e, t), GeneratorWalker(nil))
			} else {
				return p.NewBinary("&&", c, t)
			}
		}
	})
}

func (s *squeeze) rmblock(block *p.Block) p.IAst {
	var ret p.IAst
	if block != nil && block.Type() == p.Type_Block && block.Statements != nil {
		leng := len(block.Statements)
		if leng == 1 {
			ret = block.Statements[0]
		} else if leng == 0 {
			ret = p.NewBlock(nil)
		}
	}
	return ret
}

func (s *squeeze) _lambda(w *Walker, ast p.IAst) p.IAst {
	f := ast.(*p.Function)
	return p.NewFunction(ast.Type(), ast.Name(), f.Args, p.NewFuncBody(s.tighten(f.Body.Exprs, w)))
}

// this function does a few things:
// 1. discard useless blocks
// 2. join consecutive var declarations
// 3. remove obviously dead code
// 4. transform consecutive statements using the comma operator
// 5. if block_type == "lambda" and it detects constructs like if(foo) return ... - rewrite like if (!foo) { ... }
func (s *squeeze) tighten(statements []p.IAst, w *Walker) []p.IAst {
	ret := map_(w, statements)
	ret = s.tighten_clear_empty(ret)
	ret = s.tighten_each(ret)
	ret = s.tighten_clear_dead_code(ret)
	ret = s.tighten_make_seqs(ret)
	return ret
}

func (s *squeeze) tighten_clear_empty(statements []p.IAst) []p.IAst {
	ret := make([]p.IAst, 0)
	for _, v := range statements {
		if v.Type() == p.Type_Block {
			b := v.(*p.Block)
			if b.Statements != nil {
				ret = append(ret, b.Statements...)
			}
		} else {
			ret = append(ret, v)
		}
	}
	return ret
	/*ret := make([]p.IAst, 0)
	for _, v := range statements {
		if v.Type() == p.Type_Block {
			b := v.(*p.Block)
			if b.Statements != nil && len(b.Statements) > 0 {
				ret = append(ret, b)
			}
		} else {
			ret = append(ret, v)
		}
	}
	return ret*/
}

func (s *squeeze) tighten_make_seqs(statements []p.IAst) []p.IAst {
	var prev p.IAst
	ret := make([]p.IAst, 0)
	for _, v := range statements {
		if prev != nil && prev.Type() == p.Type_Stat && v.Type() == p.Type_Stat {
			stat := prev.(*p.Stat)
			stat.Statement = p.NewSeq(stat.Statement, v.(*p.Stat).Statement)
		} else {
			ret = append(ret, v)
			prev = v
		}
	}

	leng := len(ret)
	if leng >= 2 &&
		ret[leng-2].Type() == p.Type_Stat &&
		(ret[leng-1].Type() == p.Type_Return || ret[leng-1].Type() == p.Type_Thorw) {
		r := ret[leng-1].(*p.Return)
		if r.Expr != nil {
			stat := ret[leng-2].(*p.Stat)
			ret = ret[:leng-2]
			if r.Type() == p.Type_Return {
				ret = append(ret, p.NewReturn(p.NewSeq(stat.Statement, r.Expr)))
			} else {
				ret = append(ret, p.NewThrow(p.NewSeq(stat.Statement, r.Expr)))
			}
		}
	}
	return ret
}

func (s *squeeze) tighten_clear_dead_code(statements []p.IAst) []p.IAst {
	ret := make([]p.IAst, 0)
	has_quit := false
	for _, v := range statements {
		if has_quit {
			if v.Type() == p.Type_Func || v.Type() == p.Type_Defunc {
				ret = append(ret, v)
			} else if v.Type() == p.Type_Var || v.Type() == p.Type_Const {
				a := v.(*p.Var)

				a.Defs = map_def_custom(func(ast *p.VarDef) *p.VarDef {
					return p.NewDef(ast.Name(), nil)
				}, a.Defs)
				ret = append(ret, a)
			}
		} else {
			ret = append(ret, v)
			switch v.Type() {
			case p.Type_Return, p.Type_Thorw, p.Type_Break, p.Type_Coutinue:
				has_quit = true
			}
		}
	}
	return ret
	return statements
}

func (s *squeeze) tighten_each(a []p.IAst) []p.IAst {
	ret := make([]p.IAst, 0)
	var prev p.IAst

	for _, v := range a {
		if prev != nil && ((v.Type() == p.Type_Var && prev.Type() == p.Type_Var) ||
			(v.Type() == p.Type_Const && prev.Type() == p.Type_Const)) {
			c := v.(*p.Var)
			pr := prev.(*p.Var)
			pr.Defs = append(pr.Defs, c.Defs...)
		} else {
			ret = append(ret, v)
			prev = v
		}
	}
	return ret
}

//-----------[ make if ]------------
func (s *squeeze) make_if(w *Walker, ast p.IAst) p.IAst {
	ifs := ast.(*p.If)
	return WhenConstant(ifs.Cond, func(ast p.IAst, val interface{}) p.IAst {
		if val != nil {
			t := w.Walk(ifs.Body)
			if t == nil {
				return p.NewBlock(nil)
			} else {
				return t
			}
		}
		return nil
	}, func(nothing p.IAst) p.IAst {
		return s.make_real_if(w, ast)
		//return ast
	})
}

func (s *squeeze) abort_else(c, t, e p.IAst, w *Walker) p.IAst {
	ret := make([]p.IAst, 0)
	ret = append(ret, p.NewIf(s.negate(c, w), e, nil))
	if t.Type() == p.Type_Block {
		b := t.(*p.Block)
		if b.Statements != nil {
			ret = append(ret, b.Statements...)
		} else {
			ret = append(ret, t)
		}
	}
	return w.Walk(p.NewBlock(ret))
}

func (s *squeeze) make_real_if(w *Walker, ast p.IAst) p.IAst {
	ifs := ast.(*p.If)
	c := w.Walk(ifs.Cond)
	t := w.Walk(ifs.Body)
	e := w.Walk(ifs.Else)
	wk := GeneratorWalker(nil)

	if empty(e) && empty(t) {
		return p.NewStat(c)
	}

	if empty(t) {
		c = s.negate(c, w)
		t = e
		e = nil
	} else if empty(e) {
		e = nil
	} else {
		a := GenCode(c, wk)
		n := s.negate(c, w)
		b := GenCode(n, wk)

		if len(b) < len(a) {
			tmp := t
			t = e
			e = tmp
			c = n
		}
	}

	var ret p.IAst
	ret = p.NewIf(c, t, e)
	if t.Type() == p.Type_If {
		tif := t.(*p.If)
		if empty(tif.Else) && empty(e) {
			ret = best_of(ret, w.Walk(p.NewIf(p.NewBinary("&&", c, t.(*p.If).Cond), t.(*p.If).Body, nil)), wk)
		}
	} else if t.Type() == p.Type_Stat {
		if e != nil {
			if e.Type() == p.Type_Stat {
				ret = best_of(ret, p.NewStat(s.make_conditional(c, t.(*p.Stat).Statement, e.(*p.Stat).Statement, w)), wk)
			} else if aborts(e) {
				ret = s.abort_else(c, t, e, w)
			}
		} else {
			ret = best_of(ret, p.NewStat(s.make_conditional(c, t.(*p.Stat).Statement, nil, w)), wk)
		}
	} else if e != nil && t.Type() == e.Type() && (t.Type() == p.Type_Return || t.Type() == p.Type_Thorw) &&
		t.(*p.Return).Expr != nil && e.(*p.Return).Expr != nil {
		cond := s.make_conditional(c, t.(*p.Return).Expr, e.(*p.Return).Expr, w)
		if t.Type() == p.Type_Return {
			ret = best_of(ret, p.NewReturn(cond), wk)
		} else {
			ret = best_of(ret, p.NewThrow(cond), wk)
		}
	} else if e != nil && aborts(t) {
		arr := []p.IAst{p.NewIf(c, t, nil)}
		if e.Type() == p.Type_Block {
			if e.(*p.Block).Statements != nil {
				arr = append(arr, e.(*p.Block).Statements...)
			}
		} else {
			arr = append(arr, e)
		}
		ret = w.Walk(p.NewBlock(arr))
	} else if t != nil && aborts(e) {
		ret = s.abort_else(c, t, e, w)
	}
	return ret
}

func (s *squeeze) _do_while(cond, body p.IAst, w *Walker) p.IAst {
	return WhenConstant(cond, func(cond p.IAst, val interface{}) p.IAst {
		if val == nil {
			return p.NewBlock(nil)
		} else {
			return p.NewFor(p.Type_For, nil, nil, nil, w.Walk(body))
		}
	}, nil)
}

//----------[ step ]------------
func (s *squeeze) Sub(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Sub)
	if a.Ret.Type() == p.Type_String {
		name := a.Ret.Name()
		if adapter.IsIdentifier(name) {
			return p.NewDot(name, w.Walk(a.Expr))
		} else if is_number.MatchString(name) || name == "0" {
			val, _ := strconv.ParseInt(name, 10, 32)
			return p.NewSub(w.Walk(a.Expr), p.NewNumber(fmt.Sprint(val)))
		}
	}
	return nil
}

func (s *squeeze) If(w *Walker, ast p.IAst) p.IAst {
	return s.make_if(w, ast)
}

func (s *squeeze) TopLevel(w *Walker, ast p.IAst) p.IAst {
	return p.NewToplevel(s.tighten(ast.(*p.Toplevel).Statements, w)...)
}

func (s *squeeze) Switch(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Switch)

	cond := w.Walk(a.Expr)

	last := len(a.Cases) - 1
	cases := make([]*p.Case, 0)
	for k, v := range a.Cases {
		block := s.tighten(v.Body, w)
		if last == k && len(block) > 0 {
			node := block[len(block)-1]
			if node.Type() == p.Type_Break && node.Name() == "" {
				block = block[:len(block)-1]
			}
		}
		if v.Expr == nil {
			cases = append(cases, p.NewCase(nil, block))
		} else {
			cases = append(cases, p.NewCase(w.Walk(v.Expr), block))
		}
	}

	return p.NewSwitch(cond, cases)
}

func (s *squeeze) Function(w *Walker, ast p.IAst) p.IAst {
	return s._lambda(w, ast)
}

func (s *squeeze) Defun(w *Walker, ast p.IAst) p.IAst {
	return s._lambda(w, ast)
}

func (s *squeeze) Binary(w *Walker, ast p.IAst) p.IAst {
	b := ast.(*p.Binary)
	return WhenConstant(p.NewBinary(b.Name(), w.Walk(b.Left), w.Walk(b.Right)),
		func(c p.IAst, val interface{}) p.IAst {
			return best_of(w.Walk(c), c, GeneratorWalker(nil))
		}, func(this p.IAst) p.IAst {
			if b.Name() != "==" && b.Name() != "!=" {
				return nil
			}
			l := w.Walk(b.Left)
			r := w.Walk(b.Right)

			if l != nil && l.Type() == p.Type_Unary_Prefix && l.Name() == "!" && l.(*p.Unary).Expr.Type() == p.Type_Number {
				val := l.(*p.Unary).Expr.Name()
				num, _ := parse_number(val)
				switch num.(type) {
				case float64:
					fv := num.(float64)
					if fv == 0 {
						b.Left = p.NewNumber("0")
					} else {
						b.Left = p.NewNumber("1")
					}
					break
				case int:
					fv := num.(int)
					if fv == 0 {
						b.Left = p.NewNumber("0")
					} else {
						b.Left = p.NewNumber("1")
					}
					break
				}
			} else if r != nil && r.Type() == p.Type_Unary_Prefix && r.Name() == "!" && r.(*p.Unary).Expr.Type() == p.Type_Number {
				val := r.(*p.Unary).Expr.Name()
				num, _ := parse_number(val)
				switch num.(type) {
				case float64:
					fv := num.(float64)
					if fv == 0 {
						b.Right = p.NewNumber("0")
					} else {
						b.Right = p.NewNumber("1")
					}
					break
				case int:
					fv := num.(int)
					if fv == 0 {
						b.Right = p.NewNumber("0")
					} else {
						b.Right = p.NewNumber("1")
					}
					break
				}
			}
			return p.NewBinary(b.Name(), b.Left, b.Right)
		})
}

func (s *squeeze) Conditional(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Conditional)
	return s.make_conditional(c.Expr, c.True, c.False, w)
}

func (s *squeeze) Try(w *Walker, ast p.IAst) p.IAst {
	t := ast.(*p.Try)
	var c *p.Catch
	var b, f []p.IAst
	b = s.tighten(t.Body, w)

	if t.Catchs != nil {
		c = p.NewCatch(t.Catchs.Name(), s.tighten(t.Catchs.Body, w))
	}
	if t.Finally != nil {
		f = s.tighten(t.Finally, w)
	}

	return p.NewTry(b, f, c)
}

func (s *squeeze) UnaryPrefix(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Unary)
	expr := w.Walk(a.Expr)
	ret := p.NewUnaryPrefix(a.Name(), expr)
	if a.Name() == "!" {
		ret = best_of(ret, s.negate(expr, w), GeneratorWalker(nil)).(*p.Unary)
	}
	return WhenConstant(ret, func(ref p.IAst, val interface{}) p.IAst {
		return w.Walk(ref)

	}, func(ref p.IAst) p.IAst {
		return ret
	})
}

func (s *squeeze) Name(w *Walker, ast p.IAst) p.IAst {
	switch ast.Name() {
	case "true":
		return p.NewUnaryPrefix("!", p.NewNumber("0"))
	case "false":
		return p.NewUnaryPrefix("!", p.NewNumber("1"))
	}
	return nil
}

func (s *squeeze) While(w *Walker, ast p.IAst) p.IAst {
	d := ast.(*p.While)
	return s._do_while(d.Expr, d.Body, w)
}

func (s *squeeze) Assign(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Assign)
	lval := w.Walk(a.Left)
	rval := w.Walk(a.Right)

	if a.Name() == "true" && lval.Type() == p.Type_Name && rval.Type() == p.Type_Binnary &&
		member(ok_ops, rval.Name()) && rval.(*p.Binary).Left.Type() == p.Type_Name &&
		rval.(*p.Binary).Left.Name() == lval.Name() {
		return p.NewAssign(rval.Name(), lval, rval.(*p.Binary).Right)
	}
	return p.NewAssign(a.Name(), lval, rval)
}

func (s *squeeze) Call(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Call)
	expr := w.Walk(c.Expr)

	if expr.Type() == p.Type_Dot && expr.(*p.Dot).Expr.Type() == p.Type_String && expr.(*p.Dot).Expr.Name() == "toString" {
		return expr.(*p.Dot).Expr
	}
	return p.NewCall(expr, map_(w, c.List))
}

//-----------[ out func ]--------------
func Squeeze(ast p.IAst, w *Walker) p.IAst {
	ret := do_squeeze(ast, w)
	return do_squeeze(ret, w)
}

func do_squeeze(ast p.IAst, w *Walker) p.IAst {
	ret := PrepareIfs(ast, w)
	foot := &squeeze{}
	w.foots = foot
	ret = w.Walk(ret)
	w.foots = nil
	return ret
}
