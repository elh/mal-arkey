package mal

import (
	"fmt"
	"strings"
)

// PrintStr returns a string representation of the given sexpr.
func PrintStr(s Sexpr, readably bool) string {
	switch s.Type {
	case "string":
		str := s.Val.(string)
		if readably {
			str = strings.Replace(str, "\\", "\\\\", -1)
			str = strings.Replace(str, "\"", "\\\"", -1)
			str = strings.Replace(str, "\n", "\\n", -1)
			return fmt.Sprintf("\"%s\"", str)
		}
		return str
	case "symbol", "integer", "float", "boolean":
		return fmt.Sprintf("%v", s.Val)
	case "nil":
		return "nil"
	case "list":
		var elements []string
		for _, element := range s.Val.([]Sexpr) {
			elements = append(elements, PrintStr(element, readably))
		}
		return fmt.Sprintf("(%s)", strings.Join(elements, " "))
	case "function", "function-tco":
		return "#<function>"
	default:
		panic(fmt.Sprintf("cannot print unsupported type: %s", s.Type))
	}
}
