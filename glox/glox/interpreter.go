package glox

import (
	"fmt"
	"math"
)

type Interpreter struct{}

// Wraps an interpreter error to distinguish it from other errors.
type interpreterError struct{ error }

func NewInterpreter() *Interpreter {
	return &Interpreter{}
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

func (i *Interpreter) execute(stmt Stmt) interface{} {
	return stmt.accept(i)
}

func (i *Interpreter) visitBlock(stmt Block) interface{} {
	return nil
}

func (i *Interpreter) visitClass(stmt Class) interface{} {
	return nil
}

func (i *Interpreter) visitExpression(stmt Expression) interface{} {
	i.evaluate(stmt.expression)
	return nil
}

func (i *Interpreter) visitFunction(stmt Function) interface{} {
	return nil
}

func (i *Interpreter) visitIf(stmt If) interface{} {
	return nil
}

func (i *Interpreter) visitPrint(stmt Print) interface{} {
	value := i.evaluate(stmt.expression)
	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) visitReturn(stmt Return) interface{} {
	return nil
}

func (i *Interpreter) visitVar(stmt Var) interface{} {
	return nil
}

func (i *Interpreter) visitWhile(stmt While) interface{} {
	return nil
}

func (i *Interpreter) visitAssign(expr Assign) interface{} {
	return nil
}

func (i *Interpreter) visitBinary(expr Binary) interface{} {
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
		panic(interpreterError{runtimeError(expr.operator, "Operands must be numbers or strings.")})
	case SLASH:
		checkNumberOperands(expr.operator, left, right)
		rightNum := right.(float64)
		if rightNum == 0.0 {
			panic(interpreterError{runtimeError(expr.operator, "Division by zero.")})
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

func (i *Interpreter) visitCall(expr Call) interface{} {
	return nil
}

func (i *Interpreter) visitGet(expr Get) interface{} {
	return nil
}

func (i *Interpreter) visitGrouping(expr Grouping) interface{} {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) visitLiteral(expr Literal) interface{} {
	return expr.value
}

func (i *Interpreter) visitLogical(expr Logical) interface{} {
	return nil
}

func (i *Interpreter) visitSet(expr Set) interface{} {
	return nil
}

func (i *Interpreter) visitSuper(expr Super) interface{} {
	return nil
}

func (i *Interpreter) visitThis(expr This) interface{} {
	return nil
}

func (i *Interpreter) visitUnary(expr Unary) interface{} {
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

func (i *Interpreter) visitVariable(expr Variable) interface{} {
	return nil
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(i)
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
		panic(interpreterError{runtimeError(operator, "Operand must be a number.")})
	}
}

func checkNumberOperands(operator Token, left, right interface{}) {
	_, lok := left.(float64)
	_, rok := right.(float64)

	if !lok || !rok {
		panic(interpreterError{runtimeError(operator, "Operands must be numbers.")})
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
