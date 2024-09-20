package lexer

import (
	"reflect"
	"testing"

	"github.com/bootun/mini-tun/pkg/token"
)

func TestLexerParse(t *testing.T) {
	type fields struct {
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []token.Token
		wantErr bool
	}{
		{
			name: "basic_let_statement",
			fields: fields{
				input: "let x = 10",
			},
			want: []token.Token{
				token.New(token.LET, "let"),
				token.New(token.IDENTIFIER, "x"),
				token.New(token.EQUAL, "="),
				token.New(token.INT, "10"),
				token.New(token.EOF, ""),
			},
			wantErr: false,
		},
		{
			name: "func_declare_let_statement",
			fields: fields{
				input: `let x = function(bar,b,c){
	return bar+b-c
}`,
			},
			want: []token.Token{
				token.New(token.LET, "let"),
				token.New(token.IDENTIFIER, "x"),
				token.New(token.EQUAL, "="),
				token.New(token.FUNCTION, "function"),
				token.New(token.LPAREN, "("),
				token.New(token.IDENTIFIER, "bar"),
				token.New(token.COMMA, ","),
				token.New(token.IDENTIFIER, "b"),
				token.New(token.COMMA, ","),
				token.New(token.IDENTIFIER, "c"),
				token.New(token.RPAREN, ")"),
				token.New(token.LBRACE, "{"),
				token.New(token.RETURN, "return"),
				token.New(token.IDENTIFIER, "bar"),
				token.New(token.PLUS, "+"),
				token.New(token.IDENTIFIER, "b"),
				token.New(token.MINUS, "-"),
				token.New(token.IDENTIFIER, "c"),
				token.New(token.RBRACE, "}"),
				token.New(token.EOF, ""),
			},
			wantErr: false,
		},
		{
			name: "complex_let_statement",
			fields: fields{
				input: `let foo = function(a,b,c){
	return a+b-c
}
let a = 1
let b = 2
foo(a, foo(a,b))`,
			},
			want: []token.Token{
				token.New(token.LET, "let"),
				token.New(token.IDENTIFIER, "foo"),
				token.New(token.EQUAL, "="),
				token.New(token.FUNCTION, "function"),
				token.New(token.LPAREN, "("),
				token.New(token.IDENTIFIER, "a"),
				token.New(token.COMMA, ","),
				token.New(token.IDENTIFIER, "b"),
				token.New(token.COMMA, ","),
				token.New(token.IDENTIFIER, "c"),
				token.New(token.RPAREN, ")"),
				token.New(token.LBRACE, "{"),
				token.New(token.RETURN, "return"),
				token.New(token.IDENTIFIER, "a"),
				token.New(token.PLUS, "+"),
				token.New(token.IDENTIFIER, "b"),
				token.New(token.MINUS, "-"),
				token.New(token.IDENTIFIER, "c"),
				token.New(token.RBRACE, "}"),
				token.New(token.LET, "let"),
				token.New(token.IDENTIFIER, "a"),
				token.New(token.EQUAL, "="),
				token.New(token.INT, "1"),
				token.New(token.LET, "let"),
				token.New(token.IDENTIFIER, "b"),
				token.New(token.EQUAL, "="),
				token.New(token.INT, "2"),
				token.New(token.IDENTIFIER, "foo"),
				token.New(token.LPAREN, "("),
				token.New(token.IDENTIFIER, "a"),
				token.New(token.COMMA, ","),
				token.New(token.IDENTIFIER, "foo"),
				token.New(token.LPAREN, "("),
				token.New(token.IDENTIFIER, "a"),
				token.New(token.COMMA, ","),
				token.New(token.IDENTIFIER, "b"),
				token.New(token.RPAREN, ")"),
				token.New(token.RPAREN, ")"),
				token.New(token.EOF, ""),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.fields.input)
			got, err := l.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%v) got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
