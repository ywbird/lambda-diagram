package main

import (
	"fmt"
)

type ParseError struct {
	message  string
	position int
}

type Parser struct {
	toks   []Token
	pos    int
	errors []ParseError
}

func NewParser(tokens []Token) Parser {
	return Parser{
		toks: tokens,
	}
}

func (p *Parser) NextToken() Token {
	if p.pos == len(p.toks) {
		return Token{EOF, "<EOF>", [2]int{0, 0}}
	}
	token := p.toks[p.pos]
	p.pos++
	return token
}

func (p *Parser) CurrentToken() Token {
	if p.pos == len(p.toks) {
		return Token{EOF, "<EOF>", [2]int{0, 0}}
	}
	return p.toks[p.pos]
}

func (p *Parser) RaiseError(msg string, pos int) {
	p.errors = append(p.errors, ParseError{msg, pos})
}

func (p *Parser) ParseExpr() AstLambdaExpr {
	// parse every expression
	var expr AstLambdaExpr

	switch curr := p.CurrentToken(); curr.Type {
	case ILLEGAL:
		p.RaiseError(fmt.Sprintf("Illegal character `%s`.", curr.Literal), curr.Pos[0])
	case LAMBDA:
		expr = p.ParseAbstraction()
	case VARIABLE:
		expr = p.ParseVariable()
	case LPAREN:
		{
			p.NextToken()
			expr := p.ParseExpr()

			next := p.CurrentToken()
			if next.Type == RPAREN {
				p.NextToken()
				return expr
			}

			for {
				arg := p.ParseExpr()
				expr = AstApplication{
					expr, arg,
				}

				next = p.CurrentToken()
				if next.Type == RPAREN {
					p.NextToken()
					break
				}
			}
			return expr
		}
	default:
		p.RaiseError(fmt.Sprintf("Unexpected token `%s`", curr.Literal), curr.Pos[0])
	}

	return expr
}

func (p *Parser) ParseVariable() AstVariable {
	variable := p.NextToken()
	if VARIABLE != variable.Type {
		p.RaiseError("Expected variable identifier.", variable.Pos[0])
	}
	return AstVariable{
		Name: variable.Literal,
	}
}

func (p *Parser) ParseAbstraction() AstAbstraction {
	if t := p.NextToken(); LAMBDA != t.Type {
		p.RaiseError("Expected lambda (`\\`).", t.Pos[0])
	}

	param := p.NextToken()
	if VARIABLE != param.Type {
		p.RaiseError("Expected variable identifier.", param.Pos[0])
	}
	if d := p.NextToken(); DOT != d.Type {
		p.RaiseError("Expected `.` next to variable in abstraction.", d.Pos[0])
	}

	body := p.ParseExpr()

	return AstAbstraction{
		Parameter: AstVariable{ param.Literal },
		Body:      body,
	}
}

