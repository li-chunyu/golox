package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	hadError = false
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("invalid arguments")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(src string) {
	f, err := os.Open(src)
	if err != nil {
		fmt.Printf("Open file %v fail: %v.\n", src, err)
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Open file %v failL %v.\n", src, err)
		return
	}

	run(string(b))
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	for {
		line := readline(">")
		run(line)
		hadError = false
	}
}

func readline(prompt string) string {
	fmt.Printf("%v", prompt)
	var (
		c   rune
		b   []rune
		err error
	)

	for err == nil {
		_, err = fmt.Scanf("%c", &c)
		if c != '\n' {
			b = append(b, c)
		} else {
			break
		}
	}
	return string(b)
}

func run(src string) {
	s := NewScanner(src)
	s.scanTokens()
	for _, t := range s.toks {
		fmt.Println(t)
	}
}

func perror(line int, msg string) {
	report(line, "", msg)
}

func report(line int, where, msg string) {
	fmt.Println(errorMsg(line, where, msg))
	hadError = true
}

func errorMsg(line int, where, msg string) string {
	return fmt.Sprintf("[line %v] Error %v: %v.\n", line, where, msg)
}
