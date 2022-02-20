package glox

import "fmt"

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(name Token) interface{} {
	value, ok := e.values[name.lexeme]
	if !ok {
		message := fmt.Sprintf("Undefined variable '%s'.", name.lexeme)
		panic(interpreterError{runtimeError(name, message)})
	} else {
		return value
	}
}

func (e *Environment) assign(name Token, value interface{}) {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
	} else {
		message := fmt.Sprintf("Undefined variable '%s'.", name.lexeme)
		panic(interpreterError{runtimeError(name, message)})
	}
}
