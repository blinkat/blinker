package process

import (
	/*
		int left_offset(int value, int pos) {
			return value << pos;
		}

		int right_offset(int value, int pos) {
			return value >> pos;
		}
	*/
	"C"
	"fmt"
	p "github.com/blinkat/blinks/strike/parser"
	"github.com/blinkat/blinks/strike/parser/adapter"
	"math"
	"strconv"
)

const eval_err = "evaluate error"
const (
	eval_type_int = iota
	eval_type_string
	eval_type_float
	eval_type_bool
)

//--------------[ eval ast ]------------
type eval_ast struct {
	value  interface{}
	etype  int
	source interface{}
	is_num bool
}

func new_eval_ast(value interface{}, t int, need_num bool) *eval_ast {
	e := &eval_ast{}
	e.source = value
	e.etype = t
	e.is_num = false

	if need_num {
		if t == eval_type_string {
			val, err := parse_number(value.(string))
			if err != nil {
				e.value = 0
			} else {
				e.is_num = true
				e.value = val
			}
		} else if t == eval_type_bool {
			if value.(bool) {
				e.is_num = true
				e.value = 1
			} else {
				e.is_num = true
				e.value = 0
			}
		}
	} else {
		e.value = e.source
	}

	return e
}

//------------[ evaluate ]------------

func evaluate(ast p.IAst) (interface{}, error) {
	switch ast.Type() {
	case p.Type_String, p.Type_Number:
		return ast.Name(), nil
	case p.Type_Name, p.Type_Atom:
		return eval_atom(ast)
	case p.Type_Unary_Prefix:
		u := ast.(*p.Unary)
		switch ast.Name() {
		case "!":
			return eval_negation(u)
		case "~":
			return eval_non(u)
		case "-":
			return eval_minus(u)
		case "+":
			return eval_plus(u)
		}
	case p.Type_Binnary:
		b := ast.(*p.Binary)
		switch b.Name() {
		case "&&":
			return eval_and(b)
		case "||":
			return eval_or(b)
		case "|":
			return eval_single_or(b)
		case "&":
			return eval_single_and(b)
		case "^":
			return eval_single_non(b)
		case "+":
			return eval_add(b)
		case "-":
			return eval_sub(b)
		case "*":
			return eval_ride(b)
		case "/":
			return eval_divide(b)
		case "%":
			return eval_remainder(b)
		case "<<":
			return eval_left_offset(b)
		case ">>":
			return eval_right_offset(b)
		case ">>>":
			return eval_no_symbol_offset(b)
		case "==":
			return eval_equal(b)
		case "!=":
			return eval_non_equal(b)
		case ">":
			return eval_greater(b)
		case ">=":
			return eval_greater_equal(b)
		case "<":
			return eval_less(b)
		case "<=":
			return eval_less_equal(b)
		}
	}

	return nil, fmt.Errorf(eval_err)
}

func eval_atom(ast p.IAst) (interface{}, error) {
	switch ast.Name() {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	}
	return nil, fmt.Errorf(eval_err)
}

//--------------[ unary ]--------------
// !
func eval_negation(ast *p.Unary) (interface{}, error) {
	e, err := next_eval(ast.Expr, true)
	if err == nil {
		if e.etype == eval_type_string {
			return len(e.source.(string)) == 0, nil
		} else if e.etype == eval_type_int {
			return !(e.value.(int) == 0), nil
		} else if e.etype == eval_type_float {
			return !(e.value.(float64) == 0), nil
		} else if e.etype == eval_type_bool {
			return !e.value.(bool), nil
		}
	}
	return nil, err
}

// ~
func eval_non(ast *p.Unary) (interface{}, error) {
	e, err := next_eval(ast.Expr, true)
	if err == nil {
		switch e.value.(type) {
		case int:
			return ^e.value.(int), nil
		case float64:
			return ^int(e.value.(float64)), nil
		}
	}
	return nil, fmt.Errorf("error")
}

// -
func eval_minus(ast *p.Unary) (interface{}, error) {
	e, err := next_eval(ast.Expr, true)
	if err == nil {
		switch e.value.(type) {
		case int:
			return -e.value.(int), nil
		case float64:
			return -e.value.(float64), nil
		}
	}
	return nil, fmt.Errorf("error")
}

