package util

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output dir>")
		os.Exit(1)
	}
	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{"Binary : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal : Object value",
		"Unary : Token operator, Expr right"})
}

func defineAst(outputDir, basename string, types []string) error {
	path := outputDir + "/" + basename + ".go"
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file: %w", err)
	}
	file.Write()

	return nil
}
