package glox

import "fmt"

// Wraps a parsing error to distinguish it from other errors.
type parserError struct{ error }

// Wraps an interpreter error to distinguish it from other errors.
type interpreterError struct{ error }

// Wraps a resolver error to distinguish it from other errors.
type resolverError struct{ error }

func scanError(line int, message string) error {
	return report(line, "", message)
}

func gloxError(token Token, message string) error {
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
