package value

var _ Value = String("")

type String string

func (s String) String() string         { return string(s) }
func (s String) Eval(ctx Context) Value { return s }

func (s String) whichType() valueType { return stringType }
func (s String) toType(which valueType) Value {
	switch which {
	case stringType:
		return s
	case mapType:
		return Map{"key": s}
	case portType:

		// when the dot or index operator is used on a port, the
		// other operant must be upgraded from a string to a port.
		// As such we create a port that will always return
		// the key as a string so that the operator implementation
		// can be written
		return Port{stringPortCargo(s)}
	default:
		panic("type coversion from '" + s.whichType().String() + "' to '" + which.String() + "' not implemented")
	}
}
