package mal

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// Note: this recommended factoring doesn't click with me. this is like "eval non-function"?
func evalAST(sexpr Value, env *Env) Value {
	switch sexpr.Type {
	case "list", "vector":
		var elems []Value
		for _, elem := range sexpr.Val.([]Value) {
			elems = append(elems, Eval(elem, env))
		}
		return Value{Type: sexpr.Type, Val: elems}
	case "hash-map":
		kv := map[string]Value{}
		for k, v := range sexpr.Val.(map[string]Value) {
			kv[k] = Eval(v, env)
		}
		return Value{Type: "hash-map", Val: kv}
	case "symbol":
		s, err := env.Get(sexpr.Val.(string))
		if err != nil {
			panic(err)
		}
		return s
	default:
		return sexpr
	}
}

// panic if args are not of expected length and types. spec types can be "|"-separated unions.
func validateArgs(fn string, args []Value, spec []string) {
	var isVariadic bool
	if spec[len(spec)-1] == "*" {
		spec = spec[:len(spec)-1]
		isVariadic = true
	}
	if isVariadic && len(args) < len(spec) {
		panic(fmt.Sprintf("%s requires at least %d argument(s)", fn, len(spec)))
	} else if !isVariadic && len(args) != len(spec) {
		panic(fmt.Sprintf("%s requires %d argument(s)", fn, len(spec)))
	}
	for i, s := range spec {
		if s != "any" && !slices.Contains(strings.Split(s, "|"), args[i].Type) {
			panic(fmt.Sprintf("%s %d-idx argument must be of %s type", fn, i, s))
		}
	}
}

func evalDef(args []Value, env *Env) Value {
	validateArgs("def!", args, []string{"symbol", "any"})
	v := Eval(args[1], env)
	env.Set(args[0].Val.(string), v)
	return v
}

// With TCO. Return unevaluated body and new environment.
func evalLet(args []Value, env *Env) (Value, *Env) {
	validateArgs("let*", args, []string{"list|vector", "any"})
	letEnv := NewEnv(env, nil, nil)
	bindings := args[0].Val.([]Value)
	if len(bindings)%2 != 0 {
		panic("let* requires an even number of forms in bindings")
	}
	for i := 0; i < len(bindings); i += 2 {
		if bindings[i].Type != "symbol" {
			panic("let* bindings must be symbols")
		}
		letEnv.Set(bindings[i].Val.(string), Eval(bindings[i+1], letEnv))
	}

	return args[1], letEnv
}

// With TCO. Return unevaluated if/else branch form
func evalIf(args []Value, env *Env) Value {
	if len(args) != 2 && len(args) != 3 {
		panic("if requires three (or two) arguments")
	}
	cond := Eval(args[0], env)
	if (cond.Type == "boolean" && !cond.Val.(bool)) || (cond.Type == "nil") {
		if len(args) == 3 {
			return args[2]
		}
		return Value{Type: "nil", Val: nil}
	}
	return args[1]
}

// With TCO. Return unevaluated final form.
func evalDo(args []Value, env *Env) Value {
	for _, arg := range args[:len(args)-1] {
		Eval(arg, env)
	}
	return args[len(args)-1]
}

// With TCO. Return a function-tco value.
func evalFn(evalArgs []Value, env *Env) Value {
	validateArgs("fn*", evalArgs, []string{"list|vector", "any"})
	params := evalArgs[0].Val.([]Value)
	for _, param := range params {
		if param.Type != "symbol" {
			panic("fn* parameters must be symbols")
		}
	}
	body := evalArgs[1]

	return Value{Type: "function-tco", Val: FunctionTCO{
		AST:    body,
		Params: params,
		Env:    env,
		Fn: func(args ...Value) Value {
			fnEnv := NewEnv(env, params, args)
			return Eval(body, fnEnv)
		},
		IsMacro: false},
	}
}

func quasiquote(ast Value) Value {
	if ast.Type == "list" {
		elems := ast.Val.([]Value)
		if len(elems) == 2 && elems[0].Type == "symbol" && elems[0].Val.(string) == "unquote" {
			return elems[1]
		}
		if len(elems) == 0 {
			return ast
		}

		elem := elems[0]
		if elem.Type == "list" {
			children := elem.Val.([]Value)
			if len(children) > 0 && children[0].Type == "symbol" && children[0].Val.(string) == "splice-unquote" {
				return Value{Type: "list", Val: []Value{
					{Type: "symbol", Val: "concat"},
					children[1],
					quasiquote(Value{Type: "list", Val: elems[1:]})}}
			}
		}
		return Value{Type: "list", Val: []Value{
			{Type: "symbol", Val: "cons"},
			quasiquote(elem),
			quasiquote(Value{Type: "list", Val: elems[1:]}),
		}}
	}
	if ast.Type == "symbol" || ast.Type == "hash-map" {
		return Value{Type: "list", Val: []Value{
			{Type: "symbol", Val: "quote"},
			ast,
		}}
	}
	return ast
}

