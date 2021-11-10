package generator

import (
	"html/template"
	"io"

	"github.com/Masterminds/sprig"
	"github.com/nickwallen/mocksie/internal"
)

// Generator generates the mock implementation of an Interface.
type Generator struct {
	writer io.Writer
	tmpl   *template.Template
}

// New create a new Generator.
func New(writer io.Writer) (*Generator, error) {
	return &Generator{
		writer: writer,
		tmpl:   initTemplates(),
	}, nil
}

// GenerateMock generates a mock for an Interface.
func (g *Generator) GenerateMock(iface *mocksie.Interface) error {
	return g.tmpl.ExecuteTemplate(g.writer, "base", iface)
}

// initTemplates initialize the templates that are used to generate the mocks.
func initTemplates() *template.Template {
	tmpl := template.New("").Funcs(sprig.FuncMap())
	tmpl = template.Must(tmpl.New("base").Parse(baseTemplate))
	tmpl = template.Must(tmpl.New("imports").Parse(importsTemplate))
	tmpl = template.Must(tmpl.New("methods").Parse(methodsTemplate))
	tmpl = template.Must(tmpl.New("declare-params").Parse(declareParamsTemplate))
	tmpl = template.Must(tmpl.New("use-params").Parse(useParamsTemplate))
	tmpl = template.Must(tmpl.New("results").Parse(resultsTemplate))
	return tmpl
}
