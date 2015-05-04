package parser

import (
	"github.com/blinkat/blinker/strike/js/parser/scanner"
)

const (
	Type_None     = -1
	Type_TopLevel = iota
	Type_Punc
	Type_Dot
	Type_Sub
	Type_Call
	Type_Atom
	Type_New
	Type_Array
	Type_Regexp
	Type_Name
	Type_String
	Type_Number
	Type_Regexp_Mode
	Type_Unary_Prefix
	Type_Unary_Postfix
	Type_Binnary
	Type_Conditional
	Type_Assign
	Type_Seq
	Type_Object
	Type_Defunc
	Type_Func
	Type_Block
	Type_Func_Params
	Type_Label
	Type_Stat
	Type_Var
	Type_For
	Type_For_In
	Type_If
	Type_Try
	Type_Directive
	Type_Do
	Type_Return
	Type_Switch
	Type_Thorw
	Type_While
	Type_With
	Type_Break
	Type_Coutinue
	Type_Debugger
	Type_Property
	Type_Const
	Type_Lambda
	Type_Arg
	Type_Catch
	Type_Func_Body
	Type_Reserve
)

var ast_type_strings []string

func init() {
	ast_type_strings = make([]string, 0)
	ast_type_strings = append(ast_type_strings,
		"None",
		"TopLevel",
		"Punc",
		"Dot",
		"Sub",
		"Call",
		"Atom",
		"New",
		"Array",
		"Regexp",
		"Name",
		"String",
		"Number",
		"Regexp_Mode",
		"Unary_Prefix",
		"Unary_Postfix",
		"Binnary",
		"Conditional",
		"Assign",
		"Seq",
		"Object",
		"Defunc",
		"Func",
		"Block",
		"Func_Params",
		"Label",
		"Stat",
		"Var",
		"For",
		"For_In",
		"If",
		"Try",
		"Directive",
		"Do",
		"Return",
		"Switch",
		"Thorw",
		"While",
		"With",
		"Break",
		"Coutinue",
		"Debugger",
		"Property",
		"Const",
		"Lambda",
		"Arg",
		"Catch",
		"Func_Body",
		"Reserve",
	)
}

//------------------[ ast ]---------------------

type IAst interface {
	Type() int
	SetType(i int)
	Name() string
	SetName(n string)

	AtTop() bool
	SetAtTop(b bool)
	Splice() bool
	SetSplice(b bool)

	AtValue() IAst
	SetAtValue(v IAst)

	Scope() *AstScope
	SetScope(v *AstScope)

	TypeName() string
}

func TokenTypeToAstType(t int) int {
	switch t {
	case scanner.TokenName:
		return Type_Name
	case scanner.TokenString:
		return Type_String
	case scanner.TokenNumber:
		return Type_Number
	case scanner.TokenRegexp:
		return Type_Regexp
	case scanner.TokenAtom:
		return Type_Atom

	default:
		return Type_None
	}
}

func GetTypeName(t int) string {
	if t >= -1 && t < len(ast_type_strings)-1 {
		return ast_type_strings[t]
	}
	return ""
}

//----------[ asts ]-----------
type ast struct {
	t      int
	name   string
	at_top bool
	splice bool
	at_val IAst
	scope  *AstScope
}

func (a *ast) Type() int {
	return a.t
}

func (a *ast) SetType(t int) {
	a.t = t
}

func (a *ast) Name() string {
	return a.name
}

func (a *ast) SetName(n string) {
	a.name = n
}

func (a *ast) AtTop() bool {
	return a.at_top
}

func (a *ast) SetAtTop(b bool) {
	a.at_top = b
}

func (a *ast) Splice() bool {
	return a.splice
}

func (a *ast) SetSplice(b bool) {
	a.splice = b
}

func (a *ast) AtValue() IAst {
	return a.at_val
}

func (a *ast) SetAtValue(v IAst) {
	a.at_val = v
}

func (a *ast) Scope() *AstScope {
	return a.scope
}
func (a *ast) SetScope(v *AstScope) {
	a.scope = v
}

func (a *ast) TypeName() string {
	return GetTypeName(a.t)
}

//----------[ children ]------------
type Toplevel struct {
	ast
	Statements []IAst
}

type Label struct {
	ast
	Stat IAst
}

type Directive struct {
	ast
}

type Block struct {
	ast
	Statements []IAst
}

type Debugger struct {
	ast
}

type Dot struct {
	ast
	Expr IAst
}

type Sub struct {
	ast
	Expr IAst
	Ret  IAst
}

type Call struct {
	ast
	Expr IAst
	List []IAst
}

type Regexp struct {
	ast
	Mode string
}

type Unary struct {
	ast
	Expr IAst
}

type Seq struct {
	ast
	Expr1 IAst
	Expr2 IAst
}

type Assign struct {
	ast
	Left  IAst
	Right IAst
}

type Conditional struct {
	ast
	Expr  IAst
	True  IAst
	False IAst
}

type Binary struct {
	ast
	Left  IAst
	Right IAst
}

type Stat struct {
	ast
	Statement IAst
}

type New struct {
	ast
	Expr IAst
	Args []IAst
}

type Array struct {
	ast
	List []IAst
}

type Property struct {
	ast
	Expr IAst
	Oper string
}

type Object struct {
	ast
	Propertys []*Property
}

type Function struct {
	ast
	Args []IAst
	Body *FuncBody
}

type For struct {
	ast
	Init IAst
	Cond IAst
	Step IAst
	Body IAst
}

type Var struct {
	ast
	Defs []*VarDef
}

type VarDef struct {
	ast
	Expr IAst
}

type If struct {
	ast
	Cond IAst
	Body IAst
	Else IAst
}

type Switch struct {
	ast
	Expr  IAst
	Cases []*Case
}

type Case struct {
	ast
	Expr IAst
	Body []IAst
}

type Try struct {
	ast
	Body    []IAst
	Catchs  *Catch
	Finally []IAst
}

type Catch struct {
	ast
	Body []IAst
}

type Do struct {
	ast
	Cond IAst
	Body IAst
}

type Return struct {
	ast
	Expr IAst
}

type While struct {
	ast
	Expr IAst
	Body IAst
}

type FuncBody struct {
	ast
	Exprs []IAst
}
