package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func read(str string) string {
	return str
}

func eval(ast string) string {
	return ast
}

func print(exp string) string {
	return exp
}

func rep(str string) string {
	return print(eval(read(str)))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("user> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(rep(strings.TrimRight(input, "\n")))
	}
}
