package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

func main() {

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		t := NewTokenizer(line)

		t.Tokenize()

		p := NewParser(t.Tokens)

		expr := p.ParseExpr()
		fmt.Printf("AST:%s\n", stringifyAst(expr, 0))

		for _, err := range p.errors {
			fmt.Printf("%s: %d\n", err.message, err.position)
		}

		// var img image.Image
		// switch e := expr.(type) {
		// 	case AstApplication:
		// 		img = 
		// }
		
		diagram := GenDiagWrap(expr)

		file, err := os.Create("diagram.png")
		if err != nil {
			log.Fatalf("Failed creating image: %v", err)
		}
		png.Encode(file, diagram)
	}

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
