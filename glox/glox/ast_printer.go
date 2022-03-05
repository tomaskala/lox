package glox

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (a *AstPrinter) PrintExpr(expr Expr) string {
	return expr.accept(a).(string)
}

func (a *AstPrinter) visitAssign(expr Assign) interface{} {
	return a.parenthesize2("=", expr.name.lexeme, expr.value)
}

func (a *AstPrinter) visitBinary(expr Binary) interface{} {
	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (a *AstPrinter) visitCall(expr Call) interface{} {
	return a.parenthesize2("call", expr.callee, expr.arguments)
}

func (a *AstPrinter) visitGet(expr Get) interface{} {
	return a.parenthesize2(".", expr.object, expr.name.lexeme)
}

func (a *AstPrinter) visitGrouping(expr Grouping) interface{} {
	return a.parenthesize("group", expr.expression)
}

func (a *AstPrinter) visitLiteral(expr Literal) interface{} {
	return fmt.Sprintf("%v", expr.value)
}

func (a *AstPrinter) visitLogical(expr Logical) interface{} {
	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (a *AstPrinter) visitSet(expr Set) interface{} {
	return a.parenthesize2("=", expr.object, expr.name.lexeme, expr.value)
}

func (a *AstPrinter) visitSuper(expr Super) interface{} {
	return a.parenthesize2("super", expr.method)
}

func (a *AstPrinter) visitThis(expr This) interface{} {
	return a.parenthesize2("this", expr.keyword)
}

func (a *AstPrinter) visitUnary(expr Unary) interface{} {
	return a.parenthesize(expr.operator.lexeme, expr.right)
}

func (a *AstPrinter) visitVariable(expr Variable) interface{} {
	return expr.name.lexeme
}

func (a *AstPrinter) PrintStmt(stmt Stmt) string {
	return stmt.accept(a).(string)
}

func (a *AstPrinter) visitBlock(stmt Block) interface{} {
	builder := strings.Builder{}
	builder.WriteString("(block ")
	for _, statement := range stmt.statements {
		builder.WriteString(statement.accept(a).(string))
	}
	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) visitClass(stmt Class) interface{} {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("(class %s", stmt.name.lexeme))
	if stmt.superclass != nil {
		builder.WriteString(fmt.Sprintf(" < %s", a.PrintExpr(stmt.superclass)))
	}
	for _, method := range stmt.methods {
		builder.WriteString(fmt.Sprintf(" %s", a.PrintStmt(method)))
	}
	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) visitExpression(stmt Expression) interface{} {
	return a.parenthesize(";", stmt.expression)
}

func (a *AstPrinter) visitFunction(stmt Function) interface{} {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("(fun %s(", stmt.name.lexeme))
	for _, param := range stmt.params {
		if param != stmt.params[0] {
			builder.WriteRune(' ')
		}
		builder.WriteString(param.lexeme)
	}
	builder.WriteString(") ")
	for _, bodyStmt := range stmt.body {
		builder.WriteString(a.PrintStmt(bodyStmt))
	}
	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) visitIf(stmt If) interface{} {
	if stmt.elseBranch == nil {
		return a.parenthesize2("if", stmt.condition, stmt.thenBranch)
	} else {
		return a.parenthesize2("if-else", stmt.condition, stmt.thenBranch, stmt.elseBranch)
	}
}

func (a *AstPrinter) visitPrint(stmt Print) interface{} {
	return a.parenthesize("print", stmt.expression)
}

func (a *AstPrinter) visitReturn(stmt Return) interface{} {
	if stmt.value == nil {
		return "(return)"
	} else {
		return a.parenthesize("return", stmt.value)
	}
}

func (a *AstPrinter) visitVar(stmt Var) interface{} {
	if stmt.initializer == nil {
		return a.parenthesize2("var", stmt.name)
	} else {
		return a.parenthesize2("var", stmt.name, "=", stmt.initializer)
	}
}

func (a *AstPrinter) visitWhile(stmt While) interface{} {
	builder := strings.Builder{}
	builder.WriteString("(while ")
	builder.WriteString(stmt.condition.accept(a).(string))
	builder.WriteRune(' ')
	builder.WriteString(stmt.body.accept(a).(string))
	builder.WriteRune(')')
	return builder.String()
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
		builder.WriteRune(' ')
		switch p := part.(type) {
		case Expr:
			builder.WriteString(p.accept(a).(string))
		case Stmt:
			builder.WriteString(p.accept(a).(string))
		case Token:
			builder.WriteString(p.lexeme)
		case []interface{}:
			a.transform(builder, p...)
		case []Expr:
			cast := make([]interface{}, len(p))
			for i := range p {
				cast[i] = p[i]
			}
			a.transform(builder, cast)
		default:
			builder.WriteString(fmt.Sprintf("%v", p))
		}
	}
}
