package value

import (
	"fmt"
)

type valueType int

// ToNative converts a value from the jqp type system
// to the native go type system.
func ToNative(v Value) interface{} {
	switch vt := v.(type) {
	case Int:
		return int(vt)
	case String:
		return string(vt)
	case Float:
		return float64(vt)
	case Array:
		res := make([]interface{}, len(vt))
		for i := range vt {
			res[i] = ToNative(vt[i])
		}
		return res
	case Map:
		res := make(map[string]interface{}, len(vt))
		for k := range vt {
			res[k] = ToNative(vt[k])
		}
		return res
	case Func:
		return func(args ...interface{}) interface{} {
			res := make([]Value, len(args))
			for i := range args {
				res[i] = FromNative(args[i], false)
			}

			return ToNative(vt(res...))
		}
	default:
		panic("jqp/value: cant convert value type '" + v.whichType().String() + "' to a native type")
	}
}

// FromNative transforms any native supported go types to
// a value that can be evaluated in jqp. If 'port' is
// set to true, map and slice will become ports values.
// Else they will be converted to map and arrays.
func FromNative(v interface{}, port bool) Value {
	switch vt := v.(type) {
	case Value:
		return vt
	case int:
		return Int(vt)
	case float64:
		return Float(vt)
	case string:
		return String(vt)
	case func(...interface{}) interface{}:
		return Func(func(args ...Value) Value {
			res := make([]interface{}, len(args))
			for i := range args {
				res[i] = ToNative(args[i])
			}

			return FromNative(vt(res...), port)
		})

	case []interface{}:
		if port {
			return Port{slicePortCargo(vt)}
		}

		v := make(Array, len(vt))
		for i := range v {
			v[i] = FromNative(vt[i], port)
		}

		return v
	case map[string]interface{}:
		if port {
			return Port{mapPortCargo(vt)}
		}

		v := make(Map, len(vt))
		for k := range vt {
			v[k] = FromNative(vt[k], port)
		}

		return v
	default:
		panic("jqp/value: cant convert this type from native:" + fmt.Sprintf("%T", vt))
	}
}

const (
	intType valueType = iota
	floatType
	stringType
	arrayType
	mapType
	portType
	funcType

	_numTypes //number of types
)

func (vt valueType) String() string {
	var typeName = [_numTypes]string{"int", "float", "string", "array", "map", "port", "func"}
	return typeName[vt]
}

type Context struct {
	Decl map[Var]Value
}

type Expr interface {
	Eval(ctx Context) Value
}

type Value interface {
	Expr
	String() string

	whichType() valueType
	toType(valueType) Value
}

func binaryArithType(t1, t2 valueType) valueType {
	if t1 > t2 {
		return t1
	}
	return t2
}
