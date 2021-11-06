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
						ifaces = append(ifaces, buildInterface(spec.Name.String(), iface))
					}
				}
			}
		}
	}
	return ifaces, nil
}

func buildInterface(name string, iface *ast.InterfaceType) *Interface {
	var methods []Method
	for _, method := range iface.Methods.List {
		// Expect a func type
		if _, ok := method.Type.(*ast.FuncType); !ok {
			continue
		}

		// Expect the method to be named
		if len(method.Names) == 0 {
			continue
		}

		// Build the method
		funcType := method.Type.(*ast.FuncType)
		methods = append(methods, Method{
			Name:    method.Names[0].Name,
			Params:  buildParams(funcType),
			Results: buildResults(funcType),
		})
	}
	return &Interface{
		Name:    name,
		Methods: methods,
	}
}

func buildResults(funcType *ast.FuncType) []Result {
	results := make([]Result, 0)
	for i := range funcType.Results.List {
		field := funcType.Results.List[i]

		// Expect the field to be an identifier
		if _, ok := field.Type.(*ast.Ident); !ok {
			continue
		}

		// The result may not be named
		name := ""
		if len(field.Names) > 0 {
			name = field.Names[0].Name
		}

		// Build the result
		results = append(results, Result{
			Name: name,
			Type: field.Type.(*ast.Ident).Name,
		})
	}
	return results
}

func buildParams(funcType *ast.FuncType) []Param {
	params := make([]Param, 0)
	for i := range funcType.Params.List {
		field := funcType.Params.List[i]

		// Expect the field to be an identifier
		if _, ok := field.Type.(*ast.Ident); !ok {
			continue
		}

		// The param may not be named
		name := ""
		if len(field.Names) > 0 {
			name = field.Names[0].Name
		}

		// Build the param
		params = append(params, Param{
			Name: name,
			Type: field.Type.(*ast.Ident).Name,
		})
	}
	return params
}
