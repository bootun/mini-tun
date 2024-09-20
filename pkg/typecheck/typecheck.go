package typecheck

import (
	"fmt"

	"github.com/bootun/mini-tun/pkg/ast"
)

type Checker struct {
	envs    map[string]interface{}
	program ast.Program
}

type blockEnv struct {
	envs map[string]interface{}
}

func NewChecker(program ast.Program) *Checker {
	return &Checker{
		envs:    make(map[string]interface{}),
		program: program,
	}
}

func (c *Checker) Check() error {
	for _, stmt := range c.program.Statements {
		refs, err := getStatementIdentifierReference(stmt)
		if err == nil {
			for _, ref := range refs.Refs {
				if _, ok := c.envs[ref]; !ok {
					return fmt.Errorf("undefined variable: %s", ref)
				}
			}
		} else {
			return fmt.Errorf("get statement identifier reference error: %v", err)
		}
		c.envs[refs.VariableName] = struct{}{}
	}
	return nil
}

type RefInfo struct {
	VariableName string
	Refs         []string
}

func getStatementIdentifierReference(stmt ast.Statement) (*RefInfo, error) {
	switch stmt.(type) {
	case *ast.VariableAssignment:
		node := stmt.(*ast.VariableAssignment)
		refs, err := getExpressionIdentifierReference(node.Value)
		if err != nil {
			return nil, fmt.Errorf("get expression identifier reference error: %v", err)
		}
		return &RefInfo{
			VariableName: node.VariableName,
			Refs:         refs,
		}, nil
	case *ast.ReturnStatement:
		node := stmt.(*ast.ReturnStatement)
		refs, err := getExpressionIdentifierReference(node.ReturnValue)
		if err != nil {
			return nil, fmt.Errorf("get expression identifier reference from return statement error: %v", err)
		}
		return &RefInfo{
			VariableName: "",
			Refs:         refs,
		}, nil
	// case *ast.BlockStatement:
	// 	node := stmt.(*ast.BlockStatement)
	default:
		return nil, fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

func getExpressionIdentifierReference(expr ast.Expression) ([]string, error) {
	switch expr.(type) {
	case *ast.IdentifierExpression:
		node := expr.(*ast.IdentifierExpression)
		return []string{node.Value}, nil
	case *ast.ComplexExpression:
		var refs []string
		node := expr.(*ast.ComplexExpression)
		leftRefs, err := getExpressionIdentifierReference(node.Left)
		if err != nil {
			return nil, fmt.Errorf("get expression identifier reference error: %v", err)
		}
		rightRefs, err := getExpressionIdentifierReference(node.Right)
		if err != nil {
			return nil, fmt.Errorf("get expression identifier reference error: %v", err)
		}
		refs = append(refs, leftRefs...)
		refs = append(refs, rightRefs...)
		return refs, nil
	case *ast.LiteralExpression:
		return []string{}, nil
	case *ast.FunctionLiteral:
		node := expr.(*ast.FunctionLiteral)
		externalRefs, err := parseBlockIdentifierReference(node.Body)
		if err != nil {
			return nil, fmt.Errorf("parse function body error: %v", err)
		}
		parameters := make(map[string]struct{})
		for _, param := range node.Parameters {
			parameters[param.Value] = struct{}{}
		}
		for _, ref := range externalRefs {
			if _, ok := parameters[ref]; !ok {
				return nil, fmt.Errorf("undefined variable: %s", ref)
			}
		}

		return []string{}, nil
	case *ast.FunctionCall:
		node := expr.(*ast.FunctionCall)
		var refs []string
		refs = append(refs, node.FunctionName)
		for _, arg := range node.Arguments {
			argRefs, err := getExpressionIdentifierReference(arg)
			if err != nil {
				return nil, fmt.Errorf("get expression identifier reference error: %v", err)
			}
			refs = append(refs, argRefs...)
		}
		return refs, nil
	}
	return nil, nil
}

// block只会在下级作用域增加变量，不会给上级作用域增加变量
func parseBlockIdentifierReference(block *ast.BlockStatement) ([]string, error) {
	envs := make(map[string]interface{})
	var externalRefs []string
	for _, stmt := range block.Statements {
		refs, err := getStatementIdentifierReference(stmt)
		if err == nil {
			for _, ref := range refs.Refs {
				if _, ok := envs[ref]; !ok {
					externalRefs = append(externalRefs, ref)
				}
			}
		} else {
			return []string{}, fmt.Errorf("get statement identifier reference error: %v", err)
		}
		envs[refs.VariableName] = struct{}{}
	}
	return externalRefs, nil
}
