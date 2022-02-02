package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	var src string
	var s *Scanner

	src = `"abc"`
	s = NewScanner(src)
	s.scanTokens()
	assert.Equal(t, 2, len(s.toks))
	assert.Equal(t, s.toks[0].typ, STRING)
	assert.Equal(t, s.toks[0].literal.(string), "abc")
	assert.Equal(t, s.toks[1].typ, EOF)

	// test float
	src = `123.456`
	s = NewScanner(src)
	s.scanTokens()
	assert.Equal(t, 2, len(s.toks))
	assert.Equal(t, s.toks[0].typ, NUMBER)
	assert.Equal(t, s.toks[0].literal.(float64), 123.456)
	assert.Equal(t, s.toks[1].typ, EOF)

	// test identifier and keywords
	src = `_ab123c`
	s = NewScanner(src)
	s.scanTokens()
	assert.Equal(t, 2, len(s.toks))
	assert.Equal(t, s.toks[0].typ, IDENTIFIER)
	assert.Equal(t, s.toks[0].lexeme, src)

	src = `ort`
	s = NewScanner(src)
	s.scanTokens()
	assert.Equal(t, 2, len(s.toks))
	assert.Equal(t, s.toks[0].typ, IDENTIFIER)
	assert.Equal(t, s.toks[0].lexeme, src)

	src = `or`
	s = NewScanner(src)
	s.scanTokens()
	assert.Equal(t, 2, len(s.toks))
	assert.Equal(t, s.toks[0].typ, OR)

	src = `// this is a comment
var a= false`
	s = NewScanner(src)
	s.scanTokens()
	assert.Equal(t, 5, len(s.toks))
	assert.Equal(t, s.toks[0].typ, VAR)
	assert.Equal(t, s.toks[1].typ, IDENTIFIER)
	assert.Equal(t, s.toks[2].typ, EQUAL)
	assert.Equal(t, s.toks[3].typ, FALSE)
}

func TestRune(t *testing.T) {
	s := "你好"
	r := strings.NewReader(s)
	fmt.Println(r.ReadRune())
	fmt.Println(r.ReadRune())
	fmt.Println(r.ReadRune())
}
