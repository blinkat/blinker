package parser

import (
	"fmt"
	"github.com/blinkat/blinker/strike/js/parser/adapter"
	"github.com/blinkat/blinker/strike/js/parser/scanner"
)

type jsparser struct {
	input         scanner.Tokenizer
	token         *scanner.Token
	prev          *scanner.Token
	peeked        *scanner.Token
	in_func       int
	in_directives bool
	in_loop       int
	labels        []string

	//-------debuger----------
	count int
}

func (j *jsparser) Parse() IAst {
	ret := NewToplevel()
	for j.input.Eof() {
		ret.Statements = append(ret.Statements, j.statement())
	}
	return ret
}

func (j *jsparser) is_token(t int, value string) bool {
	if j.token != nil {
		return j.is(j.token, t, value)
	}
	return false
}

func (j *jsparser) is(tok *scanner.Token, t int, val string) bool {
	return t == tok.Type && (val == "" || val == tok.Value)
}

func (j *jsparser) peek() *scanner.Token {
	if j.peeked == nil {
		j.peeked = j.input.Next("")
	}
	return j.peeked
}

func (j *jsparser) next() *scanner.Token {
	j.prev = j.token
	if j.peeked != nil {
		j.token = j.peeked
		j.peeked = nil
	} else {
		j.token = j.input.Next("")
	}
	j.in_directives = j.in_directives && (j.token.Type == scanner.TokenString || j.is_token(scanner.TokenPunc, ";"))
	j.count += 1
	return j.token
}

func (j *jsparser) throw(msg string) {
	if msg == "" {
		msg = "Unexpected token"
	}
	s := fmt.Sprint(msg, "\ntoken:", j.token.Value, "\n", j.prev)
	panic(s)
}

func (j *jsparser) expect_token(ty int, val string) *scanner.Token {
	if j.is_token(ty, val) {
		return j.next()
	}
	j.throw("")
	return nil
}

func (j *jsparser) expect(val string) *scanner.Token {
	return j.expect_token(scanner.TokenPunc, val)
}

func (j *jsparser) can_insert_semicolon() bool {
	if j.token == nil {
		return true
	}
	return j.token.Nlb || !j.input.Eof() || j.is_token(scanner.TokenPunc, "}")
}

func (j *jsparser) labeled_statement(label string) IAst {
	j.labels = append(j.labels, label)
	stat := j.statement()
	j.labels = j.labels[:len(j.labels)-1]
	return NewLabel(label, stat)
}

//-------------[ statement ]---------------
func (j *jsparser) statement() IAst {
	if j.is_token(scanner.TokenOperator, "/") || j.is_token(scanner.TokenOperator, "/=") {
		j.peeked = nil
		rs := []rune(j.token.Value)
		j.token = j.input.Next(string(rs[1:]))
	}

	switch j.token.Type {
	case scanner.TokenString:
		return j.read_string()
	case scanner.TokenNumber, scanner.TokenRegexp, scanner.TokenOperator, scanner.TokenAtom:
		return j.simple_statement()
	case scanner.TokenName:
		return j.read_name()
	case scanner.TokenPunc:
		return j.read_punc()
	case scanner.TokenKeyword:
		return j.read_keyword()
	}
	j.throw("")
	return nil
}

func (j *jsparser) read_string() IAst {
	dir := j.in_directives
	stat := j.simple_statement().(*Stat)

	if dir && stat.Statement.Type() == Type_String && !j.is_token(scanner.TokenPunc, ",") {
		return NewDirective(stat.Statement.Name())
	}
	return stat
}

func (j *jsparser) read_punc() IAst {
	switch j.token.Value {
	case "{":
		return NewBlock(j.block_())
	case "[", "(":
		return j.simple_statement()
	case ";":
		j.next()
		return NewBlock(nil)
	default:
		j.throw("")
		return nil
	}
}

