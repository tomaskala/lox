package glox

type Resolver struct {
	interpreter *Interpreter
	scopes      []map[string]bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      nil,
	}
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.accept(r)
}

func (r *Resolver) visitAssign(expr Assign) interface{} {
	r.resolveExpression(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) visitBinary(expr Binary) interface{} {
	r.resolveExpression(expr.left)
	r.resolveExpression(expr.right)
	return nil
}

func (r *Resolver) visitCall(expr Call) interface{} {
	r.resolveExpression(expr.callee)
	for _, argument := range expr.arguments {
		r.resolveExpression(argument)
	}
	return nil
}

func (r *Resolver) visitGet(expr Get) interface{} {
	return nil
}

func (r *Resolver) visitGrouping(expr Grouping) interface{} {
	r.resolveExpression(expr.expression)
	return nil
}

func (r *Resolver) visitLiteral(expr Literal) interface{} {
	return nil
}

func (r *Resolver) visitLogical(expr Logical) interface{} {
	r.resolveExpression(expr.left)
	r.resolveExpression(expr.right)
	return nil
}

func (r *Resolver) visitSet(expr Set) interface{} {
	return nil
}

func (r *Resolver) visitSuper(expr Super) interface{} {
	return nil
}

func (r *Resolver) visitThis(expr This) interface{} {
	return nil
}

func (r *Resolver) visitUnary(expr Unary) interface{} {
	r.resolveExpression(expr.right)
	return nil
}

func (r *Resolver) visitVariable(expr Variable) interface{} {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		if initialized, ok := scope[expr.name.lexeme]; ok && !initialized {
			panic(resolverError{gloxError(expr.name, "Cannot read a local variable in its own initializer.")})
		}
	}
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.accept(r)
}

func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, stmt := range statements {
		r.resolveStatement(stmt)
	}
}

func (r *Resolver) visitBlock(stmt Block) interface{} {
	r.beginScope()
	r.resolveStatements(stmt.statements)
	r.endScope()
	return nil
}

func (r *Resolver) visitClass(stmt Class) interface{} {
	return nil
}

func (r *Resolver) visitExpression(stmt Expression) interface{} {
	r.resolveExpression(stmt.expression)
	return nil
}

func (r *Resolver) visitFunction(stmt Function) interface{} {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt)
	return nil
}

func (r *Resolver) visitIf(stmt If) interface{} {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStatement(stmt.elseBranch)
	}
	return nil
}

func (r *Resolver) visitPrint(stmt Print) interface{} {
	r.resolveExpression(stmt.expression)
	return nil
}

func (r *Resolver) visitReturn(stmt Return) interface{} {
	if stmt.value != nil {
		r.resolveExpression(stmt.value)
	}
	return nil
}

func (r *Resolver) visitVar(stmt Var) interface{} {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpression(stmt.initializer)
	}
	r.define(stmt.name)
	return nil
}

func (r *Resolver) visitWhile(stmt While) interface{} {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.body)
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes = r.scopes[0 : len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function Function) {
	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.body)
	r.endScope()
}