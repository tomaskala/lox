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
	if value, ok := e.values[name.lexeme]; ok {
		return value
	} else if e.enclosing != nil {
		return e.enclosing.get(name)
	} else {
		message := fmt.Sprintf("Undefined variable '%s'.", name.lexeme)
		panic(interpreterError{gloxError(name, message)})
	}
}

func (e *Environment) getAt(distance int, name string) interface{} {
	return e.ancestor(distance).values[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) assign(name Token, value interface{}) {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
	} else if e.enclosing != nil {
		e.enclosing.assign(name, value)
	} else {
		message := fmt.Sprintf("Undefined variable '%s'.", name.lexeme)
		panic(interpreterError{gloxError(name, message)})
	}
}

func (e *Environment) assignAt(distance int, name Token, value interface{}) {
	e.ancestor(distance).values[name.lexeme] = value
}
