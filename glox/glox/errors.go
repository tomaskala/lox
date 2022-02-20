package glox

import "fmt"

func scanError(line int, message string) error {
	return report(line, "", message)
}

func parseError(token Token, message string) error {
	var where string
	if token.tokenType == EOF {
		where = "at end"
	} else {
		where = fmt.Sprintf("at '%s'", token.lexeme)
	}
	return report(token.line, where, message)
}

func runtimeError(token Token, message string) error {
	var where string
	if token.tokenType == EOF {
		where = "at end"
	} else {
		where = fmt.Sprintf("at '%s'", token.lexeme)
	}
	return report(token.line, where, message)
}

func report(line int, where, message string) error {
	return fmt.Errorf("[line %d] Error %s: %s", line, where, message)
}
