package mal

import (
	"fmt"
	"strings"
)

// PrintStr returns a string representation of the given sexpr.
func PrintStr(s Sexpr) string {
	switch s.Type {
	case "symbol", "integer", "float", "boolean":
		return fmt.Sprintf("%v", s.Val)
	case "nil":
		return "nil"
	case "list":
		var elements []string
		for _, element := range s.Val.([]Sexpr) {
			elements = append(elements, PrintStr(element))
		}
		return fmt.Sprintf("(%s)", strings.Join(elements, " "))
	case "function", "function-tco":
		return "#<function>"
	default:
		panic(fmt.Sprintf("cannot print unsupported type: %s", s.Type))
	}
}
