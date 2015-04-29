package parser

import (
	"github.com/blinkat/blinks/strike/js/parser/adapter"
	"math"
)

type AstScope struct {
	Names      map[string]int       //作用于符号表
	Mangled    map[string]string    //混淆后变量表
	RevManged  map[string]string    //混淆前变量表
	CName      int                  //已混淆变量个数
	Refs       map[string]*AstScope //引用的变量名
	UsesWith   bool
	UsesEval   bool
	Directives []string
	Parent     *AstScope
	Children   []*AstScope //子作用域
	Level      int
	Labels     *AstScope
	Body       []IAst
	IsTrue     bool
}

func (a *AstScope) Has(name string) *AstScope {
	for par := a; par != nil; par = par.Parent {
		if _, ok := par.Names[name]; ok {
			return par
		}
	}
	return nil
}

func (a *AstScope) HasMangled(name string) *AstScope {
	for par := a; par != nil; par = par.Parent {
		if _, ok := par.RevManged[name]; ok {
			return par
		}
	}
	return nil
}

func (a *AstScope) NextMangled() string {
	for {
		a.CName += 1
		name := base54(a.CName)
		prior := a.HasMangled(name)
		if prior != nil && a.Refs[prior.RevManged[name]] == prior {
			continue
		}
		prior = a.Has(name)
		if prior != nil && prior != a && a.Refs[name] == prior && prior.HasMangled(name) == nil {
			continue
		}
		if v, ok := a.Refs[name]; ok && v == nil {
			continue
		}
		if !adapter.IsIdentifier(name) {
			continue
		}
		return name
	}
}

func (a *AstScope) SetMangled(name, m string) string {
	a.RevManged[m] = name
	a.Mangled[name] = m
	return m
}

func (a *AstScope) GetMangled(name string) string {
	if a.UsesEval || a.UsesWith {
		return name
	}
	s := a.Has(name)
	if s == nil {
		return name
	}
	if _, ok := s.Mangled[name]; ok {
		return s.Mangled[name]
	}
	return s.SetMangled(name, s.NextMangled())
}

func (a *AstScope) References(name string) bool {
	return name != "" && a.Parent == nil || a.UsesWith || a.UsesEval || a.Refs[name] != nil
}

func (a *AstScope) Define(name string, t int) string {
	if name != "" {
		if _, ok := a.Names[name]; t == Type_Var || ok {
			a.Names[name] = t
		}
		return name
	}
	return ""
}

func NewScope(p *AstScope) *AstScope {
	s := &AstScope{}
	s.Names = make(map[string]int)
	s.Mangled = make(map[string]string)
	s.RevManged = make(map[string]string)
	s.CName = -1
	s.Refs = make(map[string]*AstScope)
	s.UsesWith = false
	s.UsesEval = false
	s.Directives = make([]string, 0)
	s.Parent = p
	s.Children = make([]*AstScope, 0)
	s.IsTrue = false
	s.Body = make([]IAst, 0)
	s.Labels = nil

	if p == nil {
		s.Level = 0
	} else {
		s.Level = p.Level + 1
		p.Children = append(p.Children, s)
	}
	return s
}

func TrueScope() *AstScope {
	s := &AstScope{}
	s.IsTrue = true
	return s
}

//-----------[ digits ]------------

var digits []rune

func base54(num int) string {
	ret := make([]rune, 0)
	base := 54

	for {
		ret = append(ret, digits[num%base])
		num = int(math.Floor(float64(num) / float64(base)))
		base = 64
		if num <= 0 {
			break
		}
	}
	return string(ret)
}

func init_digits() {
	digits = []rune{
		'e',
		't',
		'n',
		'r',
		'i',
		's',
		'o',
		'u',
		'a',
		'f',
		'l',
		'c',
		'h',
		'p',
		'd',
		'v',
		'm',
		'g',
		'y',
		'b',
		'w',
		'E',
		'S',
		'x',
		'T',
		'N',
		'C',
		'k',
		'L',
		'A',
		'O',
		'M',
		'_',
		'D',
		'P',
		'H',
		'B',
		'j',
		'F',
		'I',
		'q',
		'R',
		'U',
		'z',
		'W',
		'X',
		'V',
		'$',
		'J',
		'K',
		'Q',
		'G',
		'Y',
		'Z',
		'0',
		'5',
		'1',
		'6',
		'3',
		'7',
		'2',
		'9',
		'8',
		'4',
	}
}
