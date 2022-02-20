package main

import (
	"bufio"
	"fmt"
	"glox/glox"
	"os"
)

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading the file %s\n", path)
		os.Exit(74)
	}
	if errors := run(string(content)); errors != nil {
		printErrors(errors)
		os.Exit(65)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Scanning error: %v\n", err)
				os.Exit(74)
			} else {
				fmt.Println()
				break
			}
		}
		line := scanner.Text()
		if errors := run(line); errors != nil {
			printErrors(errors)
		}
	}
}
func run(source string) []error {
	scanner := glox.NewScanner(source)
	tokens, errors := scanner.ScanTokens()
	if errors != nil {
		return errors
	}
	parser := glox.NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		return []error{err}
	}
	interpreter := glox.NewInterpreter()
	err = interpreter.Interpret(expr)
	if err != nil {
		return []error{err}
	}
	return nil
}

func printErrors(errors []error) {
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err)
	}
}
func main() {
	programName := os.Args[0]
	numArgs := len(os.Args[1:])
	if numArgs > 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [script]\n", programName)
		os.Exit(64)
	} else if numArgs == 1 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
