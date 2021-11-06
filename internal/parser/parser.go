package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	parser2 "github.com/nickwallen/mocksie/internal"
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
func (p *FileParser) FindInterfaces() ([]*parser2.Interface, error) {
	// Parse the file
	set := token.NewFileSet()
	f, err := parser.ParseFile(set, p.filename, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	// Find any interfaces
	ifaces := make([]*parser2.Interface, 0)
	for _, decl := range f.Decls {
		if _, ok := decl.(*ast.GenDecl); !ok {
			continue // Not a declaration
		}
		if decl.(*ast.GenDecl).Tok != token.TYPE {
			continue // Not declaring a type
		}
		for _, spec := range decl.(*ast.GenDecl).Specs {
			if _, ok := spec.(*ast.TypeSpec); !ok {
				continue // Not a proper type
			}
			if _, ok := spec.(*ast.TypeSpec).Type.(*ast.InterfaceType); !ok {
				continue // Not an interface
			}
			typ := spec.(*ast.TypeSpec)
			ifaces = append(ifaces, &parser2.Interface{
				Name:    typ.Name.String(),
				Methods: buildMethods(typ.Type.(*ast.InterfaceType)),
			})
		}
	}
	return ifaces, nil
}

func buildMethods(typ *ast.InterfaceType) []parser2.Method {
	var methods []parser2.Method
	for _, field := range typ.Methods.List {
		// Expect a function type
		if _, ok := field.Type.(*ast.FuncType); !ok {
			continue
		}

		// Expect the method to be named
		if len(field.Names) == 0 {
			continue
		}

		// Build the method
		funcType := field.Type.(*ast.FuncType)
		methods = append(methods, parser2.Method{
			Name:    field.Names[0].Name,
			Params:  buildParams(funcType),
			Results: buildResults(funcType),
		})
	}
	return methods
}

func buildResults(funcType *ast.FuncType) []parser2.Result {
	results := make([]parser2.Result, 0)
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
		results = append(results, parser2.Result{
			Name: name,
			Type: field.Type.(*ast.Ident).Name,
		})
	}
	return results
}

func buildParams(funcType *ast.FuncType) []parser2.Param {
	params := make([]parser2.Param, 0)
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
		params = append(params, parser2.Param{
			Name: name,
			Type: field.Type.(*ast.Ident).Name,
		})
	}
	return params
}