package main

import (
	"flag"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
)

func main() {
	exprPtr := flag.String("expr", "\\x.x", "Expression")
	outputPtr := flag.String("out", "diagram.png", "Output file")
	scalePtr := flag.Int("scale", 10, "Output scale")
	verbosePtr := flag.Bool("verbose", false, "Print AST")

	flag.Parse()

	t := NewTokenizer(*exprPtr)

	t.Tokenize()

	p := NewParser(t.Tokens)

	expr := p.ParseExpr()

	if *verbosePtr {
		fmt.Printf("AST:%s\n", stringifyAst(expr, 0))
	}
	//
	// for _, err := range p.errors {
	// 	fmt.Printf("%s: %d\n", err.message, err.position)
	// }

	diagram := GenDiagWrap(expr)

	scaledSize := image.Rect(0, 0, diagram.Bounds().Dx()**scalePtr, diagram.Bounds().Dy()**scalePtr) // Example: halve the size
	resizedImage := image.NewRGBA(scaledSize)

	draw.NearestNeighbor.Scale(resizedImage, scaledSize, diagram, diagram.Bounds(), draw.Over, nil)

	file, err := os.Create(*outputPtr)
	if err != nil {
		log.Fatalf("Failed creating image: %v", err)
	}
	png.Encode(file, resizedImage)
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
