package parser

import (
	"fmt"
	"github.com/blinkat/blinks/strike/parser/adapter"
	"github.com/blinkat/blinks/strike/parser/scanner"
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

func (j *jsparser) Parse() *Ast {
	ret := ast(Ast_TopLevel, "")
	for j.input.Eof() {
		ret.Attributes = append(ret.Attributes, j.statement())
	}
	return ret
}

func (j *jsparser) is_token(t int, value string) bool {
	return j.is(j.token, t, value)
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
	fmt.Println("message :", j.count, j.token)
	j.count += 1
	return j.token
}

func (j *jsparser) throw(msg string) {
	if msg == "" {
		msg = "Unexpected token"
	}
	s := fmt.Sprint(msg, "\ntoken:", j.token)
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
	return j.token.Nlb || !j.input.Eof() || j.is_token(scanner.TokenPunc, "}")
}

func (j *jsparser) labeled_statement(label string) *Ast {
	j.labels = append(j.labels, label)
	stat := j.statement()
	j.labels = j.labels[:len(j.labels)-1]
	return ast(Ast_Label, label, stat)
}

//-------------[ statement ]---------------
func (j *jsparser) statement() *Ast {
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

func (j *jsparser) read_string() *Ast {
	dir := j.in_directives
	stat := j.simple_statement()
	if dir && stat.Attributes[0].Type == scanner.TokenString && !j.is_token(scanner.TokenPunc, ",") {
		return ast(Ast_Directive, stat.Attributes[0].Name)
	}
	return stat
}

func (j *jsparser) read_punc() *Ast {
	switch j.token.Value {
	case "{":
		return ast(Ast_Block, "", j.block_())
	case "[", "(":
		return j.simple_statement()
	case ";":
		j.next()
		return ast(Ast_None, "block")
	default:
		j.throw("")
		return nil
	}
}

func (j *jsparser) read_keyword() *Ast {
	val := j.token.Value
	j.next()
	switch val {
	case "break":
		return j.break_cont(Ast_Break)
	case "continue":
		return j.break_cont(Ast_Coutinue)
	case "debugger":
		j.semicolon()
		return ast(Ast_None, "debugger")
	case "do":
		body := j.loop(j.statement)
		j.expect_token(scanner.TokenKeyword, "while")
		cond := j.parenthesised()
		j.semicolon()
		return ast(Ast_Do, "", cond, body)
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
			return ast(Ast_Return, "")
		} else if j.can_insert_semicolon() {
			return ast(Ast_Return, "")
		} else {
			a := ast(Ast_Return, "", j.expression(true, false))
			j.semicolon()
			return a
		}
	case "switch":
		return ast(Ast_Switch, "", j.parenthesised(), j.switch_block_())
	case "throw":
		if j.token.Nlb {
			j.throw("Illegal newline after 'throw'")
		}
		ret := j.expression(true, false)
		j.semicolon()
		return ast(Ast_Thorw, "", ret)

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
		return ast(Ast_While, "", j.parenthesised(), j.loop(j.statement))
	case "with":
		return ast(Ast_With, "", j.parenthesised(), j.statement())
	default:
		j.throw("")
		return nil
	}
}

func (j *jsparser) read_name() *Ast {
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
func (j *jsparser) expr_atom(allow_calls bool) *Ast {
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
		atom := ast(TokenTypeToAstType(j.token.Type), j.token.Value)
		if atom.Type == Ast_Regexp {
			atom.Attributes = append(atom.Attributes, ast(Ast_Regexp_Mode, j.token.Attributes[0]))
		}
		j.next()
		return j.subscripts(atom, allow_calls)
	}
	j.throw("")
	return nil
}

func (j *jsparser) subscripts(expr *Ast, allow_calls bool) *Ast {
	if j.is_token(scanner.TokenPunc, ".") {
		j.next()
		return j.subscripts(ast(Ast_Dot, j.as_name(), expr), allow_calls)
	}
	if j.is_token(scanner.TokenPunc, "[") {
		j.next()
		ret := j.expression(true, false)
		j.expect("]")
		return j.subscripts(ast(Ast_Sub, "", expr, ret), allow_calls)
	}
	if allow_calls && j.is_token(scanner.TokenPunc, "(") {
		j.next()
		return j.subscripts(ast(Ast_Call, "", expr, j.expr_list(")", false, false)), true)
	}
	return expr
}

func (j *jsparser) maybe_unary(allow_calls bool) *Ast {
	if j.is_token(scanner.TokenOperator, "") && adapter.UnaryPrefix(j.token.Value) {
		ret := j.token.Value
		j.next()
		return ast(Ast_Unary_Prefix, ret, j.maybe_unary(allow_calls))
	}
	val := j.expr_atom(allow_calls)
	for j.is_token(scanner.TokenOperator, "") && adapter.UnaryPostfix(j.token.Value) && !j.token.Nlb {
		val = ast(Ast_Unary_Postfix, j.token.Value, val)
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

func (j *jsparser) expression(commas, no_in bool) *Ast {
	expr := j.maybe_assign(no_in)
	if commas && j.is_token(scanner.TokenPunc, ",") {
		j.next()
		return ast(Ast_Seq, "", expr, j.expression(true, no_in))
	}
	return expr
}

func (j *jsparser) expr_list(end string, allow_trailing_comma, allow_empty bool) *Ast {
	first := true
	a := make([]*Ast, 0)
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
			a = append(a, ast(Ast_Atom, "undefined"))
		} else {
			a = append(a, j.expression(false, false))
		}
	}
	j.next()
	return ast(Ast_None, "", a...)
}

func (j *jsparser) maybe_assign(no_in bool) *Ast {
	left := j.maybe_conditional(no_in)
	val := j.token.Value
	ass := adapter.Assignment(val)
	if j.is_token(scanner.TokenOperator, "") && ass != "" {
		j.next()
		return ast(Ast_Assign, ass, left, j.maybe_assign(no_in))
	}
	return left
}

func (j *jsparser) expr_op(left *Ast, min_prec int, no_in bool) *Ast {
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
		return j.expr_op(ast(Ast_Binnary, op, left, right), min_prec, no_in)
	}
	return left
}

