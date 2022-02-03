package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type AstPrinter struct{}

func (v *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return v.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (v *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return v.parenthesize("group", expr.expression)
}

func (v *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	switch val := expr.value.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)
	case nil:
		return "Nil"
	default:
		panic(errorMsg(-1, "VisitLiteralExpr", fmt.Sprintf("unsupported type %v", reflect.TypeOf(val))))
	}
}

func (v *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return v.parenthesize(expr.operator.lexeme, expr.right)
}

func (v *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	builder := strings.Builder{}
	builder.WriteString("(")
	builder.WriteString(name)
	for _, e := range exprs {
		builder.WriteString(" ")
		builder.WriteString(e.Accept(v).(string))
	}
	builder.WriteString(")")
	return builder.String()
}
