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
		// TODO: why isnt the regex just skipping whitespace and commas
		cur := match[0]
		for {
			trimmed := strings.TrimSpace(cur)
			trimmed = strings.Trim(trimmed, ",")
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
func ReadStr(input string) Sexpr {
	reader := &Reader{Tokens: Tokenize(input)}
	s := ReadForm(reader)
	if reader.Peek() != "" {
		panic("invalid trailing tokens")
	}
	return s
}

func readList(reader *Reader) Sexpr {
	if reader.Peek() != "(" {
		panic("expected '('")
	}
	reader.Next()

	elements := []Sexpr{}
	for reader.Peek() != ")" {
		elements = append(elements, ReadForm(reader))
	}
	reader.Next()

	return Sexpr{
		Type: "list",
		Val:  elements,
	}
}

func readVector(reader *Reader) Sexpr {
	if reader.Peek() != "[" {
		panic("expected '['")
	}
	reader.Next()

	elements := []Sexpr{}
	for reader.Peek() != "]" {
		elements = append(elements, ReadForm(reader))
	}
	reader.Next()

	return Sexpr{
		Type: "vector",
		Val:  elements,
	}
}

func readHashMap(reader *Reader) Sexpr {
	if reader.Peek() != "{" {
		panic("expected '{'")
	}
	reader.Next()

	elements := []Sexpr{}
	for reader.Peek() != "}" {
		elements = append(elements, ReadForm(reader))
	}
	reader.Next()

	kv := map[string]Sexpr{}
	for i := 0; i < len(elements); i += 2 {
		kv[elements[i].Val.(string)] = elements[i+1]
	}

	return Sexpr{
		Type: "hash-map",
		Val:  kv,
	}
}

// only currently supporting integers and symbols
func readAtom(reader *Reader) Sexpr {
	token := reader.Next()
	if token == "" {
		panic("expected atom")
	}

	if i, err := strconv.ParseInt(token, 10, 0); err == nil {
		return Sexpr{
			Type: "integer",
			Val:  i,
		}
	}
	if f, err := strconv.ParseFloat(token, 64); err == nil {
		return Sexpr{
			Type: "float",
			Val:  f,
		}
	}
	if strings.HasPrefix(token, "\"") {
		str := token[1 : len(token)-1]
		str = strings.Replace(str, "\\\"", "\"", -1)
		str = strings.Replace(str, "\\n", "\n", -1)
		str = strings.Replace(str, "\\\\", "\\", -1)
		return Sexpr{
			Type: "string",
			Val:  str,
		}
	}

	switch token {
	case "true":
		return Sexpr{
			Type: "boolean",
			Val:  true,
		}
	case "false":
		return Sexpr{
			Type: "boolean",
			Val:  false,
		}
	case "nil":
		return Sexpr{
			Type: "nil",
			Val:  nil,
		}
	default:
		return Sexpr{
			Type: "symbol",
			Val:  token,
		}
	}
}

// ReadForm parses the next sexpr from the reader.
// Currently only supporting lists and atoms.
func ReadForm(reader *Reader) Sexpr {
	peekToken := reader.Peek()
	switch peekToken {
	case "@":
		reader.Next()
		return Sexpr{Type: "list", Val: []Sexpr{
			{Type: "symbol", Val: "deref"},
			ReadForm(reader),
		}}
	case "'":
		reader.Next()
		return Sexpr{Type: "list", Val: []Sexpr{
			{Type: "symbol", Val: "quote"},
			ReadForm(reader),
		}}
	case "`":
		reader.Next()
		return Sexpr{Type: "list", Val: []Sexpr{
			{Type: "symbol", Val: "quasiquote"},
			ReadForm(reader),
		}}
	case "~":
		reader.Next()
		return Sexpr{Type: "list", Val: []Sexpr{
			{Type: "symbol", Val: "unquote"},
			ReadForm(reader),
		}}
	case "~@":
		reader.Next()
		return Sexpr{Type: "list", Val: []Sexpr{
			{Type: "symbol", Val: "splice-unquote"},
			ReadForm(reader),
		}}
	case "(":
		return readList(reader)
	case "[":
		return readVector(reader)
	case "{":
		return readHashMap(reader)
	}
	return readAtom(reader)
}
