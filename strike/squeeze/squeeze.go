package squeeze

import (
	p "github.com/blinkat/blinks/strike/parser"
	"github.com/blinkat/blinks/strike/parser/adapter"
)

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

func boolean_expr(expr p.IAst) bool {
	name := expr.Name()
	return (expr.Type() == p.Type_Unary_Prefix) && member(unary_prefix_symbol, name) ||
		(expr.Type() == p.Type_Binnary && member(binary_symbol_1, name)) ||
		(expr.Type() == p.Type_Binnary && member(binary_symbol_2, name) &&
			boolean_expr(expr.(*p.Binary).Left) && boolean_expr(expr.(*p.Binary).Right)) ||
		(expr.Type() == p.Conditional &&
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
