package main

type AstMacroDefinition struct {
	Name  string
	Value AstLambdaExpr
}

type AstLambdaExpr interface{}

type AstVariable struct {
	Name string
}

type AstAbstraction struct {
	Parameter AstVariable
	Body AstLambdaExpr
}

type AstApplication struct {
	Abstraction AstLambdaExpr
	Argument AstLambdaExpr
}
