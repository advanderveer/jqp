package jqp_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/advanderveer/jqp"
	"github.com/advanderveer/jqp/token"
	"github.com/advanderveer/jqp/value"
)

func TestExprFormating(t *testing.T) {
	for i, c := range []struct {
		tokens []token.Token
		expr   string
	}{
		{[]token.Token{
			{Type: token.Ident, Text: "foo"},
			{Type: token.EOF},
		}, `<var foo>`},

		{[]token.Token{
			{Type: token.String, Text: "foo"},
			{Type: token.EOF},
		}, `<string foo>`},

		{[]token.Token{
			{Type: token.Int, Text: "1"},
			{Type: token.EOF},
		}, `<int 1>`},

		{[]token.Token{
			{Type: token.Float, Text: "1.5"},
			{Type: token.EOF},
		}, `<float 1.5E+00>`},

		{[]token.Token{
			{Type: token.Not},
			{Type: token.Int, Text: "1"},
			{Type: token.EOF},
		}, `(! <int 1>)`},

		{[]token.Token{
			{Type: token.Int, Text: "1"},
			{Type: token.Add},
			{Type: token.Int, Text: "2"},
			{Type: token.EOF},
		}, `(<int 1> + <int 2>)`},

		// {[]token.Token{
		// 	{Type: token.Int, Text: "1"},
		// 	{Type: token.Comma},
		// 	{Type: token.Int, Text: "2"},
		// 	{Type: token.EOF},
		// }, `(<int 1> , <int 2>)`},

		// deep indexing: $[foo][bar][rab]
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.LBrack},
			{Type: token.String, Text: "foo"},
			{Type: token.RBrack},
			{Type: token.LBrack},
			{Type: token.String, Text: "bar"},
			{Type: token.RBrack},
			{Type: token.LBrack},
			{Type: token.String, Text: "rab"},
			{Type: token.RBrack},
			{Type: token.EOF},
		}, `(((<var $>[<string foo>])[<string bar>])[<string rab>])`},

		// deep field reading: $.foo.foo.foo
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "foo"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "bar"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "rab"},
		}, `(((<var $> . <string foo>) . <string bar>) . <string rab>)`},

		// $.foo['bar']['bar']
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "foo"},
			{Type: token.LBrack},
			{Type: token.String, Text: "bar"},
			{Type: token.RBrack},
			{Type: token.LBrack},
			{Type: token.String, Text: "bar"},
			{Type: token.RBrack},
		}, `(((<var $> . <string foo>)[<string bar>])[<string bar>])`},

		// $.foo['bar'].foo
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "foo"},
			{Type: token.LBrack},
			{Type: token.String, Text: "bar"},
			{Type: token.RBrack},
			{Type: token.Dot},
			{Type: token.Ident, Text: "foo"},
		}, `(((<var $> . <string foo>)[<string bar>]) . <string foo>)`},

		// function calling on identifier directly
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.LParen},
			{Type: token.String, Text: "arg1"},
			{Type: token.RParen},
		}, `(<var $>(<string arg1>))`},

		// calling with two args
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.LParen},
			{Type: token.String, Text: "arg1"},
			{Type: token.Comma},
			{Type: token.Int, Text: "1"},
			{Type: token.Add},
			{Type: token.Float, Text: "5.3"},
			{Type: token.Comma},
			{Type: token.String, Text: "arg3"},
			{Type: token.RParen},
		}, `(<var $>(<string arg1>, (<int 1> + <float 5.3E+00>), <string arg3>))`},

		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.LParen},
			{Type: token.RParen},
		}, `(<var $>())`},

		// function calling on index
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.LBrack},
			{Type: token.String, Text: "bar"},
			{Type: token.RBrack},
			{Type: token.LParen},
			{Type: token.String, Text: "arg1"},
			{Type: token.RParen},
		}, `((<var $>[<string bar>])(<string arg1>))`},

		// read field on function result
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.LBrack},
			{Type: token.String, Text: "bar"},
			{Type: token.RBrack},
			{Type: token.LParen},
			{Type: token.String, Text: "arg1"},
			{Type: token.RParen},
			{Type: token.Dot},
			{Type: token.Ident, Text: "foo"},
		}, `(((<var $>[<string bar>])(<string arg1>)) . <string foo>)`},

		// read index on function result with two args
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.LBrack},
			{Type: token.String, Text: "bar"},
			{Type: token.RBrack},
			{Type: token.LParen},
			{Type: token.String, Text: "arg1"},
			{Type: token.Comma},
			{Type: token.Int, Text: "100"},
			{Type: token.RParen},
			{Type: token.LBrack},
			{Type: token.String, Text: "foo"},
			{Type: token.RBrack},
		}, `(((<var $>[<string bar>])(<string arg1>, <int 100>))[<string foo>])`},

		// function calling on field
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "foo"},
			{Type: token.LParen},
			{Type: token.String, Text: "arg1"},
			{Type: token.RParen},
		}, `((<var $> . <string foo>)(<string arg1>))`},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic while parsing tokens '%s':\n\t %v", c.tokens, r)
				}
			}()

			res := jqp.Parse(c.tokens)
			if jqp.Format(res) != c.expr {
				t.Fatalf("tokens '%s' should result in expr: \n\t %s got: \n\t %s", c.tokens, c.expr, jqp.Format(res))
			}
		})
	}
}

