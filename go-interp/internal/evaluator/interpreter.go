package evaluator

import (
	"go_interp/internal/utils"
)

type Interpreter struct {
}

func (i Interpreter) VisitLiteralExpr(expr utils.Literal) any {
	return expr.Value
}

func (i Interpreter) VisitGroupingExpr(expr utils.Grouping) utils.Expr {
	return evaluate(expr.Expression)
}

func (o Object) evaluate(expr utils.Expr) {
	return expr.Accept(o)
}
