package glox

import "fmt"

type Callable interface {
	arity() int
	call(interpreter *Interpreter, arguments []interface{}) interface{}
}

type LoxCallable struct {
	declaration Function
}

func (c LoxCallable) arity() int {
	return len(c.declaration.params)
}

func (c LoxCallable) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	environment := NewEnvironment(interpreter.globals)
	for i := 0; i < len(c.declaration.params); i++ {
		environment.define(c.declaration.params[i].lexeme, arguments[i])
	}
	interpreter.executeBlock(c.declaration.body, environment)
	return nil
}

func (c LoxCallable) String() string {
	return fmt.Sprintf("<fn %s>", c.declaration.name.lexeme)
}

type BuiltinCallable struct {
	builtinArity    int
	builtinFunction func(interpreter *Interpreter, arguments []interface{}) interface{}
}

func (b BuiltinCallable) arity() int {
	return b.builtinArity
}

func (b BuiltinCallable) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return b.builtinFunction(interpreter, arguments)
}

func (b BuiltinCallable) String() string {
	return "<builtin fn>"
}
