package generator

import (
	"html/template"
	"os"

	"github.com/Masterminds/sprig"
	"github.com/nickwallen/mocksie/internal"
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

// GenerateMock generates a mock for an Interface.
func (g *Generator) GenerateMock(iface *parser.Interface) error {
	t, err := template.New("mock").Funcs(sprig.FuncMap()).ParseGlob("templates/*.template")
	if err != nil {
		return err
	}
	return t.Execute(g.writer, iface)
}
