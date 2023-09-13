package printer

import (
	"fmt"
	"strings"

	"github.com/elh/mal-go/internal/pkg/ast"
)

// PrintStr returns a string representation of the given sexpr.
func PrintStr(s ast.Sexpr) string {
	switch s.Type {
	case "symbol", "integer", "float", "boolean":
		return fmt.Sprintf("%v", s.Val)
	case "nil":
		return "nil"
	case "list":
		var elements []string
		for _, element := range s.Val.([]ast.Sexpr) {
			elements = append(elements, PrintStr(element))
		}
		return fmt.Sprintf("(%s)", strings.Join(elements, " "))
	case "function", "function-tco":
		return "#<function>"
	default:
		panic(fmt.Sprintf("cannot print unsupported type: %s", s.Type))
	}
}
