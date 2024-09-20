package parser

import (
	"fmt"
	"strconv"

	"github.com/bootun/mini-tun/pkg/ast"
	"github.com/bootun/mini-tun/pkg/lexer"
	"github.com/bootun/mini-tun/pkg/token"
)

type Parser struct {
	tokens []token.Token
	curPos int
}

func New(l *lexer.Lexer) (*Parser, error) {
	p := &Parser{}
	tokens, err := l.Parse()
	if err != nil {
		return nil, fmt.Errorf("lexer parse token error: %v", err)
	}

	p.tokens = tokens
	return p, nil
}

func (p *Parser) Parse() (ast.Program, error) {
	var program ast.Program
	for p.curPos < len(p.tokens) {
		if p.tokens[p.curPos].GetType() == token.EOF {
			break
		}
		statement, err := p.parseStatement()
		if err != nil {
			return program, fmt.Errorf("parse statement error: %v", err)
		}
		program.Statements = append(program.Statements, statement)
	}
	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.tokens[p.curPos].GetType() {
	case token.LET:
		statement, err := p.parseLetStatement()
		if err != nil {
			return nil, fmt.Errorf("parse let statement error: %v", err)
		}
		return statement, nil
	case token.RETURN:
		statement, err := p.parseReturnStatement()
		if err != nil {
			return nil, fmt.Errorf("parse return statement error: %v", err)
		}
		return statement, nil
	default:
		return nil, fmt.Errorf("unexpected token type: %v", p.tokens[p.curPos].GetLiteral())
	}
}

func (p *Parser) parseLetStatement() (ast.Statement, error) {
	if p.tokens[p.curPos].GetType() != token.LET {
		return nil, fmt.Errorf("invalid token type, expected let, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	p.curPos++
	if p.tokens[p.curPos].GetType() != token.IDENTIFIER {
		return nil, fmt.Errorf("invalid token type, expected identifier, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	letStatement := ast.NewVariableAssignment(p.tokens[p.curPos].GetLiteral(), nil)
	p.curPos++
	if p.tokens[p.curPos].GetType() != token.EQUAL {
		return nil, fmt.Errorf("invalid token type, expected equal, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	p.curPos++
	expression, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("parse literal error: %v", err)
	}
	letStatement.Value = expression
	return letStatement, nil
}

func (p *Parser) parseReturnStatement() (ast.Statement, error) {
	if p.tokens[p.curPos].GetType() != token.RETURN {
		return nil, fmt.Errorf("invalid token type, expected return, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	p.curPos++
	expression, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("parse literal error: %v", err)
	}
	return ast.NewReturnStatement(expression), nil
}
func (p *Parser) parseExpression() (ast.Expression, error) {
	curToken := p.tokens[p.curPos]
	nextToken := p.peekToken()
	switch nextToken.GetType() {
	case token.PLUS, token.MINUS:
		return p.parseComplexExpression()
	}
	switch curToken.GetType() {
	case token.FUNCTION:
		// 解析函数
		expression, err := p.parseFunctionDeclareExpression()
		if err != nil {
			return nil, fmt.Errorf("parse function expression error: %v", err)
		}
		return expression, nil
	case token.INT:
		return p.parseLiteralExpression()
	case token.IDENTIFIER:
		if p.tokens[p.curPos+1].GetType() == token.LPAREN {
			return p.parseFunctionCallExpression()
		}
		return p.parseIdentifierExpression()
	default:
		return nil, fmt.Errorf("unexpected token type: %v", curToken.GetLiteral())
	}
}

func (p *Parser) parseFunctionCallExpression() (ast.Expression, error) {
	f := ast.NewFunctionCall(p.tokens[p.curPos].GetLiteral(), nil)
	p.curPos++
	if p.tokens[p.curPos].GetType() != token.LPAREN {
		return nil, fmt.Errorf("invalid token type, expected left parenthesis, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	p.curPos++
	for p.tokens[p.curPos].GetType() != token.RPAREN {
		if p.tokens[p.curPos].GetType() == token.COMMA {
			p.curPos++
		} else if p.tokens[p.curPos].GetType() == token.RPAREN {
			break
		}
		expression, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("parse expression error: %v", err)
		}
		f.Arguments = append(f.Arguments, expression)
	}
	p.curPos++

	return f, nil
}

func (p *Parser) parseFunctionDeclareExpression() (ast.Expression, error) {
	if p.tokens[p.curPos].GetType() != token.FUNCTION {
		return nil, fmt.Errorf("invalid token type, expected function, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	function := ast.NewFunctionLiteral(nil, nil)
	p.curPos++
	if p.tokens[p.curPos].GetType() != token.LPAREN {
		return nil, fmt.Errorf("invalid token type, expected left parenthesis, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	p.curPos++
	for p.tokens[p.curPos].GetType() != token.RPAREN {
		if p.tokens[p.curPos].GetType() == token.COMMA {

		} else if p.tokens[p.curPos].GetType() == token.IDENTIFIER {
			function.Parameters = append(function.Parameters, ast.NewIdentifierExpression(p.tokens[p.curPos].GetLiteral()))
		} else {
			return nil, fmt.Errorf("unexpected token type: %v, expected identifier", p.tokens[p.curPos].GetLiteral())
		}
		p.curPos++
	}
	p.curPos++
	if p.tokens[p.curPos].GetType() != token.LBRACE {
		return nil, fmt.Errorf("invalid token type, expected left brace, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	p.curPos++
	function.Body = ast.NewBlockStatement(nil)
	for p.tokens[p.curPos].GetType() != token.RBRACE {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, fmt.Errorf("parse statement error: %v", err)
		}
		function.Body.Statements = append(function.Body.Statements, statement)
	}
	p.curPos++
	return function, nil
}

func (p *Parser) parseIdentifierExpression() (ast.Expression, error) {
	if p.tokens[p.curPos].GetType() != token.IDENTIFIER {
		return nil, fmt.Errorf("invalid token type, expected identifier, but got %v", p.tokens[p.curPos].GetLiteral())
	}
	identifierExpression := ast.NewIdentifierExpression(p.tokens[p.curPos].GetLiteral())
	p.curPos++
	return identifierExpression, nil
}

func (p *Parser) parseComplexExpression() (ast.Expression, error) {
	complexExpression := ast.NewComplexExpression(nil, token.New(token.EOF, ""), nil)
	switch p.tokens[p.curPos].GetType() {
	case token.IDENTIFIER:
		exp, err := p.parseIdentifierExpression()
		if err != nil {
			return nil, fmt.Errorf("parse identifier expression error: %v", err)
		}
		complexExpression.Left = exp
	case token.INT:
		exp, err := p.parseLiteralExpression()
		if err != nil {
			return nil, fmt.Errorf("parse literal expression error: %v", err)
		}
		complexExpression.Left = exp
	}
	complexExpression.Operator = p.tokens[p.curPos]
	p.curPos++
	expression, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("parse expression statement error: %v", err)
	}
	complexExpression.Right = expression
	return complexExpression, nil
}

func (p *Parser) parseLiteralExpression() (ast.Expression, error) {
	literalExpression := ast.NewLiteralExpression(0)
	literal := p.tokens[p.curPos].GetLiteral()
	if v, err := strconv.ParseInt(literal, 10, 64); err != nil {
		return nil, fmt.Errorf("parse literal expression error, %v", err)
	} else {
		literalExpression.Value = int(v)
	}
	p.curPos++
	return literalExpression, nil
}

func (p *Parser) peekToken() token.Token {
	if p.curPos >= len(p.tokens)-1 {
		return token.New(token.EOF, "")
	}
	return p.tokens[p.curPos+1]
}
