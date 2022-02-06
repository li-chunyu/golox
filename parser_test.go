package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	src := `1-6/2`
	s := NewScanner(src)
	toks := s.scanTokens()
	p := NewParser(toks)
	expr := p.Parse()
	fmt.Println(expr.Accept(&AstPrinter{}))
	assert.Equal(t, "(- 1 (/ 6 2))", expr.Accept(&AstPrinter{}))
}
