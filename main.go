package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("invalid arguments")
	} else if len(os.Args) == 2{
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
}

func runPrompt() {
	for {
		fmt.Printf(">")
		var line string
		n, err := fmt.Scanln(&line)
		if err != nil {
			fmt.Println(err)
			break
		}
		if n == 0 {
			break
		}
		run(line)
	}
}

func run(src string) {
	fmt.Println(src)
}