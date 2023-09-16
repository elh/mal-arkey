package mal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env is a map of symbols to bound values.
type Env struct {
	outer    *Env
	bindings map[string]Value
}

// NewEnv creates a new environment with the given outer environment.
func NewEnv(outer *Env, bindSymbols []Value, bindValues []Value) *Env {
	bindings := map[string]Value{}
	for i, s := range bindSymbols {
		if s.Type == "symbol" && s.Val.(string) == "&" {
			bindings[bindSymbols[i+1].Val.(string)] = Value{Type: "list", Val: bindValues[i:]}
			break
		}
		bindings[s.Val.(string)] = bindValues[i]
	}

	return &Env{
		outer:    outer,
		bindings: bindings,
	}
}

// Set binds a symbol to a value in the current environment.
func (e *Env) Set(symbol string, value Value) {
	e.bindings[symbol] = value
}

// Get returns the value bound to the given symbol in the environment.
func (e *Env) Get(symbol string) (Value, error) {
	if val, ok := e.bindings[symbol]; ok {
		return val, nil
	}
	if e.outer != nil {
		return e.outer.Get(symbol)
	}
	return Value{}, fmt.Errorf("'%v' not found", symbol)
}

// hack for numerical comparison
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

// get the underlying go function from a function Value
func getFn(expr Value) func(...Value) Value {
	if expr.Type == "function" {
		return expr.Val.(func(...Value) Value)
	} else if expr.Type == "function-tco" {
		return expr.Val.(FunctionTCO).Fn
	}
	return nil
}

