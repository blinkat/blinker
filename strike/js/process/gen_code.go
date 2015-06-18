package process

import (
	"bytes"
	"fmt"
	p "github.com/blinkat/blinker/strike/js/parser"
	"github.com/blinkat/blinker/strike/js/parser/adapter"
	"regexp"
	"strings"
)

type gen_code struct {
	Default
	newline     string
	space       string
	indentation int

	add_space1 *regexp.Regexp
	add_space2 *regexp.Regexp
}

func (g *gen_code) encode_string(str string) string {
	ret := make_string(str)
	return ret
}

func (g *gen_code) make_num(w *Walker, ast p.IAst) p.IAst {
	str := ast.Name()
	a := make([]string, 0)
	a = append(a, strings.Replace(make_num_e.ReplaceAllString(str, "."), "e+", "e", -1))
	num, _ := parse_number(str)

	switch num.(type) {
	case int:
		a = append(a, g.make_int(num.(int))...)
		break

	case float64:
		num_int := int(num.(float64))
		if float64(num_int)-num.(float64) == 0 {
			a = append(a, g.make_int(num_int)...)
		} else {
			m := make_num_match1.FindAllString(str, -1)
			if m != nil && len(m) >= 3 {
				a = append(a, m[1]+"e"+fmt.Sprint(len(m[2])))
			} else {
				m = make_num_match2.FindAllString(str, -1)
				if m != nil && len(m) >= 3 {
					rs := []rune(str)
					//ret := append(make([]rune, 0), m[2]
					a = append(a, m[2]+"e-"+fmt.Sprint(len(m[1])+len(m[2]))+string(rs[strings.Index(str, "."):]))
				}
			}
		}
	}

	ret := a[0]
	for _, v := range a {
		if len(v) < len(ret) {
			ret = v
		}
	}
	return p.NewString(ret)
}

func (g *gen_code) make_int(num int) []string {
	a := make([]string, 0)
	if num >= 0 {
		a = append(a, "0x"+strings.ToLower(fmt.Sprintf("%x", num)))
		a = append(a, "0"+fmt.Sprintf("%o", num))
	} else {
		a = append(a, "-0x"+strings.ToLower(fmt.Sprintf("%x", (-num))))
		a = append(a, "-0"+fmt.Sprintf("%o", (-num)))
	}

	return a
}

func (g *gen_code) must_has_semicolon(node p.IAst) bool {
	switch node.Type() {
	case p.Type_With, p.Type_While:
		w := node.(*p.While)
		return empty(w.Expr) || g.must_has_semicolon(w.Expr)
	case p.Type_For, p.Type_For_In:
		w := node.(*p.For)
		return empty(w.Body) || g.must_has_semicolon(w.Body)
	case p.Type_If:
		w := node.(*p.If)
		if empty(w.Body) && w.Else == nil {
			return true
		}
		if w.Else != nil {
			if empty(w.Else) {
				return true
			}
			return g.must_has_semicolon(w.Else)
		}
		return g.must_has_semicolon(w.Body)
	case p.Type_Directive:
		return true
	}
	return false
}

type indent_func func() p.IAst

func (g *gen_code) with_indent(cont indent_func, incr int) p.IAst {
	g.indentation += incr
	defer g.out_indent(incr)
	return cont()
}

func (g *gen_code) out_indent(i int) {
	g.indentation -= i
}

func (g *gen_code) add_spaces(a []string) string {
	/*b := make([]string, 0)
	for k, v := range a {
		b = append(b, v)
		if k+1 < len(a) {
			next := a[k+1]
			rv := ([]rune(v))
			last := rv[len(rv)-1]
			first := ([]rune(next))[0]

			if (adapter.IsIdentifierChar(last) &&
				(adapter.IsIdentifierChar(first) || first == '\\')) ||
				(g.add_space1.MatchString(v) && g.add_space2.MatchString(next) ||
					last == '/' && first == '/') {
				b = append(b, " ")
			}
		}
	}
	*/
	b := make([]string, 0)
	for k, v := range a {
		b = append(b, v)
		if k+1 < len(a) {
			next := a[k+1]
			cs_v := []rune(v)
			cs_n := []rune(next)
			if next != "" && ((adapter.IsIdentifierChar(cs_v[len(cs_v)-1]) &&
				(adapter.IsIdentifierChar(cs_n[0]) || cs_n[0] == '\\')) ||
				(g.add_space1.MatchString(v) &&
					g.add_space2.MatchString(next) ||
					cs_v[len(cs_v)-1] == '/' && cs_n[0] == '/')) {
				b = append(b, " ")

			}
		}
	}
	return strings.Join(b, "")
}

