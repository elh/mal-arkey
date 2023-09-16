package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	mal "github.com/elh/mal-arkey"
)

// Read–eval–print. recover panics here so that REPL can continue accepting stdin
func rep(str string, env *mal.Env) (out string) {
	defer func() {
		if r := recover(); r != nil {
			const colorRed, colorReset = "\033[31m", "\033[0m"
			fmt.Printf("%sError: %s%s", colorRed, r, colorReset)
		}
	}()
	return mal.Print(mal.Eval(mal.Read(str), env), true)
}

// Starts the Mal-arkey REPL. If command line args are provided, the first arg is treated as a file to load, and the
// remaining are bound as a list to `*ARGV*`.
func main() {
	reader := bufio.NewReader(os.Stdin)
	env := mal.BuiltinEnv()

	// self-hosted fns
	rep(`(def! not (fn* (a) (if a false true)))`, env)
	rep(`(defmacro! cond (fn* (& xs) (if (> (count xs) 0) (list 'if (first xs) (if (> (count xs) 1) (nth xs 1) (throw "odd number of forms to cond")) (cons 'cond (rest (rest xs)))))))`, env)
	rep(`(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\nnil)")))))`, env)

	if len(os.Args) > 1 {
		var vals []mal.Value
		for _, arg := range os.Args[2:] {
			vals = append(vals, mal.Value{Type: "string", Val: arg})
		}
		env.Set("*ARGV*", mal.Value{Type: "list", Val: vals})
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
