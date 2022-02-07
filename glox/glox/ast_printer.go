package glox

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(expr Expr) string {
	return expr.accept(a).(string)
}

func (a *AstPrinter) visitAssign(expr *Assign) interface{} {
	return a.parenthesize2("=", expr.name.lexeme, expr.value)
}

func (a *AstPrinter) visitBinary(expr *Binary) interface{} {
	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (a *AstPrinter) visitCall(expr *Call) interface{} {
	return a.parenthesize2("call", expr.callee, expr.arguments)
}

func (a *AstPrinter) visitGet(expr *Get) interface{} {
	return a.parenthesize2(".", expr.object, expr.name.lexeme)
}

func (a *AstPrinter) visitGrouping(expr *Grouping) interface{} {
	return a.parenthesize("group", expr.expression)
}

func (a *AstPrinter) visitLiteral(expr *Literal) interface{} {
	return fmt.Sprintf("%v", expr.value)
}

func (a *AstPrinter) visitLogical(expr *Logical) interface{} {
	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (a *AstPrinter) visitSet(expr *Set) interface{} {
	return a.parenthesize2("=", expr.object, expr.name.lexeme, expr.value)
}

func (a *AstPrinter) visitSuper(expr *Super) interface{} {
	return a.parenthesize2("super", expr.method)
}

func (a *AstPrinter) visitThis(expr *This) interface{} {
	return a.parenthesize2("this", expr.keyword)
}

func (a *AstPrinter) visitUnary(expr *Unary) interface{} {
	return a.parenthesize(expr.operator.lexeme, expr.right)
}

func (a *AstPrinter) visitVariable(expr *Variable) interface{} {
	return expr.name.lexeme
}

func (a *AstPrinter) parenthesize(name string, parts ...Expr) string {
	builder := strings.Builder{}
	builder.WriteRune('(')
	builder.WriteString(name)
	for _, part := range parts {
		builder.WriteRune(' ')
		builder.WriteString(part.accept(a).(string))
	}
	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) parenthesize2(name string, parts ...interface{}) string {
	builder := strings.Builder{}
	builder.WriteRune('(')
	builder.WriteString(name)
	a.transform(&builder, parts)
	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) transform(builder *strings.Builder, parts ...interface{}) {
	for _, part := range parts {
		switch p := part.(type) {
		case Expr:
			builder.WriteString(p.accept(a).(string))
		case Token:
			builder.WriteString(p.lexeme)
		case []Expr:
			a.transform(builder, p)
		default:
			builder.WriteString(fmt.Sprintf("%v", p))
		}
	}
}
