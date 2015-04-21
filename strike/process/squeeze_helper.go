package process

import (
	"bytes"
	p "github.com/blinkat/blinks/strike/parser"
	"strings"
)

func boolean_expr(expr p.IAst) bool {
	name := expr.Name()
	return (expr.Type() == p.Type_Unary_Prefix) && member(unary_prefix_symbol, name) ||
		(expr.Type() == p.Type_Binnary && member(binary_symbol_1, name)) ||
		(expr.Type() == p.Type_Binnary && member(binary_symbol_2, name) &&
			boolean_expr(expr.(*p.Binary).Left) && boolean_expr(expr.(*p.Binary).Right)) ||
		(expr.Type() == p.Type_Conditional &&
			boolean_expr(expr.(*p.Conditional).True) && boolean_expr(expr.(*p.Conditional).False)) ||
		(expr.Type() == p.Type_Assign && expr.Name() == "true" && boolean_expr(expr.(*p.Assign).Right)) ||
		(expr.Type() == p.Type_Seq && boolean_expr(expr.(*p.Seq).Expr2))
}

func empty(b p.IAst) bool {
	if b != nil {
		if b.Type() == p.Type_Block {
			a := b.(*p.Block)
			if a.Statements != nil {
				return len(a.Statements) == 0
			} else {
				return true
			}
		}
	}
	return false
}

func is_string(ast p.IAst) bool {
	return ast.Type() == p.Type_String ||
		ast.Type() == p.Type_Unary_Prefix && ast.Name() == "typeof" ||
		ast.Type() == p.Type_Binnary && ast.Name() == "+" &&
			(is_string(ast.(*p.Binary).Left) || is_string(ast.(*p.Binary).Right))
}

func last_stat(a p.IAst) p.IAst {
	if a.Type() == p.Type_Block {
		e := a.(*p.Block)
		if e.Statements != nil && len(e.Statements) > 0 {
			return e.Statements[len(e.Statements)-1]
		}
	}
	return a
}

func aborts(ast p.IAst) bool {
	if ast != nil {
		switch last_stat(ast).Type() {
		case p.Type_Return, p.Type_Break, p.Type_Coutinue, p.Type_Thorw:
			return true
		}
	}
	return false
}

func best_of(a1, a2 p.IAst, w *Walker) p.IAst {
	ret1 := GenCode(a1, w)
	var ret2 string
	if a2.Type() == p.Type_Stat {
		ret2 = GenCode(a2.(*p.Stat).Statement, w)
	} else {
		ret2 = GenCode(a2, w)
	}

	if len(ret1) > len(ret2) {
		return a2
	} else {
		return a1
	}
}

func make_string(str string) string {
	dq := 0
	sq := 0
	rs := []rune(str)

	var rbuf bytes.Buffer
	for _, r := range rs {
		switch r {
		case '\\':
			rbuf.WriteString("\\\\")
			break
		case '\b':
			rbuf.WriteString("\\b")
			break
		case '\f':
			rbuf.WriteString("\\f")
			break
		case '\n':
			rbuf.WriteString("\\n")
			break
		case '\r':
			rbuf.WriteString("\\r")
			break
		case '\u2028':
			rbuf.WriteString("\\u2028")
			break
		case '\u2029':
			rbuf.WriteString("\\u2029")
			break
		case '"':
			dq += 1
			rbuf.WriteString("\\\"")
			break
		case '\'':
			sq += 1
			rbuf.WriteString("'")
			break
		case 0:
			rbuf.WriteString("\\0")
			break
		default:
			rbuf.WriteRune(r)
			break
		}
	}

	ret := rbuf.String()
	if dq > sq {
		ret = strings.Replace(ret, "\x27", "\\'", 0)
		return "'" + ret + "'"
	} else {
		return "\"" + strings.Replace(ret, "\x22", "\\\"", 0) + "\""
	}
}