func (j *jsparser) read_keyword() IAst {
	val := j.token.Value
	j.next()
	switch val {
	case "break":
		return j.break_cont(Type_Break)
	case "continue":
		return j.break_cont(Type_Coutinue)
	case "debugger":
		j.semicolon()
		return NewDebugger()
	case "do":
		body := j.loop(j.statement)
		j.expect_token(scanner.TokenKeyword, "while")
		cond := j.parenthesised()
		j.semicolon()
		return NewDo(cond, body)
	case "for":
		return j.for_()
	case "function":
		return j.function_(true)
	case "if":
		return j.if_()
	case "return":
		if j.in_func == 0 {
			j.throw("'return' outside of function")
		}
		if j.is_token(scanner.TokenPunc, ";") {
			j.next()
			return NewReturn(nil)
		} else if j.can_insert_semicolon() {
			return NewReturn(nil)
		} else {
			a := NewReturn(j.expression(true, false))
			j.semicolon()
			return a
		}
	case "switch":
		return NewSwitch(j.parenthesised(), j.switch_block_())
	case "throw":
		if j.token.Nlb {
			j.throw("Illegal newline after 'throw'")
		}
		ret := j.expression(true, false)
		j.semicolon()
		return NewThrow(ret)

	case "try":
		return j.try_()
	case "var":
		ret := j.var_(false)
		j.semicolon()
		return ret
	case "const":
		ret := j.const_()
		j.semicolon()
		return ret
	case "while":
		return NewWhile(j.parenthesised(), j.loop(j.statement))
	case "with":
		return NewWith(j.parenthesised(), j.statement())
	default:
		j.throw("")
		return nil
	}
}

func (j *jsparser) read_name() IAst {
	if j.is(j.peek(), scanner.TokenPunc, ":") {
		val := j.token.Value
		j.next()
		j.next()
		return j.labeled_statement(val)
	} else {
		return j.simple_statement()
	}
}

func (j *jsparser) semicolon() {
	if j.is_token(scanner.TokenPunc, ";") {
		j.next()
	} else if !j.can_insert_semicolon() {
		j.throw("")
	}
}

//------------[ expression ]---------------
func (j *jsparser) expr_atom(allow_calls bool) IAst {
	if j.is_token(scanner.TokenOperator, "new") {
		j.next()
		return j.new_()
	}
	if j.is_token(scanner.TokenPunc, "") {
		switch j.token.Value {
		case "(":
			j.next()
			ret := j.expression(true, false)
			j.expect(")")
			return j.subscripts(ret, allow_calls)
		case "[":
			j.next()
			return j.subscripts(j.array_(), allow_calls)
		case "{":
			j.next()
			return j.subscripts(j.object_(), allow_calls)
		}
		j.throw("")
	}
	if j.is_token(scanner.TokenKeyword, "function") {
		j.next()
		return j.subscripts(j.function_(false), allow_calls)
	}
	if scanner.IsAtomStartToken(j.token) {
		atom := NewAtom(TokenTypeToAstType(j.token.Type), j.token.Value)
		if atom.Type() == Type_Regexp {
			ret := atom.(*Regexp)
			ret.Mode = j.token.Attributes[0]
		}
		j.next()
		return j.subscripts(atom, allow_calls)
	}
	j.throw("")
	return nil
}

func (j *jsparser) subscripts(expr IAst, allow_calls bool) IAst {
	if j.is_token(scanner.TokenPunc, ".") {
		j.next()
		return j.subscripts(NewDot(j.as_name(), expr), allow_calls)
	}
	if j.is_token(scanner.TokenPunc, "[") {
		j.next()
		ret := j.expression(true, false)
		j.expect("]")
		return j.subscripts(NewSub(expr, ret), allow_calls)
	}
	if allow_calls && j.is_token(scanner.TokenPunc, "(") {
		j.next()
		return j.subscripts(NewCall(expr, j.expr_list(")", false, false)), true)
	}
	return expr
}

