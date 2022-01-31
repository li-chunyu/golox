package main

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	runFile("test/src1.txt")
}

func TestReport(t *testing.T) {
	report(5, "somewhere in your code", "wrong type")
}
