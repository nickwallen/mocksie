package generator

import (
	"html/template"
	"os"

	"github.com/Masterminds/sprig"
	"github.com/nickwallen/mocksie/internal/parser"
)

type writer interface {
	Write(p []byte) (n int, err error)
}

// Generator generates the mock implementation of an Interface.
type Generator struct {
	writer writer
	tmpl   *template.Template
}

// NewGenerator create a new Generator.
func NewGenerator() (*Generator, error) {
	return &Generator{
		writer: os.Stdout,
		tmpl:   initTemplates(),
	}, nil
}

// GenerateMock generates a mock for an Interface.
func (g *Generator) GenerateMock(iface *parser.Interface) error {
	return g.tmpl.ExecuteTemplate(g.writer, "base", iface)
}

func initTemplates() *template.Template {
	tmpl := template.New("").Funcs(sprig.FuncMap())
	tmpl = template.Must(tmpl.New("base").Parse(baseTemplate))
	tmpl = template.Must(tmpl.New("methods").Parse(methodsTemplate))
	tmpl = template.Must(tmpl.New("declare-params").Parse(declareParamsTemplate))
	tmpl = template.Must(tmpl.New("use-params").Parse(useParamsTemplate))
	tmpl = template.Must(tmpl.New("results").Parse(resultsTemplate))
	return tmpl
}