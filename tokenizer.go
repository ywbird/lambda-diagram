package main

import "unicode"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	INT
	VARIABLE
	ASSIGN
	PLUS
	MINUS
	MULT
	LPAREN
	RPAREN
	LAMBDA
	MACRO
	DOT
)

var TokenName = map[TokenType]string{
	ILLEGAL:  "illegal",
	EOF:      "eof",
	INT:      "int",
	VARIABLE: "indenifier",
	ASSIGN:   "assign",
	PLUS:     "plus",
	MINUS:    "minus",
	MULT:     "mult",
	LPAREN:   "lparen",
	RPAREN:   "rparen",
	LAMBDA:   "lambda",
	MACRO:    "macro",
	DOT:      "dot",
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     [2]int
}
type Tokenizer struct {
	src    string
	pos    int
	Tokens []Token
	last   string
}

func NewTokenizer(source string) Tokenizer {
	return Tokenizer{
		src: source,
		pos: 0,
	}
}

func (t *Tokenizer) Tokenize() {
	if len(t.src) <= 0 {
		return
	}
	for {
		char := byte(t.src[t.pos])

		switch char {
		case '=':
			t.Tokens = append(t.Tokens, Token{ASSIGN, "=", [2]int{t.pos, t.pos+1}})
		case '+':
			t.Tokens = append(t.Tokens, Token{PLUS, "+", [2]int{t.pos, t.pos+1}})
		case '-':
			t.Tokens = append(t.Tokens, Token{MINUS, "-", [2]int{t.pos, t.pos+1}})
		case '*':
			t.Tokens = append(t.Tokens, Token{MULT, "*", [2]int{t.pos, t.pos+1}})
		case '(':
			t.Tokens = append(t.Tokens, Token{LPAREN, "(", [2]int{t.pos, t.pos+1}})
		case ')':
			t.Tokens = append(t.Tokens, Token{RPAREN, ")", [2]int{t.pos, t.pos+1}})
		case '\\':
			t.Tokens = append(t.Tokens, Token{LAMBDA, "\\", [2]int{t.pos, t.pos+1}})
		case '.':
			t.Tokens = append(t.Tokens, Token{DOT, ".", [2]int{t.pos, t.pos+1}})
		default:
			{
				if char == '$' {
					var word []byte
					var nchar byte
					start := t.pos
					for {
						if !(t.pos+1 >= len(t.src)) {
							nchar = t.src[t.pos+1]
						} else {
							nchar = 0
						}
						char = t.src[t.pos]
						word = append(word, char)
						if t.pos+1 >= len(t.src) || !(('a' <= nchar && nchar <= 'z') ||
							('A' <= nchar && nchar <= 'Z')) {
							t.Tokens = append(t.Tokens, Token{MACRO, string(word), [2]int{start, t.pos+1}})
							break
						}
						t.pos++
					}
				} else if '0' <= char && char <= '9' {
					var word []byte
					var nchar byte
					start := t.pos
					for {
						if !(t.pos+1 >= len(t.src)) {
							nchar = t.src[t.pos+1]
						} else {
							nchar = 0
						}
						char = t.src[t.pos]
						word = append(word, char)
						if t.pos+1 >= len(t.src) || !('0' <= nchar && char <= '9') {
							t.Tokens = append(t.Tokens, Token{INT, string(word), [2]int{start, t.pos+1}})
							break
						}
						t.pos++
					}
				} else if !unicode.IsSpace(rune(char)) && (('a' <= char && char <= 'z') ||
					('A' <= char && char <= 'Z')) {
					t.Tokens = append(t.Tokens, Token{VARIABLE, string(char), [2]int{t.pos, t.pos+1} })
					// var word []byte
					// var nchar byte
					// start := t.pos
					// for {
					// 	if !(t.pos+1 >= len(t.src)) {
					// 		nchar = t.src[t.pos+1]
					// 	} else {
					// 		nchar = 0
					// 	}
					// 	char = t.src[t.pos]
					// 	word = append(word, char)
					// 	if t.pos+1 >= len(t.src) || !(('a' <= nchar && nchar <= 'z') ||
					// 		('A' <= nchar && nchar <= 'Z') || ('0' <= nchar && nchar <= '9')) {
					// 		t.Tokens = append(t.Tokens, Token{VARIABLE, string(word), [2]int{start, t.pos+1}})
					// 		break
					// 	}
					// 	t.pos++
					// }
				} else if !unicode.IsSpace(rune(char)) {
					t.Tokens = append(t.Tokens, Token{ILLEGAL, string(char), [2]int{t.pos, t.pos+1}})
				}
			}
		}

		if t.pos+1 >= len(t.src) {
			break
		}
		t.pos++
	}
}
