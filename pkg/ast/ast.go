package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bootun/mini-tun/pkg/token"
)

const (
	NodeTypeExpression = "Expression"
	NodeTypeStatement  = "Statement"
)

type NodeInfo struct {
	NodeType string // 节点类型
	NodeName string // 节点实际名称
}

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Program struct {
	Statements []Statement
}

func (p *Program) JSON() string {
	tree, _ := json.MarshalIndent(p, "", "    ")
	return string(tree)
}

type VariableAssignment struct {
	NodeInfo     NodeInfo
	VariableName string
	Value        Expression
}

func NewVariableAssignment(variableName string, value Expression) *VariableAssignment {
	return &VariableAssignment{
		NodeInfo: NodeInfo{
			NodeType: NodeTypeStatement,
			NodeName: "VariableAssignment",
		},
		VariableName: variableName,
		Value:        value,
	}
}

func (v *VariableAssignment) TokenLiteral() string {
	return fmt.Sprintf("let %s = %s", v.VariableName, v.Value.TokenLiteral())
}

func (v *VariableAssignment) Name() string {
	return "VariableAssignment"
}

// 字面值常量表达式
type LiteralExpression struct {
	NodeInfo NodeInfo
	Value    int
}

func NewLiteralExpression(value int) *LiteralExpression {
	return &LiteralExpression{
		NodeInfo: NodeInfo{
			NodeType: NodeTypeExpression,
			NodeName: "LiteralExpression",
		},
		Value: value,
	}
}

func (l *LiteralExpression) TokenLiteral() string {
	return fmt.Sprintf("%d", l.Value)
}

type FunctionCall struct {
	NodeInfo     NodeInfo
	FunctionName string
	Arguments    []Expression
}

func NewFunctionCall(functionName string, arguments []Expression) *FunctionCall {
	return &FunctionCall{
		NodeInfo: NodeInfo{
			NodeType: NodeTypeExpression,
			NodeName: "FunctionCall",
		},
		FunctionName: functionName,
		Arguments:    arguments,
	}
}

func (f *FunctionCall) TokenLiteral() string {
	return fmt.Sprintf("%s(%s)", f.FunctionName, f.Arguments)
}

type FunctionLiteral struct {
	NodeInfo   NodeInfo
	Parameters []*IdentifierExpression
	Body       *BlockStatement
}

func NewFunctionLiteral(parameters []*IdentifierExpression, body *BlockStatement) *FunctionLiteral {
	return &FunctionLiteral{
		NodeInfo: NodeInfo{
			NodeType: NodeTypeExpression,
			NodeName: "FunctionLiteral",
		},
		Parameters: parameters,
		Body:       body,
	}
}

func (f *FunctionLiteral) TokenLiteral() string {
	var buf strings.Builder
	buf.WriteString("function(")
	var params []string
	for _, param := range f.Parameters {
		params = append(params, param.TokenLiteral())
	}
	buf.WriteString(strings.Join(params, ","))
	buf.WriteString(") {")
	for _, stmt := range f.Body.Statements {
		buf.WriteString(stmt.TokenLiteral())
		buf.WriteString(";")
	}
	buf.WriteString("}")
	return buf.String()
}

type BlockStatement struct {
	NodeInfo   NodeInfo
	Statements []Statement
}

func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{
		NodeInfo: NodeInfo{
			NodeType: NodeTypeStatement,
			NodeName: "BlockStatement",
		},
		Statements: statements,
	}
}

func (b *BlockStatement) TokenLiteral() string {
	return fmt.Sprintf("{%s}", b.Statements)
}

type ReturnStatement struct {
	NodeInfo    NodeInfo
	ReturnValue Expression
}

func NewReturnStatement(returnValue Expression) *ReturnStatement {
	return &ReturnStatement{
		NodeInfo: NodeInfo{
			NodeType: NodeTypeStatement,
			NodeName: "ReturnStatement",
		},
		ReturnValue: returnValue,
	}
}

func (r *ReturnStatement) TokenLiteral() string {
	return fmt.Sprintf("return %s", r.ReturnValue.TokenLiteral())
}

type ComplexExpression struct {
	NodeInfo NodeInfo
	Left     Expression
	Operator token.Token
	Right    Expression
}

func NewComplexExpression(left Expression, operator token.Token, right Expression) *ComplexExpression {
	return &ComplexExpression{
		NodeInfo: NodeInfo{
			NodeType: NodeTypeExpression,
			NodeName: "ComplexExpression",
		},
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (c *ComplexExpression) TokenLiteral() string {
	return fmt.Sprintf("%s %s %s", c.Left.TokenLiteral(), c.Operator.Literal, c.Right.TokenLiteral())
}

type IdentifierExpression struct {
	NodeInfo NodeInfo
	Value    string // 标识符名称
}

func NewIdentifierExpression(value string) *IdentifierExpression {
	return &IdentifierExpression{
		NodeInfo: NodeInfo{
			NodeType: NodeTypeExpression,
			NodeName: "IdentifierExpression",
		},
		Value: value,
	}
}

func (i *IdentifierExpression) TokenLiteral() string {
	return i.Value
}
