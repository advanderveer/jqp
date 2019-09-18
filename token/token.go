package token

import (
	"strconv"
)

// TokenType describes the type of token
type TokenType int

const (
	// Special Tokens
	Illegal TokenType = iota
	EOF

	// Literal types
	Ident  // main
	Int    // 42
	Float  // 42.3
	String // 42

	// grouping, calls, indexing, fields and func args
	LParen // '('
	RParen // ')'
	LBrack // '['
	RBrack // ']'
	Dot    // '.'
	Comma  // ','
	// Pipe   // '|' @TODO add support

	// basic operators
	_operator_beg
	Add // +
	Sub // -
	Mul // *
	Quo // /
	Rem // %

	Equal    // ==
	NotEqual // !=
	LTE      // <=
	GTE      // >=
	LT       // <
	GT       // >
	Not      // !
	_operator_end
)

func (tt TokenType) IsOperator() bool {
	return _operator_beg < tt && tt < _operator_end
}

func (tt TokenType) String() string {
	var tokens = map[TokenType]string{
		Illegal: "ILLEGAL",
		EOF:     "EOF",

		Int:    "Int",
		Float:  "Float",
		String: "String",
		Ident:  "Ident",

		LParen: "(",
		RParen: ")",
		LBrack: "[",
		RBrack: "]",

		Add:   "+",
		Sub:   "-",
		Mul:   "*",
		Quo:   "/",
		Rem:   "%",
		Dot:   ".",
		Comma: ",",

		// Pipe: "|",

		Equal:    "==",
		NotEqual: "!=",
		LTE:      "<=",
		GTE:      ">=",
		LT:       "<",
		GT:       ">",
		Not:      "!",
	}

	s, ok := tokens[tt]
	if !ok {
		return "token(" + strconv.Itoa(int(tt)) + ")"
	}

	return s
}

type Token struct {
	Type TokenType
	Text string
	Pos  int
}

func (tok Token) String() string {
	s := strconv.Itoa(tok.Pos) + ":" + tok.Type.String()
	switch tok.Type {
	case Int, Float, String, Ident:
		return s + "(" + tok.Text + ")"
	default:
		return s
	}
}
