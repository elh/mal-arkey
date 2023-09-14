package mal

import (
	"fmt"
	"os"
)

// Env is a map of symbols to bound values.
type Env struct {
	outer    *Env
	bindings map[string]Sexpr
}

// NewEnv creates a new environment with the given outer environment.
func NewEnv(outer *Env, bindSymbols []Sexpr, bindValues []Sexpr) *Env {
	bindings := map[string]Sexpr{}
	if len(bindSymbols) != len(bindValues) {
		panic("wrong number of binding symbols and values")
	}
	for i, v := range bindValues {
		bindings[bindSymbols[i].Val.(string)] = Eval(v, outer)
	}

	return &Env{
		outer:    outer,
		bindings: bindings,
	}
}

// Set binds a symbol to a value in the current environment.
func (e *Env) Set(symbol string, value Sexpr) {
	e.bindings[symbol] = value
}

// Get returns the value bound to the given symbol in the environment.
func (e *Env) Get(symbol string) Sexpr {
	if val, ok := e.bindings[symbol]; ok {
		return val
	}
	if e.outer != nil {
		return e.outer.Get(symbol)
	}
	panic(fmt.Sprintf("symbol '%v' not found", symbol))
}

func asFloat(v interface{}) float64 {
	switch v := v.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	default:
		panic(fmt.Sprintf("cannot convert %v to float64", v))
	}
}

// BuiltInEnv creates a new default built-in namespace env.
func BuiltInEnv() *Env {
	return &Env{
		outer: nil,
		bindings: map[string]Sexpr{
			"+": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				var sum int64
				for _, arg := range args {
					sum += arg.Val.(int64)
				}
				return Sexpr{Type: "integer", Val: sum}
			}},
			"-": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) < 2 {
					panic("wrong number of arguments. `-` requires at least 2 arguments")
				}
				var diff int64
				for i, arg := range args {
					if i == 0 {
						diff = arg.Val.(int64)
					} else {
						diff -= arg.Val.(int64)
					}
				}
				return Sexpr{Type: "integer", Val: diff}
			}},
			"*": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				var product int64 = 1
				for _, arg := range args {
					product *= arg.Val.(int64)
				}
				return Sexpr{Type: "integer", Val: product}
			}},
			"/": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) < 1 {
					panic("wrong number of arguments. `/` requires at least 1 arguments")
				}
				var quotient int64 = 1
				for i, arg := range args {
					if i == 0 && len(args) > 1 {
						quotient = arg.Val.(int64)
					} else {
						quotient /= arg.Val.(int64)
					}
				}
				return Sexpr{Type: "integer", Val: quotient}
			}},
			"prn": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) < 1 {
					panic("wrong number of arguments. `prn` requires at least 1 arguments")
				}
				fmt.Println(PrintStr(args[0]))
				return Sexpr{Type: "nil", Val: nil}
			}},
			"list": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				return Sexpr{Type: "list", Val: args}
			}},
			"list?": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) > 0 && args[0].Type == "list" {
					return Sexpr{Type: "boolean", Val: true}
				}
				return Sexpr{Type: "boolean", Val: false}
			}},
			"empty?": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) > 0 && args[0].Type == "list" && len(args[0].Val.([]Sexpr)) == 0 {
					return Sexpr{Type: "boolean", Val: true}
				}
				return Sexpr{Type: "boolean", Val: false}
			}},
			"count": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) > 0 && args[0].Type == "list" {
					return Sexpr{Type: "integer", Val: int64(len(args[0].Val.([]Sexpr)))}
				}
				return Sexpr{Type: "integer", Val: int64(0)}
			}},
			"=": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `=` requires 2 arguments")
				}
				if args[0].Type == "list" && args[1].Type == "list" {
					alist := args[0].Val.([]Sexpr)
					blist := args[1].Val.([]Sexpr)
					if len(alist) != len(blist) {
						return Sexpr{Type: "boolean", Val: false}
					}
					for i, a := range alist {
						if a.Type != blist[i].Type || a.Val != blist[i].Val {
							return Sexpr{Type: "boolean", Val: false}
						}
					}
					return Sexpr{Type: "boolean", Val: true}
				}

				if args[0].Type != args[1].Type || args[0].Val != args[1].Val {
					return Sexpr{Type: "boolean", Val: false}
				}
				return Sexpr{Type: "boolean", Val: true}
			}},
			"<": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `<` requires 2 arguments")
				}
				if asFloat(args[0].Val) < asFloat(args[1].Val) {
					return Sexpr{Type: "boolean", Val: true}
				}
				return Sexpr{Type: "boolean", Val: false}
			}},
			"<=": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `<` requires 2 arguments")
				}
				if asFloat(args[0].Val) <= asFloat(args[1].Val) {
					return Sexpr{Type: "boolean", Val: true}
				}
				return Sexpr{Type: "boolean", Val: false}
			}},
			">": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `<` requires 2 arguments")
				}
				if asFloat(args[0].Val) > asFloat(args[1].Val) {
					return Sexpr{Type: "boolean", Val: true}
				}
				return Sexpr{Type: "boolean", Val: false}
			}},
			">=": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `<` requires 2 arguments")
				}
				if asFloat(args[0].Val) >= asFloat(args[1].Val) {
					return Sexpr{Type: "boolean", Val: true}
				}
				return Sexpr{Type: "boolean", Val: false}
			}},
			"read-string": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) != 1 {
					panic("wrong number of arguments. `read-string` requires 1 arguments")
				}
				return ReadStr(args[0].Val.(string))
			}},
			"slurp": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) != 1 {
					panic("wrong number of arguments. `slurp` requires 1 arguments")
				}
				s, err := os.ReadFile(args[0].Val.(string))
				if err != nil {
					panic(fmt.Sprintf("error reading file: %v", err))
				}
				return Sexpr{Type: "string", Val: string(s)}
			}},
			"eval": {Type: "function", Val: func(args ...Sexpr) Sexpr {
				if len(args) != 1 {
					panic("wrong number of arguments. `eval` requires 1 arguments")
				}
				return Eval(args[0], nil) // don't create infinite loop instantiating envs
			}},
		},
	}
}
