package glox

import "fmt"

type GloxClass struct {
	name    string
	methods map[string]*GloxCallable
}

type GloxInstance struct {
	class  *GloxClass
	fields map[string]interface{}
}

func (g *GloxClass) arity() int {
	return 0
}

func (g *GloxClass) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return &GloxInstance{
		class:  g,
		fields: make(map[string]interface{}),
	}
}

func (g *GloxClass) findMethod(name string) *GloxCallable {
	if method, ok := g.methods[name]; ok {
		return method
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
		return method
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
