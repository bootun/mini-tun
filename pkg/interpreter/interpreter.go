package interpreter

import (
	"fmt"
	"maps"

	"github.com/bootun/mini-tun/pkg/ast"
	"github.com/bootun/mini-tun/pkg/token"
)

type Interpreter struct {
	stack   functionStack
	program ast.Program
}

func NewInterpreter(program ast.Program) *Interpreter {
	return &Interpreter{
		stack: functionStack{
			envs: make(map[string]interface{}),
		},
		program: program,
	}
}

func (i *Interpreter) Exec() error {
	for _, statement := range i.program.Statements {
		switch statement.(type) {
		case *ast.VariableAssignment:
			variableAssignment := statement.(*ast.VariableAssignment)
			// 判断value是函数还是值
			i.stack.envs[variableAssignment.VariableName] = i.stack.computeExpression(variableAssignment.Value)
		}
	}
	// 运行最终态
	for k, v := range i.stack.envs {
		switch v.(type) {
		case int:
			fmt.Printf("%s = %d\n", k, v.(int))
		case *ast.FunctionLiteral:
			fmt.Printf("%s = %s\n", k, v.(*ast.FunctionLiteral).TokenLiteral())
		}
	}
	return nil
}

func (s *functionStack) computeExpression(expression ast.Expression) interface{} {
	switch expression.(type) {
	case *ast.LiteralExpression:
		return expression.(*ast.LiteralExpression).Value
	case *ast.ComplexExpression:
		node := expression.(*ast.ComplexExpression)
		left := s.computeExpression(node.Left).(int)
		right := s.computeExpression(node.Right).(int)
		switch expression.(*ast.ComplexExpression).Operator.Type {
		case token.PLUS:
			return left + right
		case token.MINUS:
			return left - right
		default:
			panic("unsupported operator")
		}
	case *ast.IdentifierExpression:
		node := expression.(*ast.IdentifierExpression)
		return s.envs[node.Value].(int)
	case *ast.FunctionCall:
		// 函数调用
		node := expression.(*ast.FunctionCall)
		funcDecl := s.envs[node.FunctionName].(*ast.FunctionLiteral)
		params := make(map[string]interface{})
		for i, param := range node.Arguments {
			name := funcDecl.Parameters[i].Value
			value := s.computeExpression(param).(int)
			params[name] = value
		}
		// 初始化函数栈
		envs := maps.Clone(s.envs)
		callStack := functionStack{envs: envs}
		return callStack.computeFunction(funcDecl)

	case *ast.FunctionLiteral:
		// 函数定义
		node := expression.(*ast.FunctionLiteral)
		return node
	}
	return 0
}

type functionStack struct {
	envs map[string]interface{}
}

func (s *functionStack) computeFunction(function *ast.FunctionLiteral) interface{} {
	if function.Body == nil {
		return nil
	}
	for _, statement := range function.Body.Statements {
		switch statement.(type) {
		case *ast.VariableAssignment:
			variableAssignment := statement.(*ast.VariableAssignment)
			// 判断value是函数还是值
			s.envs[variableAssignment.VariableName] = s.computeExpression(variableAssignment.Value)
		case *ast.ReturnStatement:
			return s.computeExpression(statement.(*ast.ReturnStatement).ReturnValue)
		}
	}
	return nil
}
