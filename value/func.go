package value

// Func type can return its value when the binary operator is
// evaluated with this type as its left operant.
type Func func(args ...Value) Value

var _ Value = Func(nil)

func (f Func) String() string         { return "func()" }
func (f Func) Eval(ctx Context) Value { return f }
func (f Func) whichType() valueType {
	return funcType
}

func (f Func) toType(which valueType) Value {
	switch which {
	case funcType:
		return f
	default:
		panic("type coversion from '" + f.whichType().String() + "' to '" + which.String() + "' not implemented")
	}
}
