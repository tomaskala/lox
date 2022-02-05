package glox

import "fmt"

func gloxError(line int, message string) error {
	return report(line, "", message)
}

func report(line int, where, message string) error {
	return fmt.Errorf("[line %d] Error %s: %s", line, where, message)
}
