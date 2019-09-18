// +build wasm

package value

import (
	"syscall/js"
)

// FromJS turns a JavaScript value into a value of
// our own type system. It wil always create port
// values for javascript objects
func FromJS(jsv js.Value) Value {
	switch jsv.Type() {
	case js.TypeNull, js.TypeUndefined:
		// @TODO implement no value
		panic("js type was null/undefined")
	case js.TypeSymbol:
		//https: //developer.mozilla.org/en-US/docs/Glossary/Symbol
		panic("js symbols are not supported")
	case js.TypeBoolean:
		// @TODO implement boolean type
		panic("what to do with boolean values?")
	case js.TypeFunction:
		return Func(func(args ...Value) Value {
			res := make([]interface{}, len(args))
			for i := range args {
				res[i] = ToNative(args[i])
			}

			jsres := jsv.Invoke(res...)
			return FromJS(jsres)
		})

	case js.TypeObject:
		return Port{jsPortCargo{jsv}}
	case js.TypeString:
		return String(jsv.String())
	case js.TypeNumber:
		return Float(jsv.Float()).shrink() //possible shrink to int
	default:
		panic("unexpected js type, cannot convert:" + jsv.Type().String())
	}
}

// jsPortCargo is a port cargo implementation that
// implements range and get on the underlying javascript
// value
type jsPortCargo struct{ js.Value }

var _ PortCargo = jsPortCargo{}

func (p jsPortCargo) Get(k string) Value   { return FromJS(p.Value.Get(k)) }
func (p jsPortCargo) Range(i, j int) Value { return FromJS(p.Value.Index(i)) }
