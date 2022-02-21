package glox

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(name Token) interface{} {
	value, ok := e.values[name.lexeme]
	if ok {
		return value
	} else if e.enclosing != nil {
		return e.enclosing.get(name)
	} else {
		message := fmt.Sprintf("Undefined variable '%s'.", name.lexeme)
		panic(interpreterError{runtimeError(name, message)})
	}
}

func (e *Environment) assign(name Token, value interface{}) {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
	} else if e.enclosing != nil {
		e.enclosing.assign(name, value)
	} else {
		message := fmt.Sprintf("Undefined variable '%s'.", name.lexeme)
		panic(interpreterError{runtimeError(name, message)})
	}
}