// +
func eval_plus(ast *p.Unary) (interface{}, error) {
	e, err := next_eval(ast.Expr, true)
	if err == nil {
		switch e.value.(type) {
		case int:
			return +e.value.(int), nil
		case float64:
			return +e.value.(float64), nil
		}
	}
	return nil, fmt.Errorf("error")
}

//------------[ binary ]----------------

// &&
func eval_and(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		b1 := get_bool(left, ast.Left.Type())
		b2 := get_bool(right, ast.Left.Type())

		if b1 && b2 {
			return right.source, nil
		} else if !b1 {
			return left.source, nil
		} else {
			return right.source, nil
		}
	}
	return nil, err
}

// ||
func eval_or(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		if get_bool(left, ast.Left.Type()) {
			return left.source, nil
		} else {
			return right.source, nil
		}
	}
	return nil, err
}

// |
func eval_single_or(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return get_int(left) | get_int(right), nil
	}
	return nil, err
}

// &
func eval_single_and(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return get_int(left) & get_int(right), nil
	}
	return nil, err
}

// ^
func eval_single_non(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return get_int(left) ^ get_int(right), nil
	}
	return nil, err
}

// +
func eval_add(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return get_float(left) + get_float(right), nil
	}
	return nil, err
}

// -
func eval_sub(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return get_float(left) - get_float(right), nil
	}
	return nil, err
}

// *
func eval_ride(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return get_float(left) * get_float(right), nil
	}
	return nil, err
}

// /
func eval_divide(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return get_float(left) / get_float(right), nil
	}
	return nil, err
}

// %
func eval_remainder(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return math.Mod(get_float(left), get_float(right)), nil
	}
	return nil, err
}

// <<
func eval_left_offset(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return int(C.left_offset(get_c_int(left), get_c_int(right))), nil
	}
	return nil, err
}

// >>
func eval_right_offset(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return int(C.right_offset(get_c_int(left), get_c_int(right))), nil
	}
	return nil, err
}

// >>>
func eval_no_symbol_offset(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		return unsigned_offset(int32(get_int(left)), int32(get_int(right))), nil
	}
	return nil, err
}

// ==
func eval_equal(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		if left.is_num && right.is_num {
			return get_float(left) == get_float(right), nil
		} else if left.etype == eval_type_string && right.etype == eval_type_string {
			return left.source.(string) == right.source.(string), nil
		} else {
			return false, nil
		}
	}
	return nil, err
}

// !=
func eval_non_equal(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		if left.is_num && right.is_num {
			return get_float(left) != get_float(right), nil
		} else if left.etype == eval_type_string && right.etype == eval_type_string {
			return left.source.(string) != right.source.(string), nil
		} else {
			return false, nil
		}
	}
	return nil, err
}

// >
func eval_greater(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		if left.is_num && right.is_num {
			return get_float(left) > get_float(right), nil
		} else if left.etype == eval_type_string && right.etype == eval_type_string {
			return left.source.(string) > right.source.(string), nil
		} else {
			return false, nil
		}
	}
	return nil, err
}

// <
func eval_less(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		if left.is_num && right.is_num {
			return get_float(left) < get_float(right), nil
		} else if left.etype == eval_type_string && right.etype == eval_type_string {
			return left.source.(string) < right.source.(string), nil
		} else {
			return false, nil
		}
	}
	return nil, err
}

// >=
func eval_greater_equal(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		if left.is_num && right.is_num {
			return get_float(left) >= get_float(right), nil
		} else if left.etype == eval_type_string && right.etype == eval_type_string {
			return left.source.(string) >= right.source.(string), nil
		} else {
			return false, nil
		}
	}
	return nil, err
}

// <=
func eval_less_equal(ast *p.Binary) (interface{}, error) {
	left, right, err := get_left_right(ast)
	if err == nil {
		if left.is_num && right.is_num {
			return get_float(left) <= get_float(right), nil
		} else if left.etype == eval_type_string && right.etype == eval_type_string {
			return left.source.(string) <= right.source.(string), nil
		} else {
			return false, nil
		}
	}
	return nil, err
}

//------------[ helper ]----------------------
const int_32_max = 0x7FFFFFFF

