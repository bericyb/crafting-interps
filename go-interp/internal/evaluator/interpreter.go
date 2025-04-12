package evaluator

import (
	"fmt"
	"go_interp/internal/utils"
	"strconv"

	"github.com/tdewolff/parse/v2/strconv"
)

type Interpreter struct {
}

func (i *Interpreter) Interpret(expr utils.Expr) {
	value := i.evaluate(expr)
	fmt.Println(value)
	return
}

type InterpreterVisitor interface {
	Accept(expr utils.Expr) any
}

func (i *Interpreter) VisitLiteralExpr(expr utils.Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr utils.Grouping) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) visitUnaryExpr(expr utils.Unary) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case utils.MINUS:
		return -right.(float64)
	case utils.BANG:
		return !isTruthy(right)
	}
	return nil
}

func (i *Interpreter) evaluate(expr utils.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(expr utils.Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case utils.MINUS:
		return left.(float64) - right.(float64)
	case utils.SLASH:
		return left.(float64) / right.(float64)
	case utils.STAR:
		return left.(float64) * right.(float64)
	case utils.PLUS:
		lres, lok := left.(string)
		rres, rok := right.(string)
		if lok && rok {
			return lres + rres
		} else {
			lres, lok := left.(float64)
			rres, rok := right.(float64)
			if lok && rok {
				return lres + rres
			}
		}
	case utils.GREATER:
		return left.(float64) > right.(float64)
	case utils.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case utils.LESS:
		return left.(float64) < right.(float64)
	case utils.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case utils.BANG_EQUAL:
		return !isEqual(left, right)
	case utils.EQUAL_EQUAL:
		return isEqual(left, right)
	}
}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if object, ok := object.(bool); ok {
		return object
	}
	return false
}

func stringify(obj any) string {
	if obj == nil {
		return "nil"
	}

	if f, ok := obj.(float64); ok {
		text := strconv.FormatFloat(f, 'f', 1, 64)
		return text
	}
	return obj.(string)
}
