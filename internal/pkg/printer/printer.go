package printer

import (
	"fmt"
	"strings"

	"github.com/elh/mal-go/internal/pkg/ast"
)

func PrintStr(s ast.Sexpr) string {
	switch s.Type {
	case "symbol":
		return s.Val.(string)
	case "integer":
		return fmt.Sprintf("%d", s.Val)
	case "list":
		var elements []string
		for _, element := range s.Val.([]ast.Sexpr) {
			elements = append(elements, PrintStr(element))
		}
		return fmt.Sprintf("(%s)", strings.Join(elements, " "))
	default:
		panic(fmt.Sprintf("cannot print unsupported type: %s", s.Type))
	}
}