// BuiltinEnv creates a new default built-in function env.
func BuiltinEnv() *Env {
	env := &Env{
		outer: nil,
		bindings: map[string]Value{
			"*host-language*": {Type: "string", Val: "Mal-arkey"},
			"+": {Type: "function", Val: func(args ...Value) Value {
				var sum int64
				for _, arg := range args {
					sum += arg.Val.(int64)
				}
				return Value{Type: "integer", Val: sum}
			}},
			"-": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("-", args, []string{"integer", "integer", "*"})
				var diff int64
				for i, arg := range args {
					if i == 0 {
						diff = arg.Val.(int64)
					} else {
						diff -= arg.Val.(int64)
					}
				}
				return Value{Type: "integer", Val: diff}
			}},
			"*": {Type: "function", Val: func(args ...Value) Value {
				var product int64 = 1
				for _, arg := range args {
					product *= arg.Val.(int64)
				}
				return Value{Type: "integer", Val: product}
			}},
			"/": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("/", args, []string{"integer", "*"})
				var quotient int64 = 1
				for i, arg := range args {
					if i == 0 && len(args) > 1 {
						quotient = arg.Val.(int64)
					} else {
						quotient /= arg.Val.(int64)
					}
				}
				return Value{Type: "integer", Val: quotient}
			}},
			"pr-str": {Type: "function", Val: func(args ...Value) Value {
				var strs []string
				for _, arg := range args {
					strs = append(strs, Print(arg, true))
				}
				return Value{Type: "string", Val: strings.Join(strs, " ")}
			}},
			"str": {Type: "function", Val: func(args ...Value) Value {
				var strs []string
				for _, arg := range args {
					strs = append(strs, Print(arg, false))
				}
				return Value{Type: "string", Val: strings.Join(strs, "")}
			}},
			"prn": {Type: "function", Val: func(args ...Value) Value {
				var strs []string
				for _, arg := range args {
					strs = append(strs, Print(arg, true))
				}
				fmt.Println(strings.Join(strs, " "))
				return Value{Type: "nil", Val: nil}
			}},
			"println": {Type: "function", Val: func(args ...Value) Value {
				var strs []string
				for _, arg := range args {
					strs = append(strs, Print(arg, false))
				}
				fmt.Println(strings.Join(strs, " "))
				return Value{Type: "nil", Val: nil}
			}},
			"list": {Type: "function", Val: func(args ...Value) Value {
				return Value{Type: "list", Val: args}
			}},
			"list?": {Type: "function", Val: func(args ...Value) Value {
				return Value{Type: "boolean", Val: len(args) > 0 && args[0].Type == "list"}
			}},
			"empty?": {Type: "function", Val: func(args ...Value) Value {
				return Value{Type: "boolean", Val: len(args) > 0 && args[0].Type == "list" && len(args[0].Val.([]Value)) == 0}
			}},
			"count": {Type: "function", Val: func(args ...Value) Value {
				if len(args) > 0 && args[0].Type == "list" {
					return Value{Type: "integer", Val: int64(len(args[0].Val.([]Value)))}
				}
				return Value{Type: "integer", Val: int64(0)}
			}},
			"=": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("=", args, []string{"any", "any"})
				if args[0].Type == "list" && args[1].Type == "list" {
					alist := args[0].Val.([]Value)
					blist := args[1].Val.([]Value)
					if len(alist) != len(blist) {
						return Value{Type: "boolean", Val: false}
					}
					for i, a := range alist {
						if a.Type != blist[i].Type || a.Val != blist[i].Val {
							return Value{Type: "boolean", Val: false}
						}
					}
					return Value{Type: "boolean", Val: true}
				}

				return Value{Type: "boolean", Val: args[0].Type == args[1].Type && args[0].Val == args[1].Val}
			}},
			"<": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("<", args, []string{"any", "any"})
				return Value{Type: "boolean", Val: asFloat(args[0].Val) < asFloat(args[1].Val)}
			}},
			"<=": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("<=", args, []string{"any", "any"})
				return Value{Type: "boolean", Val: asFloat(args[0].Val) <= asFloat(args[1].Val)}
			}},
			">": {Type: "function", Val: func(args ...Value) Value {
				validateArgs(">", args, []string{"any", "any"})
				return Value{Type: "boolean", Val: asFloat(args[0].Val) > asFloat(args[1].Val)}
			}},
			">=": {Type: "function", Val: func(args ...Value) Value {
				validateArgs(">=", args, []string{"any", "any"})
				return Value{Type: "boolean", Val: asFloat(args[0].Val) >= asFloat(args[1].Val)}
			}},
			"read-string": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("read-string", args, []string{"string"})
				return Read(args[0].Val.(string))
			}},
			"slurp": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("slurp", args, []string{"string"})
				s, err := os.ReadFile(args[0].Val.(string))
				if err != nil {
					panic(fmt.Sprintf("error reading file: %v", err))
				}
				return Value{Type: "string", Val: string(s)}
			}},
			"atom": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("atom", args, []string{"any"})
				atoms = append(atoms, args[0])
				return Value{Type: "atom", Val: len(atoms) - 1}
			}},
			"atom?": {Type: "function", Val: func(args ...Value) Value {
				return Value{Type: "boolean", Val: len(args) > 0 && args[0].Type == "atom"}
			}},
			"deref": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("deref", args, []string{"atom"})
				return atoms[args[0].Val.(int)]
			}},
			"reset!": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("reset!", args, []string{"atom", "any"})
				atoms[args[0].Val.(int)] = args[1]
				return args[1]
			}},
			"swap!": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("swap!", args, []string{"atom", "any", "*"})
				fn := getFn(args[1])
				if fn == nil {
					panic("second argument to `swap!` must be a function")
				}

				atomID := args[0].Val.(int)
				atomCur := atoms[atomID]

				fnArgs := append([]Value{atomCur}, args[2:]...)
				val := fn(fnArgs...)

				atoms[atomID] = val
				return val
			}},
			"cons": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("cons", args, []string{"any", "list"})
				return Value{Type: "list", Val: append([]Value{args[0]}, args[1].Val.([]Value)...)}
			}},
			"concat": {Type: "function", Val: func(args ...Value) Value {
				var vals []Value
				for _, arg := range args {
					if arg.Type != "list" {
						panic("all arguments to `concat` must be lists")
					}
					vals = append(vals, arg.Val.([]Value)...)
				}
				return Value{Type: "list", Val: vals}
			}},
			"nth": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("nth", args, []string{"list|vector", "integer"})
				list := args[0].Val.([]Value)
				idx := args[1].Val.(int64)
				if idx < 0 || idx >= int64(len(list)) {
					panic("index out of bounds")
				}
				return list[idx]
			}},
			"first": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("first", args, []string{"nil|list|vector"})
				if args[0].Type == "nil" {
					return Value{Type: "nil", Val: nil}
				}
				list := args[0].Val.([]Value)
				if len(list) == 0 {
					return Value{Type: "nil", Val: nil}
				}
				return list[0]
			}},
			"rest": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("rest", args, []string{"nil|list|vector"})
				if args[0].Type == "nil" {
					return Value{Type: "list", Val: []Value{}}
				}
				list := args[0].Val.([]Value)
				if len(list) == 0 {
					return Value{Type: "list", Val: []Value{}}
				}
				return Value{Type: "list", Val: list[1:]}
			}},
			"throw": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("throw", args, []string{"any"})
				panic(args[0])
			}},
			"apply": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("apply", args, []string{"function|function-tco", "any", "*"})
				if args[len(args)-1].Type != "list" {
					panic("last argument to `apply` must be a list")
				}

				fn := getFn(args[0])
				if fn == nil {
					panic("first argument to `apply` must be a function")
				}
				fnArgs := append(args[1:len(args)-1], args[len(args)-1].Val.([]Value)...)
				return fn(fnArgs...)
			}},
			"map": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("map", args, []string{"function|function-tco", "list|vector"})
				fn := getFn(args[0])
				if fn == nil {
					panic("first argument to `map` must be a function")
				}

				var res []Value
				for _, arg := range args[1].Val.([]Value) {
					res = append(res, fn(arg))
				}
				return Value{Type: "list", Val: res}
			}},
			"nil?": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("nil?", args, []string{"any"})
				return Value{Type: "boolean", Val: args[0].Type == "nil"}
			}},
			"true?": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("true?", args, []string{"any"})
				return Value{Type: "boolean", Val: args[0].Type == "boolean" && args[0].Val.(bool)}
			}},
			"false?": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("false?", args, []string{"any"})
				return Value{Type: "boolean", Val: args[0].Type == "boolean" && !args[0].Val.(bool)}
			}},
			"symbol": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("symbol", args, []string{"string"})
				return Value{Type: "symbol", Val: args[0].Val.(string)}
			}},
			"symbol?": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("symbol?", args, []string{"any"})
				return Value{Type: "boolean", Val: args[0].Type == "symbol"}
			}},
			"keyword": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("keyword", args, []string{"string|keyword"})
				if args[0].Type == "string" {
					return Value{Type: "keyword", Val: ":" + args[0].Val.(string)}
				}
				return args[0]
			}},
			"keyword?": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("keyword?", args, []string{"any"})
				return Value{Type: "boolean", Val: args[0].Type == "keyword"}
			}},
			"sequential?": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("sequential?", args, []string{"any"})
				return Value{Type: "boolean", Val: args[0].Type == "list" || args[0].Type == "vector"}
			}},
			"vec": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("vec", args, []string{"list|vector"})
				return Value{Type: "vector", Val: args[0].Val.([]Value)}
			}},
			"vector": {Type: "function", Val: func(args ...Value) Value {
				return Value{Type: "vector", Val: args}
			}},
			"vector?": {Type: "function", Val: func(args ...Value) Value {
				return Value{Type: "boolean", Val: len(args) > 0 && args[0].Type == "vector"}
			}},
			"hash-map": {Type: "function", Val: func(args ...Value) Value {
				if len(args)%2 != 0 {
					panic("wrong number of arguments. `hash-map` requires an even number of arguments")
				}
				kv := map[string]Value{}
				for i := 0; i < len(args); i += 2 {
					kv[args[i].Val.(string)] = args[i+1]
				}

				return Value{Type: "hash-map", Val: kv}
			}},
			"map?": {Type: "function", Val: func(args ...Value) Value {
				return Value{Type: "boolean", Val: len(args) > 0 && args[0].Type == "hash-map"}
			}},
			"assoc": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("assoc", args, []string{"hash-map", "any", "any", "*"})
				kv := args[0].Val.(map[string]Value)
				for i := 1; i < len(args)-1; i += 2 {
					kv[args[i].Val.(string)] = args[i+1]
				}
				return Value{Type: "hash-map", Val: kv}
			}},
			"dissoc": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("dissoc", args, []string{"hash-map", "list"})
				kv := args[0].Val.(map[string]Value)
				for _, arg := range args[1].Val.([]Value) {
					delete(kv, arg.Val.(string))
				}
				return Value{Type: "hash-map", Val: kv}
			}},
			"keys": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("keys", args, []string{"hash-map"})
				var keys []Value
				for k := range args[0].Val.(map[string]Value) {
					keys = append(keys, Value{Type: "string", Val: k})
				}
				return Value{Type: "list", Val: keys}
			}},
			"vals": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("vals", args, []string{"hash-map"})
				var values []Value
				for _, v := range args[0].Val.(map[string]Value) {
					values = append(values, v)
				}
				return Value{Type: "list", Val: values}
			}},
			"get": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("get", args, []string{"hash-map", "any"})
				kv := args[0].Val.(map[string]Value)
				for i := 1; i < len(args); i++ {
					if val, ok := kv[args[i].Val.(string)]; ok {
						return val
					}
				}
				return Value{Type: "nil", Val: nil}
			}},
			"contains?": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("contains?", args, []string{"hash-map", "any"})
				kv := args[0].Val.(map[string]Value)
				_, ok := kv[args[1].Val.(string)]
				return Value{Type: "boolean", Val: ok}
			}},
			"readline": {Type: "function", Val: func(args ...Value) Value {
				validateArgs("readline", args, []string{"string"})
				reader := bufio.NewReader(os.Stdin)

				fmt.Print(args[0].Val.(string))
				input, err := reader.ReadString('\n')
				if err != nil {
					if err.Error() == "EOF" {
						return Value{Type: "nil", Val: nil}
					}
					panic(err)
				}
				return Value{Type: "string", Val: strings.TrimRight(input, "\n")}
			}},
			"time-ms":   {Type: "function", Val: func(args ...Value) Value { panic("unimplemented") }},
			"meta":      {Type: "function", Val: func(args ...Value) Value { panic("unimplemented") }},
			"with-meta": {Type: "function", Val: func(args ...Value) Value { panic("unimplemented") }},
			"fn?":       {Type: "function", Val: func(args ...Value) Value { panic("unimplemented") }},
			"string?":   {Type: "function", Val: func(args ...Value) Value { panic("unimplemented") }},
			"number?":   {Type: "function", Val: func(args ...Value) Value { panic("unimplemented") }},
			"seq":       {Type: "function", Val: func(args ...Value) Value { panic("unimplemented") }},
			"conj":      {Type: "function", Val: func(args ...Value) Value { panic("unimplemented") }},
		},
	}

	// defined here to allow cyclic reference to env
	env.bindings["eval"] = Value{Type: "function", Val: func(args ...Value) Value {
		validateArgs("eval", args, []string{"any"})
		return Eval(args[0], env)
	}}

	return env
}