func (j *jsparser) maybe_unary(allow_calls bool) IAst {
	if j.is_token(scanner.TokenOperator, "") && adapter.UnaryPrefix(j.token.Value) {
		ret := j.token.Value
		j.next()
		return NewUnaryPrefix(ret, j.maybe_unary(allow_calls))
	}
	val := j.expr_atom(allow_calls)
	for j.is_token(scanner.TokenOperator, "") && adapter.UnaryPostfix(j.token.Value) && !j.token.Nlb {
		val = NewUnaryPostfix(j.token.Value, val)
		j.next()
	}
	return val
}

func (j *jsparser) as_name() string {
	switch j.token.Type {
	case scanner.TokenName, scanner.TokenOperator, scanner.TokenKeyword, scanner.TokenAtom:
		ret := j.token.Value
		j.next()
		return ret
	default:
		j.throw("")
	}
	return ""
}

func (j *jsparser) expression(commas, no_in bool) IAst {
	expr := j.maybe_assign(no_in)
	if commas && j.is_token(scanner.TokenPunc, ",") {
		j.next()
		return NewSeq(expr, j.expression(true, no_in))
	}
	return expr
}

func (j *jsparser) expr_list(end string, allow_trailing_comma, allow_empty bool) []IAst {
	first := true
	a := make([]IAst, 0)
	for !j.is_token(scanner.TokenPunc, end) {
		if first {
			first = false
		} else {
			j.expect(",")
		}
		if allow_trailing_comma && j.is_token(scanner.TokenPunc, end) {
			break
		}
		if j.is_token(scanner.TokenPunc, ",") && allow_empty {
			a = append(a, NewAtom(Type_Atom, "undefined"))
		} else {
			a = append(a, j.expression(false, false))
		}
	}
	j.next()
	return a
}

func (j *jsparser) maybe_assign(no_in bool) IAst {
	left := j.maybe_conditional(no_in)
	if j.token != nil {
		val := j.token.Value
		ass := adapter.Assignment(val)
		if j.is_token(scanner.TokenOperator, "") && ass != "" {
			j.next()
			return NewAssign(ass, left, j.maybe_assign(no_in))
		}
	}
	return left
}

func (j *jsparser) expr_op(left IAst, min_prec int, no_in bool) IAst {
	op := ""
	if j.is_token(scanner.TokenOperator, "") {
		op = j.token.Value
	}
	if op != "" && op == "in" && no_in {
		op = ""
	}

	prec := adapter.Precedence(op)
	if prec != -1 && prec > min_prec {
		j.next()
		right := j.expr_op(j.maybe_unary(true), prec, no_in)
		return j.expr_op(NewBinary(op, left, right), min_prec, no_in)
	}
	return left
}

func (j *jsparser) maybe_conditional(no_in bool) IAst {
	expr := j.expr_op(j.maybe_unary(true), 0, no_in)
	if j.is_token(scanner.TokenOperator, "?") {
		j.next()
		yes := j.expression(false, false)
		j.expect(":")
		return NewConditional(expr, yes, j.expression(false, no_in))
	}
	return expr
}

//loop
type loop_fn func() IAst

func (j *jsparser) loop(fn loop_fn) IAst {
	j.in_loop += 1
	ret := fn()
	j.in_loop -= 1
	return ret
}

func (j *jsparser) simple_statement() IAst {
	ret := j.expression(true, false)
	j.semicolon()
	return NewStat(ret)
}

//----------------[ types ]-----------------
func (j *jsparser) new_() IAst {
	newxp := j.expr_atom(false)
	var args []IAst
	if j.is_token(scanner.TokenPunc, "(") {
		j.next()
		args = j.expr_list(")", false, false)
	}
	return j.subscripts(NewNew(newxp, args), true)
}

func (j *jsparser) array_() IAst {
	return NewArray(j.expr_list("]", true, true))
}

