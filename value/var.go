package value

var _ Expr = Var("")

type Var string

func (s Var) String() string {
	return string(s)
}

func (s Var) Eval(ctx Context) Value {
	v, ok := ctx.Decl[s]
	if !ok {
		panic("var not declared in context: " + string(s))
	}

	return v
}
