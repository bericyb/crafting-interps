package main

import (
	"bufio"
	"fmt"
	"os"

	"go_interp/internal/utils"
)

func main() {
	fmt.Println("Zox interpreter starting...")

	fmt.Println("args: ", len(os.Args))

	if len(os.Args) > 2 {
		fmt.Println("Usage: jlox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		err := runFile(os.Args[1])
		if err != nil {
			fmt.Sprintln(fmt.Sprintf("An error has occurred: %v", err))
		}
	} else {
		err := runPrompt()
		if err != nil {
			fmt.Sprintln(fmt.Sprintf("An error has occurred: %v", err))
		}
	}
}

type lox struct {
	hadError bool
}

func runFile(filepath string) error {
	fmt.Println(fmt.Sprintf("Reading file: %s", filepath))
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return run(string(bytes))
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for true {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		err = run(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func run(source string) error {
	scanner := utils.Scanner{
		Source: source,
	}

	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.ToString())
	}

	parser := utils.Parser{
		Tokens:  tokens,
		Current: 0,
	}

	expression, error := parser.Parse()
	if error != nil {
		fmt.Println("We got an error, fuck!", error)
		return error
	}

	printer := utils.AstPrinter{}

	fmt.Println(printer.Print(expression))

	return nil
}