func get_c_int(a *eval_ast) C.int {
	return C.int(get_int(a))
}

func get_int(a *eval_ast) int {
	switch a.value.(type) {
	case float64:
		return int(a.value.(float64))
	case int:
		return a.value.(int)
	}
	return 0
}

func get_float(a *eval_ast) float64 {
	switch a.value.(type) {
	case int:
		return float64(a.value.(int))
	case float64:
		return a.value.(float64)
	}
	return 0
}

func get_bool(a *eval_ast, t int) bool {
	if t == p.Type_String {
		return len(a.source.(string)) != 0
	} else {
		switch a.value.(type) {
		case float64:
			return a.value.(float64) != 0
		case int:
			return a.value.(int) != 0
		case bool:
			return a.value.(bool)
		}
	}
	return false
}

func get_left_right(ast *p.Binary) (*eval_ast, *eval_ast, error) {
	left, err1 := next_eval(ast.Left, true)
	right, err2 := next_eval(ast.Right, true)

	if err1 != nil || err2 != nil {
		return nil, nil, fmt.Errorf("error")
	}
	return left, right, nil
}

func next_eval(ast p.IAst, need_num bool) (*eval_ast, error) {
	val, err := evaluate(ast)
	if err != nil {
		return nil, err
	}

	t := 0
	switch val.(type) {
	case bool:
		t = eval_type_bool
		break
	case float64:
		t = eval_type_float
		break
	case int:
		t = eval_type_int
		break
	case string:
		t = eval_type_string
		break
	default:
		return nil, fmt.Errorf("none type", val)
	}

	return new_eval_ast(val, t, need_num), nil
}

func unsigned_offset(value, pos int32) int64 {
	ret := C.int(value)
	if pos != 0 {
		ret = C.right_offset(ret, 1)
		ret &= int_32_max
		ret = C.right_offset(ret, C.int(pos-1))
	}
	return int64(ret)
}

func parse_number(val string) (interface{}, error) {
	if adapter.IsHexNumber(val) {
		return strconv.ParseInt(val, 16, 64)
	} else if adapter.IsOctNumber(val) {
		return strconv.ParseInt(val, 8, 64)
	} else if adapter.IsDecNumber(val) {
		return strconv.ParseFloat(val, 64)
	} else {
		return nil, fmt.Errorf(eval_err)
	}
}

//---------------[ when constant ]-------------
type WhenCall func(p.IAst, p.IAst, interface{}) p.IAst
type WhenNoCall func(p.IAst) p.IAst

func WhenConstant(expr p.IAst, yes WhenCall, no WhenNoCall) p.IAst {
	val, err := evaluate(expr)
	var ast p.IAst
	if err == nil {
		switch val.(type) {
		case string:
			ast = p.NewString(val.(string))
			break
		case float64, int:
			ast = p.NewNumber(fmt.Sprint(val))
			break
		case bool:
			ast = p.NewAtom(p.Type_Name, fmt.Sprint(val))
			break
		default:
			if val == nil {
				ast = p.NewAtom(p.Type_Atom, "null")
				break
			}
			panic(fmt.Sprint("Can't handle constant of type:", val))
		}
		return yes(expr, ast, val)
	} else {
		if expr.Type() == p.Type_Binnary {
			e := expr.(*p.Binary)
			if e.Name() == "===" || e.Name() == "!==" &&
				(is_string(e.Left) && is_string(e.Right)) || (boolean_expr(e.Left) && boolean_expr(e.Right)) {
				rs := []rune(e.Name())
				e.SetName(string(rs[:2]))
			} else if no != nil && expr.Type() == p.Type_Binnary && (e.Name() == "||" || e.Name() == "&&") {
				lval, err2 := evaluate(e.Left)

				if err2 != nil {
					if e.Name() == "&&" {
						if lval != nil {
							expr = e.Right
						} else {
							expr = p.NewAtom(p.Type_Reserve, fmt.Sprint(lval))
						}
					} else if e.Name() == "||" {
						if lval != nil {
							expr = p.NewAtom(p.Type_Reserve, fmt.Sprint(lval))
						} else {
							expr = e.Right
						}
					}
				}
			}

			if no != nil {
				return no(expr)
			}
		}
	}
	return nil
}
