package main

import (
	"image"
	"image/color"
	"image/draw"

	_ "github.com/google/uuid"
)

type ExprDiagram interface {
	Img() image.Image
	Variables() []string
}

type VariDiag struct {
	img  image.Image
	name string
}

func (d VariDiag) Img() image.Image {
	return d.img
}
func (d VariDiag) Variables() []string {
	return []string{ d.name }
}

type AbstDiag struct {
	img image.Image
	variable string
	variables []string
}
func (d AbstDiag) Img() image.Image {
	return d.img
}
func (d AbstDiag) Variables() []string {
	return d.variables
}

type AppDiag struct {
	img image.Image
	variables []string
}

func (d AppDiag) Img() image.Image {
	return d.img
}
func (d AppDiag) Variables() []string {
	return d.variables
}

func GenDiag(expr AstLambdaExpr) ExprDiagram {
	switch e := expr.(type) {
	case AstVariable:
		return GenVariDiag(e)
	case AstAbstraction:
		return GenAbstDiag(e)
	case AstApplication:
		return GenAppDiag(e)
	default:
		return VariDiag{}
	}
}

func GenAbstDiag(abst AstAbstraction) AbstDiag {
	bodyDiag := GenDiag(abst.Body)
	vari := abst.Parameter.Name

	bodyImg := bodyDiag.Img()

	bounds := bodyImg.Bounds()
	clone := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(clone, bounds, bodyImg, image.Point{}, draw.Over)

	for x := range bounds.Dx() {
		clone.Set(x,1,color.Black)
	}

	var variables []string
	for i, v := range bodyDiag.Variables() {
		if v == vari {
			variables = append(variables, "")	
			clone.Set(1+4*i, 0, color.Transparent)
		} else {
			variables = append(variables, v)
		}
	}

	wrap := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()+2))
	draw.Draw(wrap, bounds.Add(image.Pt(0,2)), clone, image.Point{}, draw.Over)

	for x := range bounds.Dx() {
		wrap.Set(x,1,clone.At(x,0))
		wrap.Set(x,0,clone.At(x,0))
	}

	return AbstDiag{
		img: wrap,
		variables: variables,
	}
}

func GenVariDiag(variable AstVariable) VariDiag {
	img := image.NewRGBA(image.Rect(0, 0, 3, 2))

	img.Set(1, 0, color.Black)
	img.Set(1, 1, color.Black)

	return VariDiag{
		img:  img,
		name: variable.Name,
	}
}

func GenAppDiag(expr AstApplication) AppDiag {
	targDiagram := GenDiag(expr.Abstraction)
	exprDiagram := GenDiag(expr.Argument)

	targBounds := targDiagram.Img().Bounds()
	exprBounds := exprDiagram.Img().Bounds()

	width := targBounds.Dx() + exprBounds.Dx() + 1
	height := max(targBounds.Dy(), exprBounds.Dy()) + 2

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// draw variable
	draw.Draw(
		img,
		image.Rect(
			0,
			0,
			targBounds.Dx(),
			targBounds.Dy(),
		),
		targDiagram.Img(),
		image.Point{},
		draw.Over,
	)

	// draw target line
	for y := range height - targBounds.Dy() {
		img.Set(1, targBounds.Dy()+y, color.Black)
	}

	// draw expression
	draw.Draw(
		img,
		image.Rect(
			targBounds.Dx()+1,
			0,
			width-1,
			exprBounds.Dy(),
		),
		exprDiagram.Img(),
		image.Point{},
		draw.Over,
	)

	// draw argument line
	for y := range height - exprBounds.Dy() {
		img.Set(targBounds.Dx()+2, exprBounds.Dy()+y, color.Black)
	}

	// draw application line
	for x := range targBounds.Dx() + 2 {
		img.Set(1+x, height-1, color.Black)
	}
	img.Set(targBounds.Dx()+2, height-2, color.Black)

	return AppDiag{
		img: img,
		variables: append(targDiagram.Variables(), exprDiagram.Variables()...),
	}
}

type ExpandDirection int

const (
	Top ExpandDirection = iota
	Bottom
	Right
	Left
)

func ImageExpandRightOrBottom(img image.Image, n int, direction ExpandDirection) image.Image {
	bounds := img.Bounds()
	wrap := image.NewRGBA(image.Rect(0, 0, bounds.Dx()+n, bounds.Dy()+n))
	draw.Draw(wrap, img.Bounds(), img, image.Point{}, draw.Over)

	return wrap
}

func GenDiagWrap(expr AstLambdaExpr) image.Image {
	diagram := GenDiag(expr)
	img := diagram.Img()
	bounds := img.Bounds()
	wrap := image.NewRGBA(image.Rect(0, 0, bounds.Dx()+2, bounds.Dy()+1))
	draw.Draw(wrap, wrap.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(wrap, img.Bounds().Add(image.Pt(1, -1)), img, image.Point{}, draw.Over)
	return wrap
}
