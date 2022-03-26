package glox

import "time"

func NewGlobals() *Environment {
	values := map[string]interface{}{
		"clock": &BuiltinCallable{
			builtinArity: 0,
			builtinFunction: func(interpreter *Interpreter, arguments []interface{}) interface{} {
				return float64(time.Now().Unix())
			},
		},
	}
	return &Environment{
		enclosing: nil,
		values:    values,
	}
}
