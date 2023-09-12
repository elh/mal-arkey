package reader

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/elh/mal-go/internal/pkg/ast"
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
		if cur == "" {
			continue
		}
		out = append(out, cur)
	}
	return out
}

// ReadStr parses input text into a sexpr.
func ReadStr(input string) ast.Sexpr {
	reader := &Reader{Tokens: Tokenize(input)}
	s := ReadForm(reader)
	if reader.Peek() != "" {
		panic("invalid trailing tokens")
	}
	return s
}

func readList(reader *Reader) ast.Sexpr {
	if reader.Peek() != "(" {
		panic("expected '('")
	}
	reader.Next()

	elements := []ast.Sexpr{}
	for reader.Peek() != ")" {
		elements = append(elements, ReadForm(reader))
	}
	reader.Next()

	return ast.Sexpr{
		Type: "list",
		Val:  elements,
	}
}

// only currently supporting integers and symbols
func readAtom(reader *Reader) ast.Sexpr {
	token := reader.Next()
	if token == "" {
		panic("expected atom")
	}
	i, err := strconv.Atoi(token)
	if err != nil {
		return ast.Sexpr{
			Type: "symbol",
			Val:  token,
		}
	}
	return ast.Sexpr{
		Type: "integer",
		Val:  i,
	}
}

// ReadForm parses the next sexpr from the reader.
// Currently only supporting lists and atoms.
func ReadForm(reader *Reader) ast.Sexpr {
	if reader.Peek() == "(" {
		return readList(reader)
	}
	return readAtom(reader)
}
