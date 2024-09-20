package parser

import (
	"reflect"
	"testing"

	"github.com/bootun/mini-tun/pkg/ast"
	"github.com/bootun/mini-tun/pkg/lexer"
	"github.com/bootun/mini-tun/pkg/token"
)

func TestParser_Parse(t *testing.T) {
	type fields struct {
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		want    ast.Program
		wantErr bool
	}{
		{
			name: "basic_let_statement",
			fields: fields{
				"let a = 1",
			},
			want: ast.Program{
				Statements: []ast.Statement{
					ast.NewVariableAssignment("a", ast.NewLiteralExpression(1)),
				},
			},
			wantErr: false,
		},
		{
			name: "func_declare_statement",
			fields: fields{
				`
let add = function(a,b) {
	return a + b
}
`,
			},
			want: ast.Program{
				Statements: []ast.Statement{
					ast.NewVariableAssignment("add",
						ast.NewFunctionLiteral(
							[]*ast.IdentifierExpression{
								ast.NewIdentifierExpression("a"),
								ast.NewIdentifierExpression("b"),
							},
							ast.NewBlockStatement([]ast.Statement{
								ast.NewReturnStatement(
									ast.NewComplexExpression(
										ast.NewIdentifierExpression("a"),
										token.New(token.PLUS, "+"),
										ast.NewIdentifierExpression("b"),
									),
								),
							}),
						),
					),
				},
			},
			wantErr: false,
		}, {
			name: "func_declare_statement",
			fields: fields{
				`
let add = function(a,b) {
	return a + b
}
let a = 1
let b = 2
let c = add(a, add(a,10))
`,
			},
			want: ast.Program{
				Statements: []ast.Statement{
					ast.NewVariableAssignment("add",
						ast.NewFunctionLiteral(
							[]*ast.IdentifierExpression{
								ast.NewIdentifierExpression("a"),
								ast.NewIdentifierExpression("b"),
							},
							ast.NewBlockStatement([]ast.Statement{
								ast.NewReturnStatement(
									ast.NewComplexExpression(
										ast.NewIdentifierExpression("a"),
										token.New(token.PLUS, "+"),
										ast.NewIdentifierExpression("b"),
									),
								),
							}),
						),
					),
					ast.NewVariableAssignment("a", ast.NewLiteralExpression(1)),
					ast.NewVariableAssignment("b", ast.NewLiteralExpression(2)),
					ast.NewVariableAssignment("c", ast.NewFunctionCall("add", []ast.Expression{
						ast.NewIdentifierExpression("a"),
						ast.NewFunctionCall("add", []ast.Expression{
							ast.NewIdentifierExpression("a"),
							ast.NewLiteralExpression(10),
						}),
					})),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.New(tt.fields.input)
			p, err := New(lex)
			if err != nil {
				t.Errorf("New() error = %v", err)
				return
			}
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got.JSON(), tt.want)
			}
		})
	}
}