func (j *jsparser) maybe_conditional(no_in bool) *Ast {
	expr := j.expr_op(j.maybe_unary(true), 0, no_in)
	if j.is_token(scanner.TokenOperator, "?") {
		j.next()
		yes := j.expression(false, false)
		j.expect(":")
		return ast(Ast_Conditional, "", expr, yes, j.expression(false, no_in))
	}
	return expr
}

//loop
type loop_fn func() *Ast

func (j *jsparser) loop(fn loop_fn) *Ast {
	j.in_loop += 1
	ret := fn()
	j.in_loop -= 1
	return ret
}

func (j *jsparser) simple_statement() *Ast {
	ret := j.expression(true, false)
	j.semicolon()
	return ast(Ast_Stat, "", ret)
}

//----------------[ types ]-----------------
func (j *jsparser) new_() *Ast {
	newxp := j.expr_atom(false)
	args := ast(Ast_None, "")
	if j.is_token(scanner.TokenPunc, "(") {
		j.next()
		args = j.expr_list(")", false, false)
	}
	return j.subscripts(ast(Ast_New, "", newxp, args), true)
}

func (j *jsparser) array_() *Ast {
	return ast(Ast_Array, "", j.expr_list("]", true, true))
}

func (j *jsparser) object_() *Ast {
	first := true
	a := make([]*Ast, 0)
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
			a = append(a, ast(Ast_None, name, j.function_(false)))
		} else {
			j.expect(":")
			a = append(a, ast(Ast_None, name, j.expression(false, false)))
		}
	}
	j.next()
	return ast(Ast_Object, "", a...)
}

func (j *jsparser) block_() *Ast {
	j.expect("{")
	a := make([]*Ast, 0)
	for !j.is_token(scanner.TokenPunc, "}") {
		if !j.input.Eof() {
			j.throw("")
		}
		a = append(a, j.statement())
	}
	j.next()
	return ast(Ast_Block, "", a...)
}

func (j *jsparser) function_(in_statement bool) *Ast {
	name := ""
	if j.is_token(scanner.TokenName, "") {
		name = j.token.Value
		j.next()
	}

	if in_statement && name == "" {
		j.throw("")
	}
	j.expect("(")
	ty := Ast_Func
	if in_statement {
		ty = Ast_Defunc
	}
	return ast(ty, name, j.function_params(), j.function_block())
}

func (j *jsparser) function_params() *Ast {
	first := true
	a := make([]*Ast, 0)
	for !j.is_token(scanner.TokenPunc, ")") {
		if first {
			first = false
		} else {
			j.expect(",")
		}
		if !j.is_token(scanner.TokenName, "") {
			j.throw("")
		}
		a = append(a, ast(Ast_Func_Params, j.token.Value))
		j.next()
	}
	j.next()
	return ast(Ast_None, "", a...)
}