type parenthesize_func func(p.IAst) bool

func (g *gen_code) parenthesize_fn(expr p.IAst, w *Walker, fn parenthesize_func) string {
	gen := w.Walk(expr)
	if fn(expr) {
		return "(" + gen.Name() + ")"
	}
	return gen.Name()
}

func (g *gen_code) parenthesize_type(expr p.IAst, w *Walker, t int) string {
	gen := w.Walk(expr)
	if expr.Type() == t {
		return "(" + gen.Name() + ")"
	}
	return gen.Name()
}

func (g *gen_code) parenthesize(expr p.IAst, w *Walker, args ...interface{}) string {
	gen := w.Walk(expr).Name()

	for _, v := range args {
		switch v.(type) {
		case int:
			if expr.Type() == v.(int) {
				return "(" + gen + ")"
			}
		case parenthesize_func:
			if v.(parenthesize_func)(expr) {
				return "(" + gen + ")"
			}
		}
	}
	return gen
}

func (g *gen_code) needs_parens(w *Walker, expr p.IAst) bool {
	if expr.Type() == p.Type_Func || expr.Type() == p.Type_Object {
		arr := w.Stack()
		leng := len(arr) - 1

		if leng > 0 {
			self := arr[leng]
			leng -= 1
			par := arr[leng]

			for leng >= 0 {
				if par.Type() == p.Type_Stat {
					return true
				}
				if g.needs_parens_ifs1(par, self) || g.needs_parens_ifs2(par, self) {
					self = par
					leng -= 1
					if leng < 0 {
						break
					}
					par = arr[leng]
				} else {
					return false
				}
			}
		}
	}
	return !member_int(dot_call_no_parens, expr.Type())
}

func (g *gen_code) needs_parens_ifs1(ast p.IAst, self p.IAst) bool {
	if (ast.Type() == p.Type_Seq && ast.(*p.Seq).Expr1 == self) ||
		(ast.Type() == p.Type_Call && ast.(*p.Call).Expr == self) ||
		(ast.Type() == p.Type_Dot && ast.(*p.Dot).Expr == self) ||
		(ast.Type() == p.Type_Sub && ast.(*p.Sub).Expr == self) ||
		(ast.Type() == p.Type_Conditional && ast.(*p.Conditional).True == self) {
		return true
	}
	return false
}

func (g *gen_code) needs_parens_ifs2(ast p.IAst, self p.IAst) bool {
	if (ast.Type() == p.Type_Binnary && ast.(*p.Binary).Left == self) ||
		(ast.Type() == p.Type_Assign && ast.(*p.Assign).Left == self) ||
		(ast.Type() == p.Type_Unary_Postfix && ast.(*p.Unary).Expr == self) {
		return true
	}
	return false
}

//--------------[ makes ]--------------

func (g *gen_code) make_then(w *Walker, th p.IAst) p.IAst {
	if th == nil {
		return p.NewString(";")
	}

	th_arr := []p.IAst{th}
	if th.Type() == p.Type_Do {
		return g.make_block(w, th_arr)
	}
	b := th

	for {
		tp := b.Type()
		if tp == p.Type_If {
			ast := b.(*p.If)

			if ast.Else == nil {
				return w.Walk(p.NewBlock(th_arr))
			}
			b = ast.Else
		} else if tp == p.Type_While || tp == p.Type_Do {
			b = b.(*p.While).Body
		} else if tp == p.Type_For || tp == p.Type_For_In {
			b = b.(*p.For).Body
		} else {
			break
		}
	}
	return w.Walk(th)
}

func (g *gen_code) make_function(w *Walker, this p.IAst, name string, args, body []p.IAst, keyword string, no_parens bool) p.IAst {
	var bufs bytes.Buffer
	bufs.WriteString(keyword)
	if name != "" {
		bufs.WriteString(" " + name)
	}
	bufs.WriteRune('(')

	for k, v := range args {
		ret := v.Name()
		if k == len(args)-1 {
			bufs.WriteString(ret)
		} else {
			bufs.WriteString(ret + ",")
		}

	}
	bufs.WriteRune(')')
	out := g.add_spaces([]string{bufs.String(), g.make_block(w, body).Name()})

	if !no_parens && g.needs_parens(w, this) {
		return p.NewString("(" + out + ")")
	} else {
		return p.NewString(out)
	}
}

