package glox

type Interpreter struct{}

func (i *Interpreter) visitAssign(expr Assign) interface{} {
	return nil
}

func (i *Interpreter) visitBinary(expr Binary) interface{} {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)
	switch expr.operator.tokenType {
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case MINUS:
		return left.(float64) - right.(float64)
	case BANG_EQUAL:
		return left != right
	case EQUAL:
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
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
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
		return !i.isTruthy(right)
	case MINUS:
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

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	} else if b, ok := object.(bool); ok {
		return b
	} else {
		return true
	}
}
