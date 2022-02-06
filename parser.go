package main

import (
	"fmt"
)

type Parser struct {
	toks    []*Token
	current int
}

func NewParser(toks []*Token) *Parser {
	return &Parser{
		toks: toks,
	}
}

func (p *Parser) Parse() (e Expr) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			p.synchronize()
			e = p.expression()
		}
	}()
	e = p.expression()
	return
}

func (p *Parser) expression() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr /* left association */, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	term := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		term = NewBinary(term, op, right)
	}
	return term
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		op := p.previous()
		r := p.factor()
		expr = NewBinary(expr, op, r)
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		op := p.previous()
		r := p.unary()
		expr = NewBinary(expr, op, r)
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(MINUS, BANG) {
		r := p.unary()
		op := p.previous()
		return NewUnary(op, r)
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().literal)
	}
	if p.match(TRUE) {
		return NewLiteral(true)
	}
	if p.match(FALSE) {
		return NewLiteral(false)
	}
	if p.match(NIL) {
		return NewLiteral(nil)
	}
	if p.match(LEFT_PAREN) {
		e := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(e)
	}
	p.throwPanic(p.peek(), "Error parser, invalid prime type.")
	return nil // avoid compiler error
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t TokenType, msg string) *Token {
	if p.check(t) {
		return p.advance()
	}
	p.throwPanic(p.peek(), msg)
	return nil // avoid compile error
}

func (p *Parser) peek() *Token {
	return p.toks[p.current]
}

func (p *Parser) check(t TokenType) bool {
	return p.toks[p.current].typ == t
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.current > len(p.toks)
}

func (p *Parser) previous() *Token {
	i := p.current - 1
	return p.toks[i]
}

func (p *Parser) throwPanic(tok *Token, msg string) {
	if tok.typ == EOF {
		panic(errorMsg(tok.line, "at end", msg))
	} else {
		panic(errorMsg(tok.line, "at '"+tok.lexeme+"'", msg))
	}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().typ == SEMICOLON {
			return
		}

		switch p.peek().typ {
		case VAR:
			fallthrough
		case FUN:
			fallthrough
		case CLASS:
			fallthrough
		case FOR:
			fallthrough
		case IF:
			fallthrough
		case WHILE:
			fallthrough
		case PRINT:
			fallthrough
		case RETURN:
			return
		}
		p.advance()
	}
}
