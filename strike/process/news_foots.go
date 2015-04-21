package process

import p "github.com/blinkat/blinks/strike/parser"

type news_foots struct {
	Default
	has_call bool
}

func (n *news_foots) Function(w *Walker, ast p.IAst) p.IAst {
	return ast
}

func (n *news_foots) Call(w *Walker, ast p.IAst) p.IAst {
	n.has_call = true
	return ast
}
