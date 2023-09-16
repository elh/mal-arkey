package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	m "github.com/elh/mal-arkey/internal/pkg/mal"
)

func printError(err any) {
	const colorRed = "\033[31m"
	const colorReset = "\033[0m"
	fmt.Printf("%sError: %s%s", colorRed, err, colorReset)
}

func read(str string) m.Sexpr {
	return m.ReadStr(str)
}

func eval(expr m.Sexpr, env *m.Env) m.Sexpr {
	return m.Eval(expr, env)
}

func print(expr m.Sexpr) string {
	return m.PrintStr(expr, true)
}

func rep(str string, env *m.Env) (out string) {
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
	env := m.BuiltInEnv()

	// self-hosted fns
	rep(`(def! not (fn* (a) (if a false true)))`, env)
	rep(`(defmacro! cond (fn* (& xs) (if (> (count xs) 0) (list 'if (first xs) (if (> (count xs) 1) (nth xs 1) (throw "odd number of forms to cond")) (cons 'cond (rest (rest xs)))))))`, env)
	rep(`(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\nnil)")))))`, env)

	if len(os.Args) > 1 {
		var vals []m.Sexpr
		for _, arg := range os.Args[2:] {
			vals = append(vals, m.Sexpr{Type: "string", Val: arg})
		}
		env.Set("*ARGV*", m.Sexpr{Type: "list", Val: vals})
		rep(fmt.Sprintf("(load-file \"%s\")", os.Args[1]), env)
		return
	}

	rep(`(println (str "Mal [" *host-language* "]"))`, env)
	for {
		fmt.Print("user> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(rep(strings.TrimRight(input, "\n"), env))
	}
}