func TestExprParsing(t *testing.T) {
	for i, c := range []struct {
		tokens []token.Token
		expr   value.Expr
	}{
		// basic value expressions
		{[]token.Token{
			{Type: token.String, Text: "foo"},
			{Type: token.EOF},
		}, value.String("foo")},

		{[]token.Token{
			{Type: token.Int, Text: "1"},
			{Type: token.EOF},
		}, value.Int(1)},

		{[]token.Token{
			{Type: token.Float, Text: "1.5"},
			{Type: token.EOF},
		}, value.Float(1.5)},

		// unary operation
		{[]token.Token{
			{Type: token.Not},
			{Type: token.Int, Text: "1"},
			{Type: token.EOF},
		}, &value.Unary{Op: token.Not, Right: value.Int(1)}},

		// binary operation
		{[]token.Token{
			{Type: token.Int, Text: "1"},
			{Type: token.Add},
			{Type: token.Int, Text: "2"},
			{Type: token.EOF},
		}, &value.Binary{
			Left:  value.Int(1),
			Op:    token.Add,
			Right: value.Int(2)}},

		// parenthese grouping
		{[]token.Token{
			{Type: token.LParen},
			{Type: token.String, Text: "foo"},
			{Type: token.RParen},
			{Type: token.EOF},
		}, value.String("foo")},

		{[]token.Token{
			{Type: token.LParen},
			{Type: token.Int, Text: "3"},
			{Type: token.Mul},
			{Type: token.Int, Text: "5"},
			{Type: token.RParen},
			{Type: token.Add},
			{Type: token.Int, Text: "2"},
			{Type: token.EOF},
		}, &value.Binary{
			Left: &value.Binary{
				Left:  value.Int(3),
				Op:    token.Mul,
				Right: value.Int(5),
			},
			Op:    token.Add,
			Right: value.Int(2)}},

		// indexing
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.LBrack},
			{Type: token.String, Text: "foo"},
			{Type: token.RBrack},
			{Type: token.LBrack},
			{Type: token.Int, Text: "0"},
			{Type: token.RBrack},
			{Type: token.LBrack},
			{Type: token.String, Text: "bar"},
			{Type: token.RBrack},
			{Type: token.EOF},
		}, &value.Binary{
			Left: &value.Binary{
				Left: &value.Binary{
					Left:  value.Var("$"),
					Op:    token.LBrack,
					Right: value.String("foo")},
				Op:    token.LBrack,
				Right: value.Int(0)},
			Op:    token.LBrack,
			Right: value.String("bar")}},

		// field reading
		{[]token.Token{
			{Type: token.Ident, Text: "$"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "foo"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "bar"},
			{Type: token.Dot},
			{Type: token.Ident, Text: "foobar"},
			{Type: token.EOF},
		}, &value.Binary{
			Left: &value.Binary{
				Left: &value.Binary{
					Left:  value.Var("$"),
					Op:    token.Dot,
					Right: value.String("foo")},
				Op:    token.Dot,
				Right: value.String("bar")},
			Op:    token.Dot,
			Right: value.String("foobar")}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic: %v", r)
				}
			}()

			res := jqp.Parse(c.tokens)
			if !reflect.DeepEqual(res, c.expr) {
				t.Fatalf("tokens '%s' should result in expr: \n\t %v got: \n\t %v", c.tokens, jqp.Format(c.expr), jqp.Format(res))
			}
		})
	}
}
