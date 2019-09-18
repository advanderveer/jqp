package jqp

import (
	"strconv"
	"strings"

	"github.com/advanderveer/jqp/token"
	"github.com/advanderveer/jqp/value"
)

// Parse the scanned tokens into an expression
func Parse(input []token.Token) value.Expr {
	return (&parser{input}).expr()
}

// Format an expression in an unambiguous form for debugging.
func Format(e interface{}) string {
	switch e := e.(type) {
	case value.Int:
		return "<int " + e.String() + ">"
	case value.Float:
		return "<float " + e.String() + ">"
	case value.String:
		return "<string " + e.String() + ">"
	case value.Var:
		return "<var " + e.String() + ">"
	case []value.Expr:
		var s []string
		for _, ee := range e {
			s = append(s, Format(ee))
		}
		return strings.Join(s, ", ")
	case *value.Call:
		return "(" + Format(e.Func) + "(" + Format(e.Args) + "))"
	case *value.Unary:
		return "(" + e.Op.String() + " " + Format(e.Right) + ")"
	case *value.Binary:
		if e.Op == token.LBrack {
			return "(" + Format(e.Left) + "[" + Format(e.Right) + "])"
		}

		return "(" + Format(e.Left) + " " + e.Op.String() + " " + Format(e.Right) + ")"
	default:
		panic("not able to format provided type as expression")
	}
}

type parser struct {
	rem []token.Token
}

func (p *parser) next() (tok token.Token) {
	tok = p.peek()
	if tok.Type != token.EOF {
		p.rem = p.rem[1:]
	}

	if tok.Type == token.Illegal {
		panic("jqp/parser: illegal token")
	}

	return
}

func (p *parser) peek() token.Token {
	if len(p.rem) == 0 {
		return token.Token{Type: token.EOF}
	}
	return p.rem[0]
}

// expr
//	[x] operand
//	[x] operand op expr
func (p *parser) expr() value.Expr {
	tok := p.next()
	expr := p.operand(tok)

	peeked := p.peek()
	switch peeked.Type {
	case token.EOF, token.RParen, token.RBrack, token.Comma:
		return expr
	}

	if peeked.Type.IsOperator() {
		p.next()
		return &value.Binary{
			Left:  expr,
			Op:    peeked.Type,
			Right: p.expr(),
		}
	}

	panic("jqp/parser: unexpected token '" + p.peek().String() + "', expression so far: " + Format(expr))
}

// operand
//	[x] op Expr
//	[x] literal
//	[x] operand [ Expr ]...
//	[x] operand ( )...
//	[x] operand . Ident ...
func (p *parser) operand(tok token.Token) value.Expr {

	// if its an operator, return as unary with the right
	// set by parsing the remaining tokens as an expression
	if tok.Type.IsOperator() {
		return &value.Unary{
			Op:    tok.Type,
			Right: p.expr(),
		}
	}

	// else, the operand is a literal or grouped expression
	var expr value.Expr
	switch tok.Type {
	case token.Int, token.Float, token.String, token.Ident, token.LParen:
		expr = p.literal(tok)
	}

	// check if the current operant has an index, call or field operator
	// following it. If so keep adding operator expresisons until it is
	// stable.
	last := expr
	for {
		expr = p.index(last)
		expr = p.call(expr)
		expr = p.field(expr)
		if last == expr {
			return expr
		}

		last = expr
	}
}

func (p *parser) index(expr value.Expr) value.Expr {
	for p.peek().Type == token.LBrack {
		p.next()
		index := p.expr()
		tok := p.next()
		if tok.Type != token.RBrack {
			panic("jqp/parser: expected right bracket, found: " + tok.String())
		}

		expr = &value.Binary{
			Op:    token.LBrack,
			Left:  expr,
			Right: index,
		}
	}

	return expr
}

func (p *parser) field(expr value.Expr) value.Expr {
	for p.peek().Type == token.Dot {
		p.next() //the dot

		tok := p.next()
		switch tok.Type {
		case token.Ident:
			expr = &value.Binary{
				Op:    token.Dot,
				Left:  expr,
				Right: value.String(tok.Text),
			}

		default:
			panic("expected identifier after dot, got: " + tok.String())
		}
	}

	return expr
}

// call
//  [x] ()
//  [x] (x, ...)
func (p *parser) call(expr value.Expr) value.Expr {
	for p.peek().Type == token.LParen {
		p.next()

		// start parsing arguments
		var args []value.Expr
		for {
			peeked := p.peek()
			if peeked.Type == token.RParen {
				p.next()
				break
			} else if p.peek().Type == token.Comma {
				p.next()
				continue
			}

			args = append(args, p.expr())
		}

		expr = &value.Call{
			Func: expr,
			Args: args,
		}
	}

	return expr
}

// literal
//  [x] var
// 	[x] string
// 	[x] int
//	[x] float
//  [x] '(' expr ')'
func (p *parser) literal(tok token.Token) value.Expr {
	switch tok.Type {
	case token.Ident:
		return value.Var(tok.Text)
	case token.String:
		return value.String(tok.Text)
	case token.Int:
		i64, err := strconv.ParseInt(tok.Text, 10, 64)
		if err != nil {
			panic("jqp/parser: couldn't parse literal token as int: " + err.Error())
		}

		return value.Int(i64)
	case token.Float:
		f64, err := strconv.ParseFloat(tok.Text, 64)
		if err != nil {
			panic("jqp/parser: couldn't parse literal token as float: " + err.Error())
		}

		return value.Float(f64)
	case token.LParen:
		expr := p.expr()
		tok := p.next()
		if tok.Type != token.RParen {
			panic("jqp/parser: expected right parentheses, found: " + tok.String())
		}
		return expr
	}

	panic("unexpected token in literal, got: " + tok.String())
}
