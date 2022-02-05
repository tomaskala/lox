package main

import (
	"bufio"
	"fmt"
	"os"
)

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading the file %s\n", path)
		os.Exit(74)
	}
	if err := run(string(content)); err != nil {
		fmt.Fprintln(os.Stderr, err)
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
		if err := run(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func run(source string) error {
	// TODO
	return nil
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
