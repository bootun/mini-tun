package lexer

import (
	"strings"

	"github.com/bootun/mini-tun/pkg/token"
)

type Lexer struct {
	input string
	pos   int
}

func New(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
	}
}

func (l *Lexer) Parse() ([]token.Token, error) {
	tokens := make([]token.Token, 0, 10)
	for tok := l.nextToken(); tok.GetType() != token.EOF; tok = l.nextToken() {
		tokens = append(tokens, tok)
	}
	tokens = append(tokens, token.New(token.EOF, ""))
	return tokens, nil
}

func (l *Lexer) nextToken() token.Token {
	if l.pos >= len(l.input) {
		return token.New(token.EOF, "")
	}

	// 吞掉空格
	for isSpace(l.input[l.pos]) {
		l.pos++
		if l.pos >= len(l.input) {
			return token.New(token.EOF, "")
		}
	}
	ch := l.readChar()
	switch ch {
	case '+':
		return token.New(token.PLUS, "+")
	case '-':
		return token.New(token.MINUS, "-")
	case '=':
		return token.New(token.EQUAL, "=")
	case '(':
		return token.New(token.LPAREN, "(")
	case ')':
		return token.New(token.RPAREN, ")")
	case '{':
		return token.New(token.LBRACE, "{")
	case '}':
		return token.New(token.RBRACE, "}")
	case ',':
		return token.New(token.COMMA, ",")
	default:
		var buf strings.Builder
		buf.WriteByte(ch)
		for i := l.pos; i < len(l.input) && !isSpace(l.input[i]); i++ {
			if !isLetter(l.input[i]) && !isNumber(l.input[i]) {
				break
			}
			buf.WriteByte(l.input[i])
			l.pos++
		}
		if buf.Len() == 0 {
			return token.New(token.EOF, "")
		}
		tk := buf.String()
		typ := token.LookupIdent(tk)
		return token.New(typ, tk)
	}

}
func (l *Lexer) readChar() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	pos := l.pos
	l.pos++
	return l.input[pos]
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t'
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}
