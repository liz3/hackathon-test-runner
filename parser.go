package main

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
)

func getIdentifier(node ast.Expression) string {
	v, ok := node.(*ast.Identifier)
	if !ok {
		return ""
	}
	return v.Name
}
func getStringValue(node ast.Expression) string {
	v, ok := node.(*ast.StringLiteral)
	if !ok {
		return ""
	}
	return v.Value
}

func ExtractTests(suite *Suite, ctx *Context, node ast.Statement) {
	block, isBlock := (node).(*ast.BlockStatement)
	if isBlock {
		for _, raw := range block.List {
			exp, isExp := raw.(*ast.ExpressionStatement)
			if isExp {
				callExpression, isCallExpression := (exp.Expression).(*ast.CallExpression)
				if isCallExpression {
					var funcName = getIdentifier(callExpression.Callee)
					if len(callExpression.ArgumentList) != 2 {
						continue
					}
					if funcName == "it" {
						// this should be a test
						ctx.CreateTest(suite, getStringValue(callExpression.ArgumentList[0]))
					} else if funcName == "describe" {
						suite := ctx.CreateSuite(getStringValue(callExpression.ArgumentList[0]))
						funcLiteral, isFunc := callExpression.ArgumentList[1].(*ast.FunctionLiteral)
						if isFunc {
							ExtractTests(suite, ctx, funcLiteral.Body)
						}
					}
				}
			}
		}
	}
}
func SuiteFromFunction(file string, ctx *Context, node ast.Statement) {
	block, isBlock := (node).(*ast.BlockStatement)
	if isBlock {
		for _, raw := range block.List {
			exp, isExp := raw.(*ast.ExpressionStatement)
			if isExp {
				callExpression, isCallExpression := (exp.Expression).(*ast.CallExpression)
				if isCallExpression {
					var funcName = getIdentifier(callExpression.Callee)
					if len(callExpression.ArgumentList) != 2 {
						continue
					}
					if funcName == "describe" {
						suite := ctx.CreateSuite(getStringValue(callExpression.ArgumentList[0]))
						funcLiteral, isFunc := callExpression.ArgumentList[1].(*ast.FunctionLiteral)
						if isFunc {
							ExtractTests(suite, ctx, funcLiteral.Body)
						}
					}
				}
			}
		}
	}
}
func Traverse(statement ast.Statement, file string, ctx *Context) {
	exp, isExp := statement.(*ast.ExpressionStatement)
	if isExp {
		callExpression, isCallExpression := (exp.Expression).(*ast.CallExpression)
		if isCallExpression {
			var funcName = getIdentifier(callExpression.Callee)

			if funcName == "describe" {
				// this is a test suite
				suite := ctx.CreateSuite(getStringValue(callExpression.ArgumentList[0]))
				funcLiteral, isFunc := callExpression.ArgumentList[1].(*ast.FunctionLiteral)
				if isFunc {
					ExtractTests(suite, ctx, funcLiteral.Body)
				}
			} else {
				for _, arg := range callExpression.ArgumentList {
					funcLiteral, isFunc := arg.(*ast.FunctionLiteral)
					if isFunc {
						SuiteFromFunction(file, ctx, funcLiteral.Body)
					}
				}
			}
		} else {
			funcLiteral, isFunc := (exp.Expression).(*ast.FunctionLiteral)
			if isFunc {
				SuiteFromFunction(file, ctx, funcLiteral.Body)
			}
		}
	} else {
		v, isVariable := statement.(*ast.VariableStatement)
		if isVariable {
			for _, entry := range v.List {
				variable, isVariable := (entry).(*ast.VariableExpression)
				if isVariable {
					if variable.Initializer != nil {
						funcLiteral, isFunc := (variable.Initializer).(*ast.FunctionLiteral)
						if isFunc {
							SuiteFromFunction(file, ctx, funcLiteral.Body)
						}
					}
				}
			}
		}
	}
}
func Parse(path string, ctx *Context) {
	var res = ReadFile(path)
	if res == nil {
		return
	}
	program, err := parser.ParseFile(nil, res.Name, res.Content, 0)
	if err != nil {
		fmt.Println(path, "failed", err)
		return
	}
	for _, entry := range program.Body {
		Traverse(entry, path, ctx)
	}
}
