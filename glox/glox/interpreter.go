package glox

import (
	"fmt"
	"math"
)

type returnSignal struct {
	value interface{}
}

type Interpreter struct {
	globals     *Environment
	environment *Environment
	locals      map[Expr]int
	interactive bool
}

func NewInterpreter(interactive bool) *Interpreter {
	environment := NewGlobals()
	return &Interpreter{
		globals:     environment,
		environment: environment,
		locals:      make(map[Expr]int),
		interactive: interactive,
	}
}

func (i *Interpreter) Interpret(statements []Stmt) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if ie, ok := r.(interpreterError); ok {
				err = ie
			} else {
				panic(r)
			}
		}
	}()
	for _, statement := range statements {
		i.execute(statement)
	}
	return nil
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(i)
}

func (i *Interpreter) visitAssign(expr *Assign) interface{} {
	value := i.evaluate(expr.value)
	if distance, ok := i.locals[expr]; ok {
		i.environment.assignAt(distance, expr.name, value)
	} else {
		i.globals.assign(expr.name, value)
	}
	return value
}

func (i *Interpreter) visitBinary(expr *Binary) interface{} {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)
	switch expr.operator.tokenType {
	case GREATER:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) <= right.(float64)
	case MINUS:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) - right.(float64)
	case BANG_EQUAL:
		return left != right
	case EQUAL_EQUAL:
		return left == right
	case PLUS:
		lf, lok := left.(float64)
		rf, rok := right.(float64)
		if lok && rok {
			return lf + rf
		}
		ls, lok := left.(string)
		rs, rok := right.(string)
		if lok && rok {
			return ls + rs
		}
		panic(interpreterError{gloxError(expr.operator, "Operands must be numbers or strings.")})
	case SLASH:
		checkNumberOperands(expr.operator, left, right)
		rightNum := right.(float64)
		if rightNum == 0.0 {
			panic(interpreterError{gloxError(expr.operator, "Division by zero.")})
		} else {
			return left.(float64) / rightNum
		}
	case STAR:
		checkNumberOperands(expr.operator, left, right)
		return left.(float64) * right.(float64)
	}
	// Unreachable.
	return nil
}

func (i *Interpreter) visitCall(expr *Call) interface{} {
	callee := i.evaluate(expr.callee)
	arguments := make([]interface{}, 0)
	for _, argument := range expr.arguments {
		arguments = append(arguments, i.evaluate(argument))
	}
	checkCallable(callee, expr.paren)
	callable := callee.(Callable)
	if len(arguments) != callable.arity() {
		message := fmt.Sprintf("Expected %d arguments but got %d.", callable.arity, len(arguments))
		panic(interpreterError{gloxError(expr.paren, message)})
	}
	return callable.call(i, arguments)
}

func (i *Interpreter) visitGet(expr *Get) interface{} {
	object := i.evaluate(expr.object)
	if i, ok := object.(*GloxInstance); ok {
		return i.get(expr.name)
	} else {
		panic(interpreterError{gloxError(expr.name, "Only class instances have properties.")})
	}
}

func (i *Interpreter) visitGrouping(expr *Grouping) interface{} {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) visitLiteral(expr *Literal) interface{} {
	return expr.value
}

func (i *Interpreter) visitLogical(expr *Logical) interface{} {
	left := i.evaluate(expr.left)
	if expr.operator.tokenType == OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}
	return i.evaluate(expr.right)
}

func (i *Interpreter) visitSet(expr *Set) interface{} {
	object := i.evaluate(expr.object)
	if g, ok := object.(*GloxInstance); ok {
		value := i.evaluate(expr.value)
		g.set(expr.name, value)
		return value
	} else {
		panic(interpreterError{gloxError(expr.name, "Only class instances have fields.")})
	}
}

func (i *Interpreter) visitSuper(expr *Super) interface{} {
	distance := i.locals[expr]
	superclass := i.environment.getAt(distance, "super").(*GloxClass)
	object := i.environment.getAt(distance-1, "this").(*GloxInstance)
	method := superclass.findMethod(expr.method.lexeme)
	if method == nil {
		message := fmt.Sprintf("Undefined property '%s'.", expr.method.lexeme)
		panic(interpreterError{gloxError(expr.method, message)})
	}
	return method.bind(object)
}