func (g *gen_code) make_block_statments(statements []p.IAst, w *Walker) p.IAst {
	last := len(statements) - 1
	var bufs bytes.Buffer

	for i, stat := range statements {
		code := w.Walk(stat).Name()
		if code != ";" {
			if i == last && !g.must_has_semicolon(stat) {
				code = make_block_code.ReplaceAllString(code, "")
			}
			bufs.WriteString(code + g.newline)
		}
	}

	return p.NewString(bufs.String())
}

func (g *gen_code) make_block(w *Walker, statements []p.IAst) p.IAst {
	//block := ast.(*p.Block)
	if len(statements) == 0 {
		return p.NewString("{}")
	}
	return p.NewString("{" + g.newline + g.make_block_statments(statements, w).Name() + g.newline + "}")
}

func (g *gen_code) make_vardef1(w *Walker, ast p.IAst) p.IAst {
	def := ast.(*p.VarDef)
	if def.Expr != nil {
		return p.NewString(g.add_spaces([]string{def.Name(), "=", g.parenthesize_type(def.Expr, w, p.Type_Seq)}) + g.newline)
	}
	return p.NewString(def.Name() + g.newline)
}

func (g *gen_code) make_switch_block(w *Walker, ast p.IAst) p.IAst {
	node := ast.(*p.Switch)
	n := len(node.Cases)
	if n == 0 {
		return p.NewString("{}")
	}

	var bufs bytes.Buffer
	bufs.WriteRune('{')

	for k, v := range node.Cases {
		has_body := len(v.Body) > 0
		code := g.make_switch_case(w, v)
		if has_body {
			code += g.make_block_statments(v.Body, w).Name() + g.newline
		}

		if has_body && k < n-1 {
			code += ";"
		}
		bufs.WriteString(code + g.newline)
	}
	bufs.WriteRune('}')
	return p.NewString(bufs.String())
}

func (g *gen_code) make_switch_case(w *Walker, v *p.Case) string {
	if v.Expr != nil {
		return g.add_spaces([]string{"case", w.Walk(v.Expr).Name(), ":"})
	} else {
		return "default:"
	}
}

//--------------[ foots ]-------------
func (g *gen_code) String(w *Walker, ast p.IAst) p.IAst {
	return p.NewString(g.encode_string(ast.Name()))
}

func (g *gen_code) Number(w *Walker, ast p.IAst) p.IAst {
	return g.make_num(w, ast)
}

func (g *gen_code) Name(w *Walker, ast p.IAst) p.IAst {
	return p.NewString(ast.Name())
}

func (g *gen_code) Debugger(w *Walker, ast p.IAst) p.IAst {
	return p.NewString("debugger;")
}

func (g *gen_code) TopLevel(w *Walker, ast p.IAst) p.IAst {
	return g.make_block_statments(ast.(*p.Toplevel).Statements, w)
}

func (g *gen_code) Block(w *Walker, ast p.IAst) p.IAst {
	return g.make_block(w, ast.(*p.Block).Statements)
}

func (g *gen_code) Var(w *Walker, ast p.IAst) p.IAst {
	return p.NewString("var " + g.add_commas(w, ast) + ";")
}

func (g *gen_code) Const(w *Walker, ast p.IAst) p.IAst {
	return p.NewString("const " + g.add_commas(w, ast) + ";")
}

func (g *gen_code) Try(w *Walker, ast p.IAst) p.IAst {
	try := ast.(*p.Try)
	out := []string{"try", g.make_block(w, try.Body).Name()}
	if try.Catchs != nil {
		out = append(out, "catch"+"("+try.Catchs.Name()+")"+g.make_block(w, try.Catchs.Body).Name())
	}
	if try.Finally != nil {
		out = append(out, "finally"+g.make_block(w, try.Finally).Name())
	}
	return p.NewString(g.add_spaces(out))
}

func (g *gen_code) Throw(w *Walker, ast p.IAst) p.IAst {
	return p.NewString(g.add_spaces([]string{"throw", w.Walk(ast.(*p.Return).Expr).Name()}) + ";")
}

