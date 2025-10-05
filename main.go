package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"
	"flag"
	// "github.com/chzyer/readline"
)

func main() {
	exprPtr := flag.String("expr", "\\x.x", "Expression")
	outputPtr := flag.String("out", "diagram.png", "Output file")

	flag.Parse()

	t := NewTokenizer(*exprPtr)

	t.Tokenize()

	p := NewParser(t.Tokens)

	expr := p.ParseExpr()
	// fmt.Printf("AST:%s\n", stringifyAst(expr, 0))
	//
	// for _, err := range p.errors {
	// 	fmt.Printf("%s: %d\n", err.message, err.position)
	// }

	diagram := GenDiagWrap(expr)

	file, err := os.Create(*outputPtr)
	if err != nil {
		log.Fatalf("Failed creating image: %v", err)
	}
	png.Encode(file, diagram)
}

func stringifyAst(ast AstLambdaExpr, indent int) string {
	switch a := ast.(type) {
	case AstApplication:
		return fmt.Sprintf(
			"\n%s|-Application:\n%s|-abstraction: %s\n%s|-argument:%s",
			strings.Repeat("| ", indent),
			strings.Repeat("| ", indent+1),
			stringifyAst(a.Abstraction, indent+2),
			strings.Repeat("| ", indent+1),
			stringifyAst(a.Argument, indent+2),
		)
	case AstVariable:
		return fmt.Sprintf(
			"\n%s|-Variable(%s)",
			strings.Repeat("| ", indent),
			a.Name,
		)
	case AstAbstraction:
		return fmt.Sprintf(
			"\n%s|-Abstraction:\n%s|-parameter: %s\n%s|-body:%s",
			strings.Repeat("| ", indent),
			strings.Repeat("| ", indent+1),
			stringifyAst(a.Parameter, indent+2),
			strings.Repeat("| ", indent+1),
			stringifyAst(a.Body, indent+2),
		)
	}
	return "\n"
}
