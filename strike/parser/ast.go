package parser

import (
	"github.com/blinkat/blinks/strike/parser/scanner"
)

const (
	Ast_None     = -1
	Ast_TopLevel = iota
	Ast_Punc
	Ast_Dot
	Ast_Sub
	Ast_Call
	Ast_Atom
	Ast_New
	Ast_Array
	Ast_Regexp
	Ast_Name
	Ast_String
	Ast_Number
	Ast_Regexp_Mode
	Ast_Unary_Prefix
	Ast_Unary_Postfix
	Ast_Binnary
	Ast_Conditional
	Ast_Assign
	Ast_Seq
	Ast_Object
	Ast_Defunc
	Ast_Func
	Ast_Block
	Ast_Func_Params
	Ast_Label
	Ast_Stat
	Ast_Var
	Ast_For
	Ast_For_In
	Ast_If
	Ast_Try
	Ast_Directive
	Ast_Do
	Ast_Return
	Ast_Switch
	Ast_Thorw
	Ast_While
	Ast_With
	Ast_Break
	Ast_Coutinue
)

type Ast struct {
	Type       int
	Name       string
	Attributes []*Ast
	AtTop      bool
	Splice     bool
	AtValue    *Ast
}

func ast(t int, name string, attr ...*Ast) *Ast {
	a := &Ast{}
	a.Type = t
	a.Name = name
	a.AtTop = false
	a.Splice = false
	a.AtValue = nil
	a.Attributes = append(make([]*Ast, 0), attr...)

	return a
}

func TokenTypeToAstType(t int) int {
	switch t {
	case scanner.TokenName:
		return Ast_Name
	case scanner.TokenString:
		return Ast_String
	case scanner.TokenNumber:
		return Ast_Number
	case scanner.TokenRegexp:
		return Ast_Regexp
	case scanner.TokenAtom:
		return Ast_Atom

	default:
		return Ast_None
	}
}

func NewAst(t int, name string, attr ...*Ast) *Ast {
	return ast(t, name, attr...)
}