func (g *gen_code) New(w *Walker, ast p.IAst) p.IAst {
	arr := make([]string, 0)
	arg := ast.(*p.New)
	var ret string

	if arg.Args != nil && len(arg.Args) > 0 {
		for _, v := range arg.Args {
			arr = append(arr, g.parenthesize_type(v, w, p.Type_Seq))
		}
		ret = "(" + strings.Join(arr, ",") + ")"
	}

	str := g.add_spaces([]string{"new", g.parenthesize(arg.Expr, w, p.Type_Seq, p.Type_Binnary, p.Type_Conditional, p.Type_Assign,
		func(a p.IAst) bool {
			foots := &news_foots{}
			foots.has_call = false
			wk := GeneratorWalker(nil)
			wk.foots = foots
			wk.Walk(a)
			return foots.has_call
		}) + ret})
	return p.NewString(str)
}

func (g *gen_code) Switch(w *Walker, ast p.IAst) p.IAst {
	s := ast.(*p.Switch)
	return p.NewString(g.add_spaces([]string{"switch", "(" + w.Walk(s.Expr).Name() + ")", g.make_switch_block(w, s).Name()}))
}

func (g *gen_code) Break(w *Walker, ast p.IAst) p.IAst {
	out := "break"
	if ast.Name() != "" {
		out += " " + ast.Name()
	}
	return p.NewString(out + ";")
}

func (g *gen_code) Continue(w *Walker, ast p.IAst) p.IAst {
	out := "continue"
	if ast.Name() != "" {
		out += " " + ast.Name()
	}
	return p.NewString(out + ";")
}

func (g *gen_code) Conditional(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Conditional)
	return p.NewString(g.add_spaces([]string{
		g.parenthesize(c.Expr, w, p.Type_Assign, p.Type_Seq, p.Type_Conditional),
		"?",
		g.parenthesize(c.True, w, p.Type_Seq),
		":",
		g.parenthesize(c.False, w, p.Type_Seq),
	}))
}

func (g *gen_code) Assign(w *Walker, ast p.IAst) p.IAst {
	op := ast.Name()
	if op != "" && op != "true" {
		op += "="
	} else {
		op = "="
	}
	a := ast.(*p.Assign)
	return p.NewString(g.add_spaces([]string{w.Walk(a.Left).Name(), op, g.parenthesize(a.Right, w, p.Type_Seq)}))
}

func (g *gen_code) Dot(w *Walker, ast p.IAst) p.IAst {
	d := ast.(*p.Dot)
	out := w.Walk(d.Expr).Name()

	if d.Expr.Type() == p.Type_Number {
		if !squeeze_dot.MatchString(out) {
			out += "."
		}
	} else if d.Expr.Type() != p.Type_Func && g.needs_parens(w, d.Expr) {
		out = "(" + out + ")"
	}

	if d.Name() != "" {
		out += "." + d.Name()
	}

	return p.NewString(out)
}

func (g *gen_code) Call(w *Walker, ast p.IAst) p.IAst {
	c := ast.(*p.Call)
	f := w.Walk(c.Expr).Name()
	ready := c.Expr.Type() == p.Type_Func && ([]rune(f))[0] == '('
	if !ready && g.needs_parens(w, c.Expr) {
		f = "(" + f + ")"
	}

	ret := make([]string, 0)
	for _, v := range c.List {
		ret = append(ret, g.parenthesize(v, w, p.Type_Seq))
	}
	return p.NewString(f + "(" + strings.Join(ret, ",") + ")")
}

func (g *gen_code) Function(w *Walker, ast p.IAst) p.IAst {
	f := ast.(*p.Function)
	return g.make_function(w, ast, f.Name(), f.Args, f.Body.Exprs, "function", false)
}

func (g *gen_code) Defun(w *Walker, ast p.IAst) p.IAst {
	f := ast.(*p.Function)
	return g.make_function(w, ast, f.Name(), f.Args, f.Body.Exprs, "function", false)
}

func (g *gen_code) If(w *Walker, ast p.IAst) p.IAst {
	f := ast.(*p.If)
	out := []string{"if", "(" + w.Walk(f.Cond).Name() + ")"}
	//if f.Else != nil && f.Else.Type() == p.Type_Block && len(f.Else.(*p.Block).Statements) != 0 {
	if f.Else != nil {
		out = append(out, g.make_then(w, f.Body).Name(), "else", w.Walk(f.Else).Name())
	} else {
		out = append(out, w.Walk(f.Body).Name())
	}
	return p.NewString(g.add_spaces(out))
}

