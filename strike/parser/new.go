package parser

func NewToplevel(as ...IAst) *Toplevel {
	t := &Toplevel{}
	t.t = Type_TopLevel
	t.Statements = make([]IAst, 0)
	t.Statements = append(t.Statements, as...)
	return t
}

func NewString(val string) IAst {
	a := &ast{}
	a.t = Type_String
	a.name = val
	return a
}

func NewLabel(v []string, stat IAst) *Label {
	a := &Label{}
	a.t = Type_Label
	a.Labels = v
	a.Stat = stat
	return a
}

func NewDirective(v string) *Directive {
	a := &Directive{}
	a.t = Type_Directive
	a.name = v
	return a
}

func NewBlock(b []IAst) *Block {
	r := &Block{}
	r.t = Type_Block
	r.Statements = b
	return r
}

func NewDebugger() *Debugger {
	d := &Debugger{}
	d.t = Type_Debugger
	return d
}

func NewDot(name string, expr IAst) *Dot {
	d := &Dot{}
	d.t = Type_Dot
	d.Expr = expr
	d.name = name
	return d
}

func NewSub(expr, ret IAst) *Sub {
	s := &Sub{}
	s.t = Type_Sub
	s.Expr = expr
	s.Ret = ret
	return s
}

func NewCall(expr IAst, list []IAst) *Call {
	c := &Call{}
	c.t = Type_Call
	c.Expr = expr
	c.List = list
	return c
}

func NewAtom(t int, name string) IAst {
	if t == Type_Regexp {
		r := &Regexp{}
		r.t = t
		r.name = name
		return r
	} else {
		a := &ast{}
		a.t = t
		a.name = name
		return a
	}
}

func NewUnaryPrefix(name string, expr IAst) *Unary {
	a := &Unary{}
	a.t = Type_Unary_Prefix
	a.name = name
	a.Expr = expr
	return a
}

func NewUnaryPostfix(name string, expr IAst) *Unary {
	a := &Unary{}
	a.t = Type_Unary_Postfix
	a.name = name
	a.Expr = expr
	return a
}

func NewSeq(e1, e2 IAst) *Seq {
	a := &Seq{}
	a.t = Type_Seq
	a.Expr1 = e1
	a.Expr2 = e2
	return a
}

func NewAssign(n string, l, r IAst) *Assign {
	a := &Assign{}
	a.t = Type_Assign
	a.name = n
	a.Left = l
	a.Right = r
	return a
}

func NewConditional(expr, y, f IAst) *Conditional {
	c := &Conditional{}
	c.t = Type_Conditional
	c.Expr = expr
	c.True = y
	c.False = f
	return c
}

func NewBinary(op string, l, r IAst) *Binary {
	c := &Binary{}
	c.t = Type_Binnary
	c.Left = l
	c.Right = r
	c.name = op
	return c
}

func NewStat(e IAst) *Stat {
	s := &Stat{}
	s.t = Type_Stat
	s.Statement = e
	return s
}

func NewNew(e IAst, args []IAst) *New {
	c := &New{}
	c.t = Type_New
	c.Expr = e
	c.Args = args
	return c
}

func NewArray(l []IAst) *Array {
	a := &Array{}
	a.t = Type_Array
	a.List = l
	return a
}

func NewGetSet(name, op string, expr IAst) *Property {
	a := &Property{}
	a.name = name
	a.Oper = op
	a.Expr = expr
	return a
}

func NewProperty(name string, expr IAst) *Property {
	a := &Property{}
	a.name = name
	a.Expr = expr
	a.Oper = "none"
	return a
}

func NewObject(p []*Property) *Object {
	a := &Object{}
	a.t = Type_Object
	a.Propertys = p
	return a
}

func NewFunction(t int, name string, args []IAst, b []IAst) *Function {
	f := &Function{}
	f.t = t
	f.name = name
	f.Args = args
	f.Body = b
	return f
}

func NewFor(t int, init, cond, step, body IAst) *For {
	f := &For{}
	f.Init = init
	f.Cond = cond
	f.Step = step
	f.Body = body
	f.t = t
	return f
}

func NewVar(defs []*VarDef) *Var {
	a := &Var{}
	a.t = Type_Var
	a.Defs = defs
	return a
}

func NewDef(name string, expr IAst) *VarDef {
	a := &VarDef{}
	a.name = name
	a.Expr = expr
	return a
}

func NewIf(c, b, e IAst) *If {
	f := &If{}
	f.Cond = c
	f.Body = b
	f.Else = e
	f.t = Type_If
	return f
}

func NewCase(expr IAst, b []IAst) *Case {
	c := &Case{}
	c.Expr = expr
	c.Body = b
	return c
}

func NewSwitch(e IAst, cases []*Case) *Switch {
	a := &Switch{}
	a.Expr = e
	a.Cases = cases
	a.t = Type_Switch
	return a
}

func NewCatch(name string, b []IAst) *Catch {
	c := &Catch{}
	c.name = name
	c.Body = b
	return c
}

func NewTry(b, f []IAst, c *Catch) *Try {
	t := &Try{}
	t.t = Type_Try
	t.Finally = f
	t.Body = b
	t.Catchs = c
	return t
}

func NewDo(cond, body IAst) *Do {
	d := &Do{}
	d.Cond = cond
	d.Body = body
	d.t = Type_Do
	return d
}

func NewReturn(expr IAst) *Return {
	r := &Return{}
	r.Expr = expr
	r.t = Type_Return
	return r
}

func NewThrow(expr IAst) *Return {
	r := &Return{}
	r.Expr = expr
	r.t = Type_Thorw
	return r
}

func NewWhile(expr, body IAst) *While {
	w := &While{}
	w.Expr = expr
	w.Body = body
	w.t = Type_While
	return w
}

func NewWith(expr, body IAst) *While {
	w := &While{}
	w.Expr = expr
	w.Body = body
	w.t = Type_With
	return w
}
