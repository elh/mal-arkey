# mal-arkey - my "Make a Lisp" impl in Go

> 1. Mal is a Clojure inspired Lisp interpreter
> 2. Mal is a learning tool

```
> make repl
Mal [Mal-arkey (go)]
user> (defmacro! when (fn* [test & body] `(if ~test (do ~@body))))
#<function>
user> (let* [name "Malarkey"] (when (not (nil? name)) (println "Begin" name)))
Begin Malarkey
nil
```

Status: This passes all non-optional, non-deferred tests in `kanaka/mal`. Currently, I'm working on getting it self-hosting.

### Usage

* `make repl` to start a Mal-arkey REPL.
* `make bin` to build binaries that can be copied into the `kanaka/mal` test harness. My setup is at `elh/mal`.