func (g *gen_code) For(w *Walker, ast p.IAst) p.IAst {
	out := []string{"for"}
	f := ast.(*p.For)

	init, cond, step := g.get_for_init_string(w, f)
	init = squeeze_for.ReplaceAllString(init, ";")
	cond = squeeze_for.ReplaceAllString(cond, ";")
	step = squeeze_for.ReplaceAllString(step, "")

	args := init + cond + step
	if args == "; ; " {
		args = ";;"
	}
	out = append(out, "("+args+")", w.Walk(f.Body).Name())
	return p.NewString(g.add_spaces(out))
}

func (g *gen_code) ForIn(w *Walker, ast p.IAst) p.IAst {
	f := ast.(*p.For)
	out := []string{"for"}
	if f.Init != nil {
		out = append(out, "("+squeeze_for_in.ReplaceAllString(w.Walk(f.Init).Name(), ""))
	} else {
		out = append(out, "("+w.Walk(f.Cond).Name())
	}
	out = append(out, "in", w.Walk(f.Step).Name()+")", w.Walk(f.Body).Name())
	return p.NewString(g.add_spaces(out))
}

func (g *gen_code) While(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.While)
	return p.NewString(g.add_spaces([]string{"while", "(" + w.Walk(a.Expr).Name() + ")", w.Walk(a.Body).Name()}))
}

func (g *gen_code) Do(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Do)
	return p.NewString(g.add_spaces([]string{"do", w.Walk(a.Body).Name(), "while", "(" + w.Walk(a.Cond).Name() + ")"}) + ";")
}

func (g *gen_code) Return(w *Walker, ast p.IAst) p.IAst {
	out := []string{"return"}
	r := ast.(*p.Return)
	if r.Expr != nil {
		out = append(out, w.Walk(r.Expr).Name())
	}
	return p.NewString(g.add_spaces(out) + ";")
}

func (g *gen_code) Binary(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Binary)
	left := w.Walk(a.Left).Name()
	right := w.Walk(a.Right).Name()

	ti := a.Left.Type()
	if member_int(squeeze_binary_arr, ti) ||
		ti == p.Type_Binnary && adapter.Precedence(a.Name()) > adapter.Precedence(a.Left.Name()) ||
		ti == p.Type_Func && g.needs_parens(w, ast) {
		left = "(" + left + ")"
	}
	ti = a.Right.Type()
	if member_int(squeeze_binary_arr, ti) ||
		ti == p.Type_Binnary && adapter.Precedence(a.Name()) >= adapter.Precedence(a.Right.Name()) &&
			!(a.Right.Name() == a.Name() && member(squeeze_binary_arr2, a.Name())) {
		right = "(" + right + ")"
	}
	return p.NewString(g.add_spaces([]string{left, a.Name(), right}))
}

func (g *gen_code) UnaryPrefix(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Unary)
	val := w.Walk(a.Expr).Name()
	if !(a.Expr.Type() == p.Type_Number || (a.Expr.Type() == p.Type_Unary_Prefix && !adapter.Operator(ast.Name()+a.Expr.Name())) || !g.needs_parens(w, a.Expr)) {
		val = "(" + val + ")"
	}

	if adapter.IsAlphanumericChar(([]rune(ast.Name()))[0]) {
		return p.NewString(ast.Name() + " " + val)
	} else {
		return p.NewString(ast.Name() + val)
	}
}

func (g *gen_code) UnaryPostfix(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Unary)
	val := w.Walk(a.Expr).Name()
	if !(a.Expr.Type() == p.Type_Number || (a.Expr.Type() == p.Type_Unary_Prefix && !adapter.Operator(ast.Name()+a.Expr.Name())) || !g.needs_parens(w, a.Expr)) {
		val = "(" + val + ")"
	}

	return p.NewString(val + a.Name())
}

func (g *gen_code) Sub(w *Walker, ast p.IAst) p.IAst {
	s := ast.(*p.Sub)
	hash := w.Walk(s.Expr).Name()
	if g.needs_parens(w, s.Expr) {
		hash = "(" + hash + ")"
	}
	return p.NewString(hash + "[" + w.Walk(s.Ret).Name() + "]")
}

