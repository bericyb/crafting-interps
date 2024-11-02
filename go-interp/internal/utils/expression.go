package utils

import "fmt"

type AstPrinter struct{}

func (AstP *AstPrinter) Print(expr Expr) string {
	return expr.accept(AstP)
}

type ExprVisitor interface {
	visitBinary(expr Binary) string
	visitGrouping(expr Grouping) string
	visitLiteral(expr Literal) string
	visitUnary(expr Unary) string
}

type Expr interface {
	accept(visitor ExprVisitor) string
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (B Binary) accept(visitor ExprVisitor) string {
	return visitor.visitBinary(B)
}

type Grouping struct {
	Expression Expr
}

func (G Grouping) accept(visitor ExprVisitor) string {
	return visitor.visitGrouping(G)
}

type Literal struct {
	Value any
}

func (L Literal) accept(visitor ExprVisitor) string {
	return visitor.visitLiteral(L)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (U Unary) accept(visitor ExprVisitor) string {
	return visitor.visitUnary(U)
}

func (AstP *AstPrinter) visitBinary(expr Binary) string {
	return AstP.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (AstP *AstPrinter) visitGrouping(expr Grouping) string {
	return AstP.parenthesize("group", expr.Expression)
}

func (AstP *AstPrinter) visitLiteral(expr Literal) string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%#v", expr.Value)
}

func (AstP *AstPrinter) visitUnary(expr Unary) string {
	return AstP.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (AstP *AstPrinter) parenthesize(name string, exprs ...Expr) string {

	s := ""
	s = s + "(" + name

	for _, expr := range exprs {
		s = s + " "
		s = s + expr.accept(AstP)

	}

	s = s + ")"
	return s
}
