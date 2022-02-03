package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestTrim(t *testing.T) {
	src := `type Binary struct {
    left Expr
    operator Token
    right Expr
}
func NewBinary(left Expr,operator Token,right Expr,`
	fmt.Println(strings.TrimRight(src, ","))
}
