package glox

import "fmt"

type GloxClass struct {
	name       string
	superclass *GloxClass
	methods    map[string]*GloxCallable
}

type GloxInstance struct {
	class  *GloxClass
	fields map[string]interface{}
}

func (g *GloxClass) arity() int {
	if initializer := g.findMethod("init"); initializer != nil {
		return initializer.arity()
	} else {
		return 0
	}
}

func (g *GloxClass) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	instance := &GloxInstance{
		class:  g,
		fields: make(map[string]interface{}),
	}
	if initializer := g.findMethod("init"); initializer != nil {
		initializer.bind(instance).call(interpreter, arguments)
	}
	return instance
}

func (g *GloxClass) findMethod(name string) *GloxCallable {
	if method, ok := g.methods[name]; ok {
		return method
	} else if g.superclass != nil {
		return g.superclass.findMethod(name)
	} else {
		return nil
	}
}

func (g GloxClass) String() string {
	return g.name
}

func (g *GloxInstance) get(name Token) interface{} {
	if field, ok := g.fields[name.lexeme]; ok {
		return field
	} else if method := g.class.findMethod(name.lexeme); method != nil {
		return method.bind(g)
	} else {
		message := fmt.Sprintf("Undefined property '%s'.", name.lexeme)
		panic(interpreterError{gloxError(name, message)})
	}
}

func (g *GloxInstance) set(name Token, value interface{}) {
	g.fields[name.lexeme] = value
}

func (g GloxInstance) String() string {
	return fmt.Sprintf("%s instance", g.class.name)
}