func (g *gen_code) Object(w *Walker, ast p.IAst) p.IAst {
	obj_need_parens := g.needs_parens(w, ast)
	obj := ast.(*p.Object)
	if len(obj.Propertys) == 0 {
		if obj_need_parens {
			return p.NewString("({})")
		} else {
			return p.NewString("{}")
		}
	}

	strs := make([]string, 0)
	for _, v := range obj.Propertys {
		if v.Oper == "get" || v.Oper == "set" {
			f := v.Expr.(*p.Function)
			return g.make_function(w, ast, v.Name(), f.Args, f.Body.Exprs, v.Oper, true)
		}
		key := v.Name()
		val := g.parenthesize(v.Expr, w, p.Type_Seq)
		if !adapter.IsIdentifier(key) {
			key = g.encode_string(key)
		}
		strs = append(strs, g.add_spaces([]string{key + ":", val}))
	}
	out := "{" + g.newline + strings.Join(strs, ","+g.newline) + g.newline + "}"
	if obj_need_parens {
		return p.NewString("(" + out + ")")
	}
	return p.NewString(out)
}

func (g *gen_code) Regexp(w *Walker, ast p.IAst) p.IAst {
	rx := ast.(*p.Regexp)
	return p.NewString("/" + rx.Name() + "/" + rx.Mode)
}

func (g *gen_code) Array(w *Walker, ast p.IAst) p.IAst {
	arr := ast.(*p.Array)
	if len(arr.List) == 0 {
		return p.NewString("[]")
	}

	strs := make([]string, 0)
	for k, v := range arr.List {
		if v.Type() == p.Type_Atom && v.Name() == "undefined" {
			if k == len(arr.List)-1 {
				strs = append(strs, ",")
			}
		}
		strs = append(strs, g.parenthesize(v, w, p.Type_Seq))
	}

	out := g.add_spaces([]string{"[", strings.Join(strs, ","), "]"})
	return p.NewString(out)
}

func (g *gen_code) Stat(w *Walker, ast p.IAst) p.IAst {
	s := ast.(*p.Stat)
	if s.Statement != nil {
		return p.NewString(squeeze_for.ReplaceAllString(w.Walk(s.Statement).Name(), ";"))
	} else {
		return p.NewString(";")
	}
}

func (g *gen_code) Seq(w *Walker, ast p.IAst) p.IAst {
	s := ast.(*p.Seq)
	s1 := w.Walk(s.Expr1).Name()
	s2 := w.Walk(s.Expr2).Name()
	return p.NewString(s1 + "," + s2)
}

func (g *gen_code) Label(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.Label)
	return p.NewString(g.add_spaces([]string{a.Name(), ":", w.Walk(a.Stat).Name()}))
}

func (g *gen_code) With(w *Walker, ast p.IAst) p.IAst {
	a := ast.(*p.While)
	return p.NewString(g.add_spaces([]string{"with", "(" + w.Walk(a.Expr).Name() + ")", w.Walk(a.Body).Name()}))
}

func (g *gen_code) Atom(w *Walker, ast p.IAst) p.IAst {
	return p.NewString(ast.Name())
}

func (g *gen_code) Directive(w *Walker, ast p.IAst) p.IAst {
	return p.NewString(make_string(ast.Name()) + ";")
}

func (g *gen_code) get_for_init_string(w *Walker, f *p.For) (string, string, string) {
	var init, cond, step string
	if f.Init != nil {
		init = w.Walk(f.Init).Name()
	}
	if f.Cond != nil {
		cond = w.Walk(f.Cond).Name()
	}
	if f.Step != nil {
		step = w.Walk(f.Step).Name()
	}
	return init, cond, step
}

//--------------[ out ]--------------
func (g *gen_code) add_commas(w *Walker, ast p.IAst) string {
	v := ast.(*p.Var)
	ret := make([]string, 0)
	for _, v := range v.Defs {
		ret = append(ret, g.make_vardef1(w, v).Name())
	}
	return strings.Join(ret, ",")
}

func new_gen_code() *gen_code {
	g := &gen_code{
		newline:    "",
		space:      "",
		add_space1: regexp.MustCompile("[\\+\\-]$"),
		add_space2: regexp.MustCompile("^[\\+\\-]"),
	}
	return g
}

func GenCode(ast p.IAst, w *Walker) string {
	g := new_gen_code()
	w.foots = g
	ret := w.Walk(ast)
	w.foots = nil
	return ret.Name()
}