func (j *jsparser) object_() IAst {
	first := true
	a := make([]*Property, 0)
	for !j.is_token(scanner.TokenPunc, "}") {
		if first {
			first = false
		} else {
			j.expect(",")
		}
		if j.is_token(scanner.TokenPunc, "}") {
			break
		}
		t := j.token.Type
		name := j.as_property_name()
		if t == scanner.TokenName && (name == "get" || name == "set") && !j.is_token(scanner.TokenPunc, ":") {
			a = append(a, NewGetSet(j.as_name(), name, j.function_(false)))
		} else {
			j.expect(":")
			a = append(a, NewProperty(name, j.expression(false, false)))
		}
	}
	j.next()
	return NewObject(a)
}

func (j *jsparser) block_() []IAst {
	j.expect("{")
	a := make([]IAst, 0)
	for !j.is_token(scanner.TokenPunc, "}") {
		if !j.input.Eof() {
			j.throw("")
		}
		a = append(a, j.statement())
	}
	j.next()
	return a
}

func (j *jsparser) function_(in_statement bool) IAst {
	name := ""
	if j.is_token(scanner.TokenName, "") {
		name = j.token.Value
		j.next()
	}

	if in_statement && name == "" {
		j.throw("")
	}
	j.expect("(")
	ty := Type_Func
	if in_statement {
		ty = Type_Defunc
	}
	return NewFunction(ty, name, j.function_params(), NewFuncBody(j.function_block()))
}

func (j *jsparser) function_params() []IAst {
	first := true
	a := make([]IAst, 0)
	for !j.is_token(scanner.TokenPunc, ")") {
		if first {
			first = false
		} else {
			j.expect(",")
		}
		if !j.is_token(scanner.TokenName, "") {
			j.throw("")
		}
		a = append(a, NewString(j.token.Value))
		j.next()
	}
	j.next()
	return a
}

func (j *jsparser) function_block() []IAst {
	j.in_func += 1
	loop := j.in_loop
	j.in_directives = true
	j.in_loop = 0
	a := j.block_()
	j.in_func -= 1
	j.in_loop = loop
	return a
}

func (j *jsparser) as_property_name() string {
	switch j.token.Type {
	case scanner.TokenNumber, scanner.TokenString:
		ret := j.token.Value
		j.next()
		return ret
	}
	return j.as_name()
}

func (j *jsparser) break_cont(t int) IAst {
	name := ""
	if !j.can_insert_semicolon() {
		if j.is_token(scanner.TokenName, "") {
			name = j.token.Value
		}
	}
	if name != "" {
		j.next()
		if !j.rember(name) {
			j.throw(fmt.Sprint("Label", name, "without matching loop or statement"))
		}
	} else if j.in_loop == 0 {
		j.throw(fmt.Sprint(t, " not inside a loop or switch"))
	}
	j.semicolon()
	return NewAtom(t, name)
}

func (j *jsparser) rember(name string) bool {
	for i := len(j.labels) - 1; i >= 0; i-- {
		if j.labels[i] == name {
			return true
		}
	}
	return false
}

func (j *jsparser) for_() IAst {
	j.expect("(")
	var init IAst
	if !j.is_token(scanner.TokenPunc, ";") {
		if j.is_token(scanner.TokenKeyword, "var") {
			j.next()
			init = j.var_(true)
		} else {
			init = j.expression(true, true)
		}
		if j.is_token(scanner.TokenOperator, "in") {
			if init.Type() == Type_Var && len(init.(*Var).Defs) > 1 {
				j.throw("Only one variable declaration allowed in for..in loop")
			}
			return j.for_in(init)
		}
	}
	return j.regular_for(init)
}

func (j *jsparser) regular_for(init IAst) IAst {
	j.expect(";")
	var test IAst
	var step IAst
	if !j.is_token(scanner.TokenPunc, ";") {
		test = j.expression(true, false)
	}
	j.expect(";")
	if !j.is_token(scanner.TokenPunc, ")") {
		step = j.expression(true, false)
	}
	j.expect(")")
	return NewFor(Type_For, init, test, step, j.loop(j.statement))
}

