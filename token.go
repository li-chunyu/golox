package main

import (
	"fmt"
)

func init() {
	keywords = make(map[string]TokenType)
	keywords["and"] = AND
	keywords["class"] = CLASS
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["for"] = FOR
	keywords["fun"] = FUN
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["var"] = VAR
	keywords["while"] = WHILE
}

type TokenType int

const (
	// Single-character tokens.

	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL

	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.

	IDENTIFIER
	STRING
	NUMBER

	// Keywords.

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var keywords map[string]TokenType

func NewToken(t TokenType, lit interface{}, lex string, l int) *Token {
	tok := Token{
		typ:     t,
		literal: lit,
		lexeme:  lex,
		line:    l,
	}
	return &tok
}

type Token struct {
	typ     TokenType
	literal interface{}
	lexeme  string
	line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", t.typ, t.literal, t.lexeme)
}
