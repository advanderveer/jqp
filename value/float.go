package value

import (
	"strconv"
)

var _ Value = Float(0)

type Float float64

func (f Float) String() string         { return strconv.FormatFloat(float64(f), 'E', -1, 64) }
func (f Float) Eval(ctx Context) Value { return f }

func (f Float) whichType() valueType { return floatType }
func (f Float) toType(which valueType) Value {
	switch which {
	case floatType:
		return f
	case arrayType:
		return Array{f}
	default:
		panic("type coversion from '" + f.whichType().String() + "' to '" + which.String() + "' not implemented")
	}
}

func (f Float) shrink() Value {
	if float64(f) == float64(int(f)) {
		return Int(f)
	}
	return f
}
