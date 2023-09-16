package mal

import (
	"fmt"
	"strings"
)

// PrintStr returns a string representation of the given value.
func PrintStr(s Value, readably bool) string {
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
	case "symbol", "integer", "float", "boolean", "keyword":
		return fmt.Sprintf("%v", s.Val)
	case "nil":
		return "nil"
	case "list", "vector":
		var elements []string
		for _, element := range s.Val.([]Value) {
			elements = append(elements, PrintStr(element, readably))
		}
		if s.Type == "list" {
			return fmt.Sprintf("(%s)", strings.Join(elements, " "))
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, " "))
	case "hash-map":
		var elements []string
		for k, v := range s.Val.(map[string]Value) {
			elements = append(elements, fmt.Sprintf("\"%s\"", k), PrintStr(v, readably))
		}
		return fmt.Sprintf("{%s}", strings.Join(elements, " "))
	case "function", "function-tco":
		return "#<function>"
	case "atom":
		id := s.Val.(int)
		return fmt.Sprintf("(atom %s)", PrintStr(atoms[id], readably))
	default:
		panic(fmt.Sprintf("cannot print unsupported type: %s", s.Type))
	}
}
