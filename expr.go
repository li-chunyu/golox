package main

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

type Expr interface {
	Accept(v ExprVisitor) interface{}
}

type Binary struct {
	left     Expr
	operator *Token
	right    Expr
}

func NewBinary(left Expr, operator *Token, right Expr) Expr {

	expr := &Binary{
		left:     left,
		operator: operator,
		right:    right,
	}
	return expr
}

func (expr *Binary) Accept(v ExprVisitor) interface{} {
	return v.VisitBinaryExpr(expr)
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) Expr {

	expr := &Grouping{
		expression: expression,
	}
	return expr
}

func (expr *Grouping) Accept(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(expr)
}

type Literal struct {
	value interface{}
}

func NewLiteral(value interface{}) Expr {

	expr := &Literal{
		value: value,
	}
	return expr
}

func (expr *Literal) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteralExpr(expr)
}

type Unary struct {
	operator *Token
	right    Expr
}

func NewUnary(operator *Token, right Expr) Expr {

	expr := &Unary{
		operator: operator,
		right:    right,
	}
	return expr
}

func (expr *Unary) Accept(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(expr)
}
