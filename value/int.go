package value

import (
	"strconv"
)

var _ Value = Int(0)

type Int int64

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Int) Eval(ctx Context) Value {
	return i
}

func (i Int) whichType() valueType {
	return intType
}

func (i Int) toType(which valueType) Value {
	switch which {
	case intType:
		return i
	case floatType:
		return Float(float64(i))
	case arrayType:
		return Array{i}
	case portType:

		// when the index operator is used  on a port with an int
		// operant, it must be upgraded from to a port.
		// As such we create a port that will always return
		// the key as a string so that the operator implementation
		// can be written
		return Port{intPortCargo(i)}
	default:
		panic("type coversion from '" + i.whichType().String() + "' to '" + which.String() + "' not implemented")
	}
}
