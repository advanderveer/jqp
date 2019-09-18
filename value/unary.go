package value

import (
	"github.com/advanderveer/jqp/token"
)

// parsed unary operation
type Unary struct {
	Op    token.TokenType
	Right Expr
}

func (u *Unary) Eval(ctx Context) Value {
	return nil
}
