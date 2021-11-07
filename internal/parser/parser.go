package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/nickwallen/mocksie/internal"
)

// Parser parses a file containing Go source code.
type Parser struct {
	filename string
}

// New constructs a new Parser.
func New(filename string) (*Parser, error) {
	// Find the absolute path to the file
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	// Ensure the file exists
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}

	return &Parser{filename: filename}, nil
}

// FindInterfaces returns all interfaces defined in the file.
func (p *Parser) FindInterfaces() ([]*mocksie.Interface, error) {
	// Parse the file
	f, err := parser.ParseFile(token.NewFileSet(), p.filename, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	// Find any interfaces
	ifaces := make([]*mocksie.Interface, 0)
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
			ifaces = append(ifaces, &mocksie.Interface{
				Name:    typ.Name.String(),
				Package: f.Name.Name,
				Methods: buildMethods(typ.Type.(*ast.InterfaceType)),
			})
		}
	}
	return ifaces, nil
}

func buildMethods(typ *ast.InterfaceType) []mocksie.Method {
	methods := make([]mocksie.Method, 0)
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
		methods = append(methods, mocksie.Method{
			Name:    field.Names[0].Name,
			Params:  buildParams(funcType),
			Results: buildResults(funcType),
		})
	}
	return methods
}

func buildResults(funcType *ast.FuncType) []mocksie.Result {
	results := make([]mocksie.Result, 0)
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
		results = append(results, mocksie.Result{
			Name: name,
			Type: field.Type.(*ast.Ident).Name,
		})
	}
	return results
}

func buildParams(funcType *ast.FuncType) []mocksie.Param {
	params := make([]mocksie.Param, 0)
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
		params = append(params, mocksie.Param{
			Name: name,
			Type: field.Type.(*ast.Ident).Name,
		})
	}
	return params
}
