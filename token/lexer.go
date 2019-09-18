package token

import (
	"unicode"
	"unicode/utf8"
)

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

func lexAny(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(EOF)
		return nil //done
	case unicode.IsSpace(r):
		return lexSpace
	case unicode.IsDigit(r):
		return lexNumber
	case isAlphaNum(r):
		return lexIdent
	case l.lexOperator(r):
		return lexAny
	case r == '\'':
		return lexString
	case r == '[':
		l.emit(LBrack)
		return lexAny
	case r == ']':
		l.emit(RBrack)
		return lexAny
	case r == '(':
		l.emit(LParen)
		return lexAny
	case r == ')':
		l.emit(RParen)
		return lexAny
	default:
		panic("unrecognized character: " + string(r))
	}
}

// @TODO allow for rational numbers
// @TODO allow for negative numbers
func lexNumber(l *lexer) stateFn {
	var isFloat bool
	for {
		p := l.peek()
		if p != '.' && !unicode.IsDigit(p) {
			if isFloat {
				l.emit(Float)
			} else {
				l.emit(Int)
			}

			return lexAny
		}

		if p == '.' {
			isFloat = true
		}

		l.next()
	}
}

func lexString(l *lexer) stateFn {
	l.ignore()
	for {
		if l.peek() == '\'' {
			l.emit(String)
			l.next()
			l.ignore()
			return lexAny
		}

		l.next()
	}
}

func lexIdent(l *lexer) stateFn {
	for {
		if !isAlphaNum(l.peek()) {
			l.emit(Ident)
			return lexAny
		}

		l.next()
	}
}

func lexSpace(l *lexer) stateFn {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
	l.ignore()
	return lexAny
}

// Lex the input into token tokens
func Lex(input string) []Token {
	l := &lexer{
		input: input,
	}

	for state := lexAny; state != nil; {
		state = state(l)
	}

	return l.tokens
}

// lexer holds the state the lexing process
type lexer struct {
	input  string  // the string that is being scanned
	tokens []Token // the resulting tokens

	pos   int // zero-based index into the input
	start int // start position of this item
	width int // width of last rune read from input
}

var eof rune = -1

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() rune {
	rune := l.next()
	l.prev()
	return rune
}

// prev goes back one rune
func (l *lexer) prev() {
	l.pos -= l.width
}

// next returns the next rune and moves the
// current lexer position forward by its width
func (l *lexer) next() (rune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	rune, l.width =
		utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

// emit will append a result item while annotating the
// last scanned text with the provided token
func (l *lexer) emit(t TokenType) {
	l.tokens = append(l.tokens, Token{
		Type: t,
		Text: l.input[l.start:l.pos],
		Pos:  l.start,
	})

	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// lexOperator report whether r is the start of an operator the
// lexer ecountered an operator.
func (l *lexer) lexOperator(r rune) bool {
	switch r {

	// only support as single char operators:
	case '+':
		l.emit(Add)
		return true
	case '-':
		l.emit(Sub)
		return true
	case '*':
		l.emit(Mul)
		return true
	case '/':
		l.emit(Quo)
		return true
	case '%':
		l.emit(Rem)
		return true
	case '.':
		l.emit(Dot)
		return true
	case ',':
		l.emit(Comma)
		return true

	// supported as also as multi char operators
	case '>':
		if l.peek() == '=' {
			l.next()
			l.emit(GTE)
			return true
		}
		l.emit(GT)
		return true
	case '<':
		if l.peek() == '=' {
			l.next()
			l.emit(LTE)
			return true
		}
		l.emit(LT)
		return true
	case '!':
		if l.peek() == '=' {
			l.next()
			l.emit(NotEqual)
			return true
		}
		l.emit(Not)
		return true

	// supported only as multi char
	case '=':
		if l.peek() != '=' {
			return false
		}
		l.next()
		l.emit(Equal)
		return true

	default:
		return false
	}
}

// IsAlphaNum returns whether the rune is a valid start
// of a javascript identifier.
func isAlphaNum(r rune) bool {
	return r == '_' || r == '$' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