func (j *jsparser) function_block() *Ast {
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

func (j *jsparser) break_cont(t int) *Ast {
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
	return ast(t, name)
}

func (j *jsparser) rember(name string) bool {
	for i := len(j.labels) - 1; i >= 0; i-- {
		if j.labels[i] == name {
			return true
		}
	}
	return false
}

func (j *jsparser) for_() *Ast {
	j.expect("(")
	var init *Ast
	if !j.is_token(scanner.TokenPunc, ";") {
		if j.is_token(scanner.TokenKeyword, "var") {
			j.next()
			init = j.var_(true)
		} else {
			init = j.expression(true, true)
		}
		if j.is_token(scanner.TokenOperator, "in") {
			if init.Type == Ast_Var && len(init.Attributes) > 1 {
				j.throw("Only one variable declaration allowed in for..in loop")
			}
			return j.for_in(init)
		}
	}
	return j.regular_for(init)
}

func (j *jsparser) regular_for(init *Ast) *Ast {
	j.expect(";")
	var test *Ast
	var step *Ast
	if !j.is_token(scanner.TokenPunc, ";") {
		test = j.expression(true, false)
	}
	j.expect(";")
	if !j.is_token(scanner.TokenPunc, ")") {
		step = j.expression(true, false)
	}
	j.expect(")")
	return ast(Ast_For, "", init, test, step, j.loop(j.statement))
}

func (j *jsparser) for_in(init *Ast) *Ast {
	var lhs *Ast
	if init.Type == Ast_Var {
		lhs = ast(Ast_Name, init.Attributes[0].Name)
	} else {
		lhs = init
	}
	j.next()
	obj := j.expression(true, false)
	j.expect(")")
	return ast(Ast_For_In, "", init, lhs, obj, j.loop(j.statement))
}

func (j *jsparser) vardefs(no_in bool) []*Ast {
	a := make([]*Ast, 0)
	for {
		if !j.is_token(scanner.TokenName, "") {
			j.throw("")
		}
		name := j.token.Value
		j.next()
		if j.is_token(scanner.TokenOperator, "=") {
			j.next()
			a = append(a, ast(Ast_None, name, j.expression(false, no_in)))
		} else {
			a = append(a, ast(Ast_None, name))
		}
		if !j.is_token(scanner.TokenPunc, ",") {
			break
		}
		j.next()
	}
	return a
}

func (j *jsparser) var_(no_in bool) *Ast {
	return ast(Ast_Var, "", j.vardefs(no_in)...)
}

func (j *jsparser) const_() *Ast {
	return ast(Ast_Var, "", j.vardefs(false)...)
}

func (j *jsparser) if_() *Ast {
	cond := j.parenthesised()
	body := j.statement()
	var belse *Ast
	if j.is_token(scanner.TokenKeyword, "else") {
		j.next()
		belse = j.statement()
	}
	return ast(Ast_If, "", cond, body, belse)
}

func (j *jsparser) parenthesised() *Ast {
	j.expect("(")
	ex := j.expression(true, false)
	j.expect(")")
	return ex
}

func (j *jsparser) switch_block_() *Ast {
	return j.loop(j.switch_block_loop)
}

func (j *jsparser) switch_block_loop() *Ast {
	j.expect("{")
	a := make([]*Ast, 0)
	var cur *Ast
	for !j.is_token(scanner.TokenPunc, "}") {
		if !j.input.Eof() {
			j.throw("")
		}
		if j.is_token(scanner.TokenKeyword, "case") {
			j.next()
			cur = ast(Ast_None, "")
			j.expect(":")
			a = append(a, ast(Ast_None, "", j.expression(true, false), cur))
		} else if j.is_token(scanner.TokenKeyword, "default") {
			j.next()
			j.expect(":")
			cur = ast(Ast_None, "")
			a = append(a, ast(Ast_None, "", nil, cur))
		} else {
			if cur == nil {
				j.throw("")
			}
			cur.Attributes = append(cur.Attributes, j.statement())
		}
	}
	j.next()
	return ast(Ast_None, "", a...)
}

func (j *jsparser) try_() *Ast {
	body := j.block_()
	var catch, finally *Ast
	if j.is_token(scanner.TokenKeyword, "catch") {
		j.next()
		j.expect("(")
		if j.is_token(scanner.TokenName, "") {
			j.throw("Name expected")
		}
		name := j.token.Value
		j.next()
		j.expect(")")
		catch = ast(Ast_None, name, j.block_())
	}
	if j.is_token(scanner.TokenKeyword, "finally") {
		j.next()
		finally = j.block_()
	}
	if catch == nil && finally == nil {
		j.throw("miss catch/finally blocks")
	}
	return ast(Ast_Try, "", body, catch, finally)
}
