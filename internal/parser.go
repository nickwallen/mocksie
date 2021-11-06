package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// FileParser parses a file containing Go source code.
type FileParser struct {
	filename string
}

// NewFileParser constructs a new FileParser.
func NewFileParser(filename string) (*FileParser, error) {
	// Find the absolute path to the file
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	// Ensure the file exists
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}

	return &FileParser{filename: filename}, nil
}

// FindInterfaces returns all interfaces defined in the file.
func (p *FileParser) FindInterfaces() ([]*Interface, error) {
	// Parse the file
	set := token.NewFileSet()
	f, err := parser.ParseFile(set, p.filename, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	// Find any interfaces
	ifaces := make([]*Interface, 0)
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			if d.Tok != token.TYPE {
				continue
			}
			for _, spec := range d.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if iface, ok := spec.Type.(*ast.InterfaceType); ok {
						ifaces = append(ifaces, createInterface(spec.Name.String(), iface))
					}
				}
			}
		}
	}
	return ifaces, nil
}

func createInterface(name string, iface *ast.InterfaceType) *Interface {
	var methods []Method
	for _, method := range iface.Methods.List {
		if len(method.Names) == 0 {
			continue
		}
		methods = append(methods, Method{
			Name: method.Names[0].Name,
		})
	}
	return &Interface{
		Name:    name,
		Methods: methods,
	}
}
