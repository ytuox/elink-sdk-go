package util

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// EvaluateExpression 解析并计算表达式
// 例子：offset := "(1*6)+1994"
// 结果：2000
func EvaluateExpression(expr string) (int, error) {
	exprAST, err := parser.ParseExpr(expr)
	if err != nil {
		return 0, fmt.Errorf("解析表达式出错: %v", err)
	}

	result, err := evalAST(exprAST)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// evalAST 递归计算 AST 表达式
func evalAST(node ast.Node) (int, error) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		left, err := evalAST(n.X)
		if err != nil {
			return 0, err
		}
		right, err := evalAST(n.Y)
		if err != nil {
			return 0, err
		}
		return applyOperator(left, right, n.Op)
	case *ast.BasicLit:
		val, err := strconv.Atoi(n.Value)
		if err != nil {
			return 0, err
		}
		return val, nil
	case *ast.ParenExpr:
		return evalAST(n.X)
	default:
		return 0, fmt.Errorf("不支持的表达式类型: %T", node)
	}
}

// applyOperator 应用操作符进行计算
func applyOperator(a, b int, op token.Token) (int, error) {
	switch op {
	case token.ADD:
		return a + b, nil
	case token.SUB:
		return a - b, nil
	case token.MUL:
		return a * b, nil
	case token.QUO:
		if b == 0 {
			return 0, fmt.Errorf("除数不能为0")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("不支持的操作符: %v", op)
	}
}
