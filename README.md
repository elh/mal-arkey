# mal-arkey - "Make a Lisp" in Go

> Mal is a Clojure inspired Lisp interpreter

Mal-arkey is my complete [mal](https://github.com/kanaka/mal) implementation in Go biasing towards readability and simplicity. It passes all non-optional tests and is capable of self-hosting mal.

```
> make repl
Mal [Mal-arkey]
user> (defmacro! when (fn* [test & body] `(if ~test (do ~@body))))
#<function>
user> (let* [name "Mal-arkey"] (when (not (nil? name)) (println "Begin" name)))
Begin Mal-arkey
nil
```

### Usage

* `make repl` to start a Mal-arkey REPL.
* `make bin` to build binaries that can be copied into the `kanaka/mal` test harness. My setup is at `elh/mal`.
