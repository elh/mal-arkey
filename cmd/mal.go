package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/elh/mal-go/internal/pkg/ast"
	"github.com/elh/mal-go/internal/pkg/evaluator"
	"github.com/elh/mal-go/internal/pkg/printer"
	"github.com/elh/mal-go/internal/pkg/reader"
)

func printError(err any) {
	const colorRed = "\033[31m"
	const colorReset = "\033[0m"
	fmt.Printf("%sError: %s%s", colorRed, err, colorReset)
}

func read(str string) ast.Sexpr {
	return reader.ReadStr(str)
}

func eval(expr ast.Sexpr, env *evaluator.Env) ast.Sexpr {
	return evaluator.Eval(expr, env)
}

func print(expr ast.Sexpr) string {
	return printer.PrintStr(expr)
}

func rep(str string, env *evaluator.Env) (out string) {
	// read, eval, print functions panic
	// recover here so that repl main loop can continue accepting all of stdin.
	defer func() {
		if r := recover(); r != nil {
			printError(r)
			out = ""
		}
	}()

	return print(eval(read(str), env))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	env := evaluator.BuiltInEnv()
	for {
		fmt.Print("user> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(rep(strings.TrimRight(input, "\n"), env))
	}
}
