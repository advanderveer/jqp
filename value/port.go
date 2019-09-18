package value

// slicePortCargo is a cargo implementation that
// provides range access for the index operator
type slicePortCargo []interface{}

func (p slicePortCargo) Get(k string) Value   { panic("get on slice port cargo is not supported") }
func (p slicePortCargo) Range(i, j int) Value { return FromNative(p[i], true) }

var _ PortCargo = slicePortCargo{}

// mapPortCargo provides a concrete port cargo
// for interface values mapped by strings
type mapPortCargo map[string]interface{}

func (p mapPortCargo) Range(i, j int) Value { panic("range on map port cargo is not supported") }
func (p mapPortCargo) Get(k string) Value   { return FromNative(p[k], true) }

var _ PortCargo = mapPortCargo{}

// provides a cargo implementation such that the dot
// string operant can be upgraded for the operators
// implementation
type stringPortCargo String

func (p stringPortCargo) Range(i, j int) Value { panic("range on int port cargo is not supported") }
func (p stringPortCargo) Get(k string) Value {
	return String(p)
}

var _ PortCargo = stringPortCargo("")

// provides a cargo implementation such that the index
// int operant can be upgraded for the operators
// implementation
type intPortCargo Int

func (p intPortCargo) Range(i, j int) Value { panic("range on int port cargo is not supported") }
func (p intPortCargo) Get(k string) Value   { return Int(p) }

var _ PortCargo = intPortCargo(0)

// PortCargo is the value held in the port
type PortCargo interface {
	Range(i, j int) Value
	Get(k string) Value
}

// Port is a value type that holds a reference to
// another (opaque) value while still being able to
// providing '.' and '[]' operator implementations
type Port struct{ cargo PortCargo }

var _ Value = Port{}

func (o Port) String() string         { return "{ port }" }
func (o Port) Eval(ctx Context) Value { return o }

func (o Port) whichType() valueType { return portType }
func (o Port) toType(which valueType) Value {
	switch which {
	case portType:
		return o
	default:
		panic("type coversion from '" + o.whichType().String() + "' to '" + which.String() + "' not implemented")
	}
}
