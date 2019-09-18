package value

import (
	"strings"
)

// Array is an indexable and iterable collection of values
type Array []Value

var _ Value = Array{}

func (a Array) String() string {
	vals := make([]string, len(a))
	for i := range a {
		vals[i] = a[i].String()
	}

	return "[" + strings.Join(vals, ", ") + "]"
}

func (a Array) Eval(ctx Context) Value {
	return a
}

func (a Array) whichType() valueType {
	return arrayType
}

func (a Array) toType(which valueType) Value {
	switch which {
	case arrayType:
		return a
	default:
		panic("type coversion from '" + a.whichType().String() + "' to '" + which.String() + "' not implemented")
	}
}

func (a Array) shrink() Value {
	if len(a) == 1 {
		return a[0]
	}
	return a
}
