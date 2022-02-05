package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading the file %s\n", path)
	}
	run(string(content))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				log.Fatalf("Unexpected error: %v\n", err)
			} else {
				fmt.Println()
				break
			}
		}
		line := scanner.Text()
		run(line)
	}
}

func run(source string) {
	// TODO
}

func main() {
	programName := os.Args[0]
	numArgs := len(os.Args[1:])
	if numArgs > 1 {
		log.Fatalf("Usage: %s [script]\n", programName)
	} else if numArgs == 1 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
