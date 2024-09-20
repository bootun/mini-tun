package token

import (
	"strconv"
)

type TokenType string

const (
	// 特殊
	EOF        TokenType = "EOF"
	IDENTIFIER TokenType = "IDENTIFIER" // 标识符
	FUNCTION   TokenType = "FUNCTION"   // function
	LET        TokenType = "LET"        // let
	RETURN     TokenType = "RETURN"     // return
	EQUAL      TokenType = "EQUAL"      // =
	LPAREN     TokenType = "LPAREN"     // (
	RPAREN     TokenType = "RPAREN"     // )
	LBRACE     TokenType = "LBRACE"     // {
	RBRACE     TokenType = "RBRACE"     // }
	COMMA      TokenType = "COMMA"      // ,
	COMMENT    TokenType = "COMMENT"    // //

	// 操作符
	PLUS  TokenType = "PLUS"  // +
	MINUS TokenType = "MINUS" // -

	// 类型
	INT TokenType = "INT" // int
)

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) GetType() TokenType {
	return t.Type
}

func (t Token) GetLiteral() string {
	return t.Literal
}

func New(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal}
}

var keywords = map[string]TokenType{
	"function": FUNCTION,
	"let":      LET,
	"return":   RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	if isInt(ident) {
		return INT
	}
	return IDENTIFIER
}

func isInt(ident string) bool {
	_, err := strconv.ParseInt(ident, 10, 64)
	return err == nil
}
