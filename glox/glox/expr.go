package glox

type ExprVisitor interface {
	visitAssign(expr *Assign) interface{}
	visitBinary(expr *Binary) interface{}
	visitCall(expr *Call) interface{}
	visitGet(expr *Get) interface{}
	visitGrouping(expr *Grouping) interface{}
	visitLiteral(expr *Literal) interface{}
	visitLogical(expr *Logical) interface{}
	visitSet(expr *Set) interface{}
	visitSuper(expr *Super) interface{}
	visitThis(expr *This) interface{}
	visitUnary(expr *Unary) interface{}
	visitVariable(expr *Variable) interface{}
}

type Expr interface {
	accept(v ExprVisitor) interface{}
}

type Assign struct {
	Expr
	name  Token
	value Expr
}

func (expr *Assign) accept(v ExprVisitor) interface{} {
	return v.visitAssign(expr)
}

type Binary struct {
	Expr
	left     Expr
	operator Token
	right    Expr
}

func (expr *Binary) accept(v ExprVisitor) interface{} {
	return v.visitBinary(expr)
}

type Call struct {
	Expr
	callee    Expr
	paren     Token
	arguments []Expr
}

func (expr *Call) accept(v ExprVisitor) interface{} {
	return v.visitCall(expr)
}

type Get struct {
	Expr
	object Expr
	name   Token
}

func (expr *Get) accept(v ExprVisitor) interface{} {
	return v.visitGet(expr)
}

type Grouping struct {
	Expr
	expression Expr
}

func (expr *Grouping) accept(v ExprVisitor) interface{} {
	return v.visitGrouping(expr)
}

type Literal struct {
	Expr
	value interface{}
}

func (expr *Literal) accept(v ExprVisitor) interface{} {
	return v.visitLiteral(expr)
}

type Logical struct {
	Expr
	left     Expr
	operator Token
	right    Expr
}

func (expr *Logical) accept(v ExprVisitor) interface{} {
	return v.visitLogical(expr)
}

type Set struct {
	Expr
	object Expr
	name   Token
	value  Expr
}

func (expr *Set) accept(v ExprVisitor) interface{} {
	return v.visitSet(expr)
}

type Super struct {
	Expr
	keyword Token
	method  Token
}

func (expr *Super) accept(v ExprVisitor) interface{} {
	return v.visitSuper(expr)
}

type This struct {
	Expr
	keyword Token
}

func (expr *This) accept(v ExprVisitor) interface{} {
	return v.visitThis(expr)
}

type Unary struct {
	Expr
	operator Token
	right    Expr
}

func (expr *Unary) accept(v ExprVisitor) interface{} {
	return v.visitUnary(expr)
}

type Variable struct {
	Expr
	name Token
}

func (expr *Variable) accept(v ExprVisitor) interface{} {
	return v.visitVariable(expr)
}