func (j *jsparser) for_in(init IAst) IAst {
	var lhs IAst
	if init.Type() == Type_Var {
		lhs = NewAtom(Type_Name, init.(*Var).Defs[0].Name())
	} else {
		lhs = init
	}
	j.next()
	obj := j.expression(true, false)
	j.expect(")")
	return NewFor(Type_For_In, init, lhs, obj, j.loop(j.statement))
}

func (j *jsparser) vardefs(no_in bool) []*VarDef {
	a := make([]*VarDef, 0)
	for {
		if !j.is_token(scanner.TokenName, "") {
			j.throw("")
		}
		name := j.token.Value
		j.next()
		if j.is_token(scanner.TokenOperator, "=") {
			j.next()
			a = append(a, NewDef(name, j.expression(false, no_in)))
		} else {
			a = append(a, NewDef(name, nil))
		}
		if !j.is_token(scanner.TokenPunc, ",") {
			break
		}
		j.next()
	}
	return a
}

func (j *jsparser) var_(no_in bool) IAst {
	return NewVar(j.vardefs(no_in))
}

func (j *jsparser) const_() IAst {
	return NewVar(j.vardefs(false))
}

func (j *jsparser) if_() IAst {
	cond := j.parenthesised()
	body := j.statement()
	var belse IAst
	if j.is_token(scanner.TokenKeyword, "else") {
		j.next()
		belse = j.statement()
	}
	return NewIf(cond, body, belse)
}

func (j *jsparser) parenthesised() IAst {
	j.expect("(")
	ex := j.expression(true, false)
	j.expect(")")
	return ex
}

func (j *jsparser) switch_block_() []*Case {
	/*j.in_func += 1
	loop := j.in_loop
	j.in_directives = true
	j.in_loop = 0
	a := j.switch_block_loop()
	j.in_func -= 1
	j.in_loop = loop
	return a*/
	j.in_loop += 1
	ret := j.switch_block_loop()
	j.in_loop -= 1
	return ret
}

func (j *jsparser) switch_block_loop() []*Case {
	j.expect("{")
	a := make([]*Case, 0)
	/*cur := make([]IAst, 0)
	for !j.is_token(scanner.TokenPunc, "}") {
		if !j.input.Eof() {
			j.throw("")
		}
		if j.is_token(scanner.TokenKeyword, "case") {
			j.next()
			a = append(a, NewCase(j.expression(true, false), cur))
			cur = make([]IAst, 0)
			j.expect(":")
		} else if j.is_token(scanner.TokenKeyword, "default") {
			j.next()
			j.expect(":")
			a = append(a, NewCase(nil, cur))
			cur = make([]IAst, 0)
		} else {
			if cur == nil {
				j.throw("")
			}
			cur = append(cur, j.statement())
		}
	}*/
	var cur *Case
	for !j.is_token(scanner.TokenPunc, "}") {
		if !j.input.Eof() {
			j.throw("")
		}
		if j.is_token(scanner.TokenKeyword, "case") {
			j.next()
			cur = NewCase(j.expression(true, false), make([]IAst, 0))
			a = append(a, cur)
			j.expect(":")
		} else if j.is_token(scanner.TokenKeyword, "default") {
			j.next()
			j.expect(":")
			cur = NewCase(nil, make([]IAst, 0))
			a = append(a, cur)
		} else {
			if cur == nil {
				j.throw("")
			}
			cur.Body = append(cur.Body, j.statement())
		}
	}
	j.next()
	return a
}

func (j *jsparser) try_() IAst {
	body := j.block_()
	var catch *Catch
	var finally []IAst
	if j.is_token(scanner.TokenKeyword, "catch") {
		j.next()
		j.expect("(")
		if !j.is_token(scanner.TokenName, "") {
			j.throw("Name expected")
		}
		name := j.token.Value
		j.next()
		j.expect(")")
		catch = NewCatch(name, j.block_())
	}
	if j.is_token(scanner.TokenKeyword, "finally") {
		j.next()
		finally = j.block_()
	}
	if catch == nil && finally == nil {
		j.throw("miss catch/finally blocks")
	}
	return NewTry(body, finally, catch)
}
