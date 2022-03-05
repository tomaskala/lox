package glox

type StmtVisitor interface {
	visitBlock(stmt Block) interface{}
	visitClass(stmt Class) interface{}
	visitExpression(stmt Expression) interface{}
	visitFunction(stmt Function) interface{}
	visitIf(stmt If) interface{}
	visitPrint(stmt Print) interface{}
	visitReturn(stmt Return) interface{}
	visitVar(stmt Var) interface{}
	visitWhile(stmt While) interface{}
}

type Stmt interface {
	accept(v StmtVisitor) interface{}
}

type Block struct {
	statements []Stmt
}

func (stmt Block) accept(v StmtVisitor) interface{} {
	return v.visitBlock(stmt)
}

type Class struct {
	name       Token
	superclass *Variable
	methods    []Function
}

func (stmt Class) accept(v StmtVisitor) interface{} {
	return v.visitClass(stmt)
}

type Expression struct {
	expression Expr
}

func (stmt Expression) accept(v StmtVisitor) interface{} {
	return v.visitExpression(stmt)
}

type Function struct {
	name   Token
	params []Token
	body   []Stmt
}

func (stmt Function) accept(v StmtVisitor) interface{} {
	return v.visitFunction(stmt)
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (stmt If) accept(v StmtVisitor) interface{} {
	return v.visitIf(stmt)
}

type Print struct {
	expression Expr
}

func (stmt Print) accept(v StmtVisitor) interface{} {
	return v.visitPrint(stmt)
}

type Return struct {
	keyword Token
	value   Expr
}

func (stmt Return) accept(v StmtVisitor) interface{} {
	return v.visitReturn(stmt)
}

type Var struct {
	name        Token
	initializer Expr
}

func (stmt Var) accept(v StmtVisitor) interface{} {
	return v.visitVar(stmt)
}

type While struct {
	condition Expr
	body      Stmt
}

func (stmt While) accept(v StmtVisitor) interface{} {
	return v.visitWhile(stmt)
}
