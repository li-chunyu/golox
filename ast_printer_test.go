package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinter(t *testing.T) {
	expr := NewBinary(
		NewUnary(
			NewToken(MINUS, nil, "-", 1),
			NewLiteral(123.0),
		),
		NewToken(STAR, nil, "*", 1),
		NewGrouping(NewLiteral(45.67)),
	)
	s := expr.Accept(&AstPrinter{}).(string)
	assert.Equal(t, s, "(* (- 123) (group 45.67))")
}
