package value_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/advanderveer/jqp"
	"github.com/advanderveer/jqp/token"
	"github.com/advanderveer/jqp/value"
)

func testContext(v value.Value) value.Context {
	return value.Context{Decl: map[value.Var]value.Value{
		value.Var("$"): v,
	}}
}

func TestBinaryEvalWithoutInput(t *testing.T) {
	type c struct {
		ctx  value.Context
		expr value.Expr
		out  value.Value
	}

	var evalCases = []c{
		{value.Context{}, value.Int(1), value.Int(1)},
	}

	var addOpCases = []c{

		// adding with literals
		{value.Context{}, &value.Binary{
			Left:  value.Int(1),
			Right: value.Int(2),
			Op:    token.Add,
		}, value.Int(3)},
		{value.Context{}, &value.Binary{
			Left:  value.String("foo"),
			Right: value.String("bar"),
			Op:    token.Add,
		}, value.String("foobar")},
		{value.Context{}, &value.Binary{
			Left:  value.Int(1),
			Right: value.Float(2.1),
			Op:    token.Add,
		}, value.Float(3.1)},

		// adding with context identifier
		{value.Context{Decl: map[value.Var]value.Value{
			value.Var("foo"): value.Int(4),
			value.Var("bar"): value.Int(3),
		}}, &value.Binary{
			Left:  value.Var("foo"),
			Right: value.Var("bar"),
			Op:    token.Add,
		}, value.Int(7)},
	}

	var indexCases = []c{
		// index reading on array
		{testContext(value.Array{value.Int(10)}), &value.Binary{
			Left:  value.Var("$"),
			Op:    token.LBrack,
			Right: value.Int(0),
		}, value.Int(10)},

		// index reading on map
		{testContext(value.Map{"foo": value.Int(12)}), &value.Binary{
			Left:  value.Var("$"),
			Op:    token.LBrack,
			Right: value.String("foo"),
		}, value.Int(12)},

		// nested index reading
		{testContext(
			value.Map{"foo": value.Array{value.Map{"bar": value.Int(100)}}},
		), &value.Binary{
			Left: &value.Binary{
				Left: &value.Binary{
					Left:  value.Var("$"),
					Op:    token.LBrack,
					Right: value.String("foo")},
				Op:    token.LBrack,
				Right: value.Int(0)},
			Op:    token.LBrack,
			Right: value.String("bar")}, value.Int(100)},
	}

	var fieldCases = []c{

		// field reading on map
		{testContext(value.Map{"foo": value.Int(13)}), &value.Binary{
			Left:  value.Var("$"),
			Op:    token.Dot,
			Right: value.String("foo"),
		}, value.Int(13)},

		// nested field reading
		{testContext(
			value.Map{"foo": value.Map{"foo2": value.Map{"bar": value.Int(101)}}},
		), &value.Binary{
			Left: &value.Binary{
				Left: &value.Binary{
					Left:  value.Var("$"),
					Op:    token.Dot,
					Right: value.String("foo")},
				Op:    token.Dot,
				Right: value.String("foo2")},
			Op:    token.Dot,
			Right: value.String("bar")}, value.Int(101)},

		// ported nested field reading
		{testContext(
			value.FromNative(map[string]interface{}{
				"foo": map[string]interface{}{
					"foo2": map[string]interface{}{
						"bar": 102,
					},
				},
			}, true),
		), &value.Binary{
			Left: &value.Binary{
				Left: &value.Binary{
					Left:  value.Var("$"),
					Op:    token.Dot,
					Right: value.String("foo")},
				Op:    token.Dot,
				Right: value.String("foo2")},
			Op:    token.Dot,
			Right: value.String("bar")}, value.Int(102)},

		// ported nested field and index reading
		{testContext(
			value.FromNative(map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"bar": 102,
					},
				},
			}, true),
		), &value.Binary{
			Left: &value.Binary{
				Left: &value.Binary{
					Left:  value.Var("$"),
					Op:    token.Dot,
					Right: value.String("foo")},
				Op:    token.LBrack,
				Right: value.Int(0)},
			Op:    token.Dot,
			Right: value.String("bar")}, value.Int(102)},
	}

	var callCases = []c{

		// function type calling
		{testContext(
			value.FromNative(func(args ...interface{}) interface{} {
				return args[0].(int) + args[1].(int) + 103
			}, false),
		), &value.Call{
			Func: value.Var("$"),
			Args: []value.Expr{value.Int(5), value.Int(15)}}, value.Int(123)},

		// nested function calling
		{testContext(
			value.FromNative(map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"bar": func(...interface{}) interface{} { return 100 },
					},
				},
			}, true),
		), &value.Call{
			Func: &value.Binary{
				Left: &value.Binary{
					Left: &value.Binary{
						Left:  value.Var("$"),
						Op:    token.Dot,
						Right: value.String("foo")},
					Op:    token.LBrack,
					Right: value.Int(0)},
				Op:    token.Dot,
				Right: value.String("bar")},
			Args: []value.Expr{}}, value.Int(100)},
	}

	for i, c := range [][]c{
		evalCases,
		addOpCases,
		indexCases,
		fieldCases,
		callCases,
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			for _, c := range c {
				func() {
					defer func() {
						if r := recover(); r != nil {
							t.Fatalf("panic while evaluating '%s':\n\t %v", jqp.Format(c.expr), r)
						}
					}()

					res := c.expr.Eval(c.ctx)
					if !reflect.DeepEqual(res, c.out) {
						t.Fatalf("evaluating '%s' gave: '%v' expected: '%v' ", jqp.Format(c.expr), res, c.out)
					}

				}()
			}
		})
	}
}
