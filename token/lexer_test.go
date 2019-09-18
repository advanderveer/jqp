package token_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/advanderveer/jqp/token"
)

// Test the lexing of character streams into tokens
func TestLexing(t *testing.T) {
	for i, c := range []struct {
		input  string
		tokens string
	}{
		{"", "[0:EOF]"},
		{" \t\n", "[3:EOF]"},
		{" \t(\n[)\t]", "[2:( 4:[ 5:) 7:] 8:EOF]"},
		{"11", "[0:Int(11) 2:EOF]"},
		{"1 1", "[0:Int(1) 2:Int(1) 3:EOF]"},
		{"1.1", "[0:Float(1.1) 3:EOF]"},
		{"'foo'", "[1:String(foo) 5:EOF]"},
		{"'fo o'", "[1:String(fo o) 6:EOF]"},
		{" $", "[1:Ident($) 2:EOF]"},
		{"a1 π", "[0:Ident(a1) 3:Ident(π) 5:EOF]"},
		{"+", "[0:+ 1:EOF]"}, {"-", "[0:- 1:EOF]"}, {"*", "[0:* 1:EOF]"}, {"/", "[0:/ 1:EOF]"}, {".", "[0:. 1:EOF]"}, {",", "[0:, 1:EOF]"},
		{"+%", "[0:+ 1:% 2:EOF]"},
		{">", "[0:> 1:EOF]"},
		{">='xyz'", "[0:>= 3:String(xyz) 7:EOF]"},
		{"<", "[0:< 1:EOF]"},
		{"<=10", "[0:<= 2:Int(10) 4:EOF]"},
		{"==1.0", "[0:== 2:Float(1.0) 5:EOF]"},
		{"!=$", "[0:!= 2:Ident($) 3:EOF]"},
		{"!1", "[0:! 1:Int(1) 2:EOF]"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tokens := token.Lex(c.input)

			//@TODO assert that token positions don't overlap

			if fmt.Sprintf("%s", tokens) != c.tokens {
				t.Fatalf("lexing '%s' caused: \n\t'%v' but expected:\n\t'%s'", c.input, tokens, c.tokens)
			}
		})
	}
}
