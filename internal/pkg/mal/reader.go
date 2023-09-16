package mal

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
		cur := match[0]
		// note: regex does not drop leading whitespaces and commas
		for {
			trimmed := strings.Trim(strings.TrimSpace(cur), ",")
			if trimmed == cur {
				break
			}
			cur = trimmed
		}
		if cur == "" || strings.HasPrefix(cur, ";") {
			continue
		}
		out = append(out, cur)
	}
	return out
}

// ReadStr parses input text into a sexpr.
func ReadStr(input string) Value {
	reader := &Reader{Tokens: Tokenize(input)}
	s := ReadForm(reader)
	if reader.Peek() != "" {
		panic("invalid trailing tokens")
	}
	return s
}

func readSequence(reader *Reader, start string) Value {
	stopToken := map[string]string{"(": ")", "[": "]"}[start]
	seqType := map[string]string{"(": "list", "[": "vector"}[start]

	reader.Next()

	var elements []Value
	for reader.Peek() != stopToken {
		elements = append(elements, ReadForm(reader))
	}
	reader.Next()

	return Value{Type: seqType, Val: elements}
}

func readHashMap(reader *Reader) Value {
	if reader.Peek() != "{" {
		panic("expected '{'")
	}
	reader.Next()

	var elements []Value
	for reader.Peek() != "}" {
		elements = append(elements, ReadForm(reader))
	}
	reader.Next()

	kv := map[string]Value{}
	for i := 0; i < len(elements); i += 2 {
		kv[elements[i].Val.(string)] = elements[i+1]
	}

	return Value{Type: "hash-map", Val: kv}
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

// ReadForm parses the next form from the reader.
// Currently only supporting lists and atoms.
func ReadForm(reader *Reader) Value {
	peekToken := reader.Peek()
	switch peekToken {
	case "@":
		reader.Next()
		return Value{Type: "list", Val: []Value{
			{Type: "symbol", Val: "deref"},
			ReadForm(reader),
		}}
	case "'":
		reader.Next()
		return Value{Type: "list", Val: []Value{
			{Type: "symbol", Val: "quote"},
			ReadForm(reader),
		}}
	case "`":
		reader.Next()
		return Value{Type: "list", Val: []Value{
			{Type: "symbol", Val: "quasiquote"},
			ReadForm(reader),
		}}
	case "~":
		reader.Next()
		return Value{Type: "list", Val: []Value{
			{Type: "symbol", Val: "unquote"},
			ReadForm(reader),
		}}
	case "~@":
		reader.Next()
		return Value{Type: "list", Val: []Value{
			{Type: "symbol", Val: "splice-unquote"},
			ReadForm(reader),
		}}
	case "(", "[":
		return readSequence(reader, peekToken)
	case "{":
		return readHashMap(reader)
	}
	return readAtom(reader)
}
