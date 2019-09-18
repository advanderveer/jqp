package value

// A Call operations takes two operations
type Call struct {
	Func Expr
	Args []Expr
}

// Eval will evaluate the binary operation
func (c *Call) Eval(ctx Context) Value {
	f := c.Func.Eval(ctx).(Func)
	argv := make([]Value, len(c.Args))
	for i := range c.Args {
		argv[i] = c.Args[i].Eval(ctx)
	}

	return f(argv...)
}
