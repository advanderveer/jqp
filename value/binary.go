package value

import (
	"github.com/advanderveer/jqp/token"
)

// A Binary operations takes two operations
type Binary struct {
	Op    token.TokenType
	Left  Expr
	Right Expr
}

// Eval will evaluate the binary operation
func (b *Binary) Eval(ctx Context) Value {
	op := binaryOps[b.Op]
	if op == nil {
		panic("binary op not implemented: " + b.Op.String())
	}

	// eval both sides
	rhs := b.Right.Eval(ctx)
	lhs := b.Left.Eval(ctx)

	// determine the bigger type
	bigger := op.biggerType(
		lhs.whichType(),
		rhs.whichType())

	// promote both sides to the bigger type
	lhs = lhs.toType(bigger)
	rhs = rhs.toType(bigger)

	// lookup the implementaiton for the type
	impl := op.impl[bigger]
	if impl == nil {
		panic("no implementation for op: " + b.Op.String() + " and the bigger type: " + bigger.String())
	}

	// call the actual implementation
	return impl(lhs, rhs)
}

// binaryOps holds all implementations for the binary operations
var binaryOps = map[token.TokenType]*binaryOp{

	// addition and string concat
	token.Add: &binaryOp{binaryArithType, [_numTypes]func(u, v Value) Value{
		intType:    func(u, v Value) Value { return Int(u.(Int) + v.(Int)) },
		floatType:  func(u, v Value) Value { return Float(u.(Float) + v.(Float)).shrink() },
		stringType: func(u, v Value) Value { return String(u.(String) + v.(String)) },
	}},

	// field reading
	token.Dot: &binaryOp{binaryArithType, [_numTypes]func(u, v Value) Value{
		mapType: mapKeyReadingImpl,
		portType: func(u, v Value) Value {
			up := u.(Port)
			vp := v.(Port)
			return up.cargo.Get(string(vp.cargo.Get("").(String)))
		},
	}},

	// index reading
	token.LBrack: &binaryOp{binaryArithType, [_numTypes]func(u, v Value) Value{
		arrayType: func(u, v Value) Value {
			ua := u.(Array)
			va := v.(Array)

			//we expect the right side 'v' to be an array of index integers
			//which will be used to read from the right side 'u' array of values
			var vals = make(Array, len(va))
			for i, iv := range va {
				ii, ok := iv.(Int)
				if !ok {
					panic("non integer index")
				}

				//@TODO make sure it is not out-of-range
				vals[i] = ua[int(ii)]
			}

			// shrink before returning
			return vals.shrink()
		},

		portType: func(u, v Value) Value {
			up := u.(Port)
			vp := v.(Port)
			idx := int(vp.cargo.Get("").(Int))
			return up.cargo.Range(idx, idx+1)
		},

		mapType: mapKeyReadingImpl,
	}},
}

func mapKeyReadingImpl(u, v Value) Value {
	um := u.(Map)
	vm := v.(Map)

	key := string(vm["key"].(String))
	val, ok := um[key]
	if !ok {
		panic("object doesn't have key: " + key)
	}

	return val
}

type binaryOp struct {
	biggerType func(a, b valueType) valueType    // which of the two binary operants to promote
	impl       [_numTypes]func(u, v Value) Value // implementations for each value type
}
