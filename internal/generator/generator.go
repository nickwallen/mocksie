package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/nickwallen/mocksie/internal"
)

const (
	structTemplate = `
type mock%s struct {
}
`
)

type writer interface {
	Write(p []byte) (n int, err error)
}

// Generator generates the mock implementation of an Interface.
type Generator struct {
	writer writer
}

// NewGenerator create a new Generator.
func NewGenerator() *Generator {
	return &Generator{
		writer: os.Stdout,
	}
}

func (g *Generator) GenerateMock(iface *parser.Interface) error {
	_, err := fmt.Fprintf(g.writer, structTemplate, strings.Title(iface.Name))
	return err
}
