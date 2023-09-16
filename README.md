# mal-arkey - "Make a Lisp" in Go

[![CI](https://github.com/elh/mal-arkey/actions/workflows/ci.yaml/badge.svg)](https://github.com/elh/mal-arkey/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/elh/mal-arkey)](https://goreportcard.com/report/github.com/elh/mal-arkey)

> Mal is a Clojure inspired Lisp interpreter

Mal-arkey is a complete [Mal](https://github.com/kanaka/mal) implementation in Go aiming for readability and simplicity. It passes all non-optional tests and can self-host Mal.

```
> make repl
Mal [Mal-arkey]
user> (defmacro! when (fn* [test & body] `(if ~test (do ~@body))))
#<function>
user> (let* [name "Mal-arkey"] (when (not (nil? name)) (println "Begin" (str name "!"))))
Begin Mal-arkey!
nil
```

### Usage

* `make repl` to start a Mal-arkey REPL.
* `make bin` to build binaries that can be copied into the `kanaka/mal` test harness. My setup is at `elh/mal`.