func (i *Interpreter) visitThis(expr *This) interface{} {
	return i.lookupVariable(expr.keyword, expr)
}

func (i *Interpreter) visitUnary(expr *Unary) interface{} {
	right := i.evaluate(expr.right)
	switch expr.operator.tokenType {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		checkNumberOperand(expr.operator, right)
		return -right.(float64)
	}
	// Unreachable.
	return nil
}

func (i *Interpreter) visitVariable(expr *Variable) interface{} {
	return i.lookupVariable(expr.name, expr)
}

func (i *Interpreter) execute(stmt Stmt) interface{} {
	return stmt.accept(i)
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) visitBlock(stmt *Block) interface{} {
	i.executeBlock(stmt.statements, NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) visitClass(stmt *Class) interface{} {
	var superclass *GloxClass
	if stmt.superclass != nil {
		var ok bool
		superclassValue := i.evaluate(stmt.superclass)
		if superclass, ok = superclassValue.(*GloxClass); !ok {
			panic(interpreterError{gloxError(stmt.superclass.name, "A superclass must be a class.")})
		}
	}
	i.environment.define(stmt.name.lexeme, nil)
	if stmt.superclass != nil {
		i.environment = NewEnvironment(i.environment)
		i.environment.define("super", superclass)
	}
	methods := make(map[string]*GloxCallable)
	for _, method := range stmt.methods {
		function := &GloxCallable{
			declaration:   method,
			closure:       i.environment,
			isInitializer: method.name.lexeme == "init",
		}
		methods[method.name.lexeme] = function
	}
	class := &GloxClass{
		name:       stmt.name.lexeme,
		superclass: superclass,
		methods:    methods,
	}
	if superclass != nil {
		i.environment = i.environment.enclosing
	}
	i.environment.assign(stmt.name, class)
	return nil
}

func (i *Interpreter) visitExpression(stmt *Expression) interface{} {
	value := i.evaluate(stmt.expression)
	if i.interactive {
		fmt.Println(value)
	}
	return nil
}

func (i *Interpreter) visitFunction(stmt *Function) interface{} {
	function := &GloxCallable{
		declaration:   stmt,
		closure:       i.environment,
		isInitializer: false,
	}
	i.environment.define(stmt.name.lexeme, function)
	return nil
}

func (i *Interpreter) visitIf(stmt *If) interface{} {
	if isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
	return nil
}

func (i *Interpreter) visitPrint(stmt *Print) interface{} {
	value := i.evaluate(stmt.expression)
	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) visitReturn(stmt *Return) interface{} {
	var value interface{}
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	} else {
		value = nil
	}
	panic(returnSignal{value})
	return nil
}

func (i *Interpreter) visitVar(stmt *Var) interface{} {
	var value interface{}
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}
	i.environment.define(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) visitWhile(stmt *While) interface{} {
	for isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}
	return nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()
	i.environment = environment
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) lookupVariable(name Token, expr Expr) interface{} {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.getAt(distance, name.lexeme)
	} else {
		return i.globals.get(name)
	}
}

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	} else if b, ok := object.(bool); ok {
		return b
	} else {
		return true
	}
}

func checkNumberOperand(operator Token, operand interface{}) {
	if _, ok := operand.(float64); !ok {
		panic(interpreterError{gloxError(operator, "Operand must be a number.")})
	}
}

func checkNumberOperands(operator Token, left, right interface{}) {
	_, lok := left.(float64)
	_, rok := right.(float64)

	if !lok || !rok {
		panic(interpreterError{gloxError(operator, "Operands must be numbers.")})
	}
}

func checkCallable(callee interface{}, paren Token) {
	if _, ok := callee.(Callable); !ok {
		panic(interpreterError{gloxError(paren, "Only functions and classes are callable.")})
	}
}

func stringify(value interface{}) string {
	if num, ok := value.(float64); ok {
		var text string
		if math.Trunc(num) == num {
			text = fmt.Sprintf("%d", int64(num))
		} else {
			text = fmt.Sprintf("%f", num)
		}
		return text
	} else {
		return fmt.Sprintf("%v", value)
	}
}
