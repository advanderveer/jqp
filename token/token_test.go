package token_test

import (
	"fmt"
	"testing"

	"github.com/advanderveer/jqp/token"
)

func TestTokenIsOperator(t *testing.T) {
	if token.EOF.IsOperator() == true {
		t.Fatal("eof should not be operator")
	}

	if token.Add.IsOperator() == false {
		t.Fatal("add should be operator")
	}
}

func TestTokenTypeString(t *testing.T) {
	if fmt.Sprint(token.Illegal) != "ILLEGAL" {
		t.Fatalf("should print token type correctly")
	}
	if fmt.Sprint(token.EOF) != "EOF" {
		t.Fatalf("should print token type correctly")
	}

	if fmt.Sprint(token.TokenType(999)) != "token(999)" {
		t.Fatalf("should print token type correctly")
	}

	tokens := []token.Token{
		{Type: token.Dot, Pos: 0},
		{Type: token.Ident, Text: "foo", Pos: 1},
		{Type: token.Dot, Pos: 2},
		{Type: token.Ident, Text: "bar", Pos: 3},
		{Type: token.EOF, Pos: 4},
	}
	if fmt.Sprint(tokens) != "[0:. 1:Ident(foo) 2:. 3:Ident(bar) 4:EOF]" {
		t.Fatalf("got: %s", tokens)
	}
}

func TestTokenString(t *testing.T) {
	if fmt.Sprint(token.Token{token.String, "foo", 10}) != "10:String(foo)" {
		t.Fatalf("expected correct token string")
	}

	if fmt.Sprint(token.Token{token.LBrack, "bogus", 5}) != "5:[" {
		t.Fatalf("expected correct token string")
	}
}
