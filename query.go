package jqp

import (
	"github.com/advanderveer/jqp/token"
	"github.com/advanderveer/jqp/value"
)

func Query(q string, v interface{}) interface{} {
	tokens := token.Lex(q)                         // lex
	expr := Parse(tokens)                          // parse
	return value.ToNative(expr.Eval(value.Context{ // eval
		Decl: map[value.Var]value.Value{"$": value.FromNative(v, false)},
	}))
}