func evalDefMacro(args []Value, env *Env) Value {
	validateArgs("defmacro!", args, []string{"symbol", "any"})
	v := Eval(args[1], env)
	if v.Type != "function-tco" {
		panic("defmacro! requires a macro fn as second argument")
	}
	// need to re-wrap the non-ptr Value to update IsMacro = true
	f := v.Val.(FunctionTCO)
	f.IsMacro = true
	env.Set(args[0].Val.(string), Value{Type: "function-tco", Val: f})
	return v
}

func isMacroCall(ast Value, env *Env) bool {
	if ast.Type != "list" {
		return false
	}
	list := ast.Val.([]Value)
	if len(list) == 0 {
		return false
	}
	if list[0].Type != "symbol" {
		return false
	}
	symbol := list[0].Val.(string)
	v, err := env.Get(symbol)
	if err != nil {
		return false
	}
	if v.Type != "function-tco" {
		return false
	}
	return v.Val.(FunctionTCO).IsMacro
}

func macroExpand(ast Value, env *Env) Value {
	for isMacroCall(ast, env) {
		elems := ast.Val.([]Value)
		symbol := elems[0].Val.(string)
		macro, err := env.Get(symbol)
		if err != nil {
			continue
		}
		ast = macro.Val.(FunctionTCO).Fn(elems[1:]...)
	}
	return ast
}

func try(expr Value, env *Env) (value, exceptionValue *Value) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case string:
				exceptionValue = &Value{Type: "string", Val: v}
			case error:
				exceptionValue = &Value{Type: "string", Val: v.Error()}
			case Value:
				exceptionValue = &v
			}
		}
	}()
	out := Eval(expr, env)
	return &out, nil
}

func evalTryCatch(args []Value, env *Env) Value {
	validateArgs("try*", args, []string{"any", "list"})
	catchForm := args[1].Val.([]Value)
	validateArgs("catch*", catchForm, []string{"symbol", "symbol", "any"})
	if catchForm[0].Val.(string) != "catch*" {
		panic("try* requires a symbol 'catch* as first element of second argument")
	}

	value, exception := try(args[0], env)
	if exception == nil {
		return *value
	}
	catchEnv := NewEnv(env, []Value{catchForm[1]}, []Value{*exception})
	return Eval(catchForm[2], catchEnv)
}

// Eval evaluates an expression in the given environment.
func Eval(expr Value, env *Env) Value {
	// Tail call optimization prevents nested function calls.
	for {
		if expr.Type != "list" {
			return evalAST(expr, env)
		}
		if len(expr.Val.([]Value)) == 0 {
			return expr
		}

		// macro expansion
		expr = macroExpand(expr, env)
		if expr.Type != "list" {
			return evalAST(expr, env)
		}
		list := expr.Val.([]Value)

		// special forms
		if list[0].Type == "symbol" {
			args := list[1:]
			switch list[0].Val.(string) {
			case "def!":
				return evalDef(args, env)
			case "defmacro!":
				expr = evalDefMacro(args, env)
				continue
			case "let*":
				expr, env = evalLet(args, env)
				continue
			case "if":
				expr = evalIf(args, env)
				continue
			case "do":
				expr = evalDo(args, env)
				continue
			case "fn*":
				return evalFn(args, env)
			case "quote":
				return args[0]
			case "quasiquote":
				expr = quasiquote(args[0])
				continue
			case "macroexpand":
				return macroExpand(args[0], env)
			case "try*":
				return evalTryCatch(args, env)
			}
		}

		// function call
		evaluatedList := evalAST(expr, env)
		elems := evaluatedList.Val.([]Value)
		switch elems[0].Type {
		case "function":
			fn := elems[0].Val.(func(args ...Value) Value)
			return fn(elems[1:]...)
		case "function-tco":
			args := elems[1:]
			fn := elems[0].Val.(FunctionTCO)

			expr = fn.AST
			env = NewEnv(fn.Env, fn.Params, args)
			continue
		default:
			panic("first element of list must be a function")
		}
	}
}
