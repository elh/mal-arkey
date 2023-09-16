package malarkey

import (
	"regexp"
	"strconv"
	"strings"
)

var tokenRegex = regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)

// Reader reads tokens.
type Reader struct {
	Tokens   []string
	Position int
}

// Next returns the next token and advances the reader.
func (r *Reader) Next() string {
	if r.Position >= len(r.Tokens) {
		return ""
	}
	r.Position++
	return r.Tokens[r.Position-1]
}

// Peek returns the next token without advancing the reader.
func (r *Reader) Peek() string {
	if r.Position >= len(r.Tokens) {
		return ""
	}
	return r.Tokens[r.Position]
}

// Tokenize splits a input text into tokens.
func Tokenize(input string) []string {
	matches := tokenRegex.FindAllStringSubmatch(input, -1)
	var out []string
	for _, match := range matches {
		// note: regex does not drop leading whitespaces and commas. trim until no change
		var cur string
		trimmed := match[0]
		for trimmed != cur {
			cur = trimmed
			trimmed = strings.Trim(strings.TrimSpace(cur), ",")
		}
		if cur == "" || strings.HasPrefix(cur, ";") {
			continue
		}
		out = append(out, cur)
	}
	return out
}

// Read parses input text into an AST.
func Read(input string) Value {
	reader := &Reader{Tokens: Tokenize(input)}
	s := readForm(reader)
	if reader.Peek() != "" {
		panic("invalid trailing tokens")
	}
	return s
}

func readCollection(reader *Reader, peeked string) Value {
	stopToken := map[string]string{"(": ")", "[": "]", "{": "}"}[peeked]
	seqType := map[string]string{"(": "list", "[": "vector", "{": "hash-map"}[peeked]

	reader.Next()
	var elements []Value
	for reader.Peek() != stopToken {
		elements = append(elements, readForm(reader))
	}
	reader.Next()

	if seqType == "hash-map" {
		kv := map[string]Value{}
		for i := 0; i < len(elements); i += 2 {
			kv[elements[i].Val.(string)] = elements[i+1]
		}
		return Value{Type: "hash-map", Val: kv}
	}
	return Value{Type: seqType, Val: elements}
}

// only currently supporting integers and symbols
func readAtom(reader *Reader) Value {
	token := reader.Next()
	if token == "" {
		panic("expected atom")
	}

	if i, err := strconv.ParseInt(token, 10, 0); err == nil {
		return Value{Type: "integer", Val: i}
	}
	if f, err := strconv.ParseFloat(token, 64); err == nil {
		return Value{Type: "float", Val: f}
	}
	if strings.HasPrefix(token, "\"") {
		str := token[1 : len(token)-1]
		str = strings.Replace(str, "\\\"", "\"", -1)
		str = strings.Replace(str, "\\n", "\n", -1)
		str = strings.Replace(str, "\\\\", "\\", -1)
		return Value{Type: "string", Val: str}
	}
	if strings.HasPrefix(token, ":") {
		return Value{Type: "keyword", Val: token}
	}

	switch token {
	case "true":
		return Value{Type: "boolean", Val: true}
	case "false":
		return Value{Type: "boolean", Val: false}
	case "nil":
		return Value{Type: "nil", Val: nil}
	default:
		return Value{Type: "symbol", Val: token}
	}
}

// readForm parses the next form from the reader.
// Currently only supporting lists and atoms.
func readForm(reader *Reader) Value {
	peekToken := reader.Peek()
	switch peekToken {
	case "@", "'", "`", "~", "~@": // reader macros
		syms := map[string]string{"@": "deref", "'": "quote", "`": "quasiquote", "~": "unquote", "~@": "splice-unquote"}
		reader.Next()
		return Value{Type: "list", Val: []Value{{Type: "symbol", Val: syms[peekToken]}, readForm(reader)}}
	case "(", "[", "{":
		return readCollection(reader, peekToken)
	}
	return readAtom(reader)
}
