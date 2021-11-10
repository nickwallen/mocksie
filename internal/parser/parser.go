package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

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
			// Found an interface
			typ := spec.(*ast.TypeSpec)
			ifaces = append(ifaces, &mocksie.Interface{
				Name:    typ.Name.String(),
				Package: buildPackage(f),
				Imports: buildImports(f),
				Methods: buildMethods(typ.Type.(*ast.InterfaceType)),
			})
		}
	}
	return ifaces, nil
}

func buildPackage(f *ast.File) mocksie.Package {
	return mocksie.Package(f.Name.Name)
}

func buildImports(f *ast.File) []mocksie.Import {
	imports := make([]mocksie.Import, 0)
	for _, impSpec := range f.Imports {
		imports = append(imports, mocksie.Import{
			Path: strings.ReplaceAll(impSpec.Path.Value, "\"", ""),
		})
	}
	return imports
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
	if funcType.Results == nil {
		return results // No function results
	}
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
	for _, field := range funcType.Params.List {
		// The param may not be named
		name := ""
		if len(field.Names) > 0 {
			name = field.Names[0].Name
		}

		switch typ := field.Type.(type) {
		case *ast.Ident:
			// Build the param
			params = append(params, mocksie.Param{
				Name: name,
				Type: typ.Name,
			})

		case *ast.SelectorExpr:
			i, ok := typ.X.(*ast.Ident)
			if !ok {
				// TODO fix me
				log.Fatalf("expected *ast.Ident, but got something else.")
			}
			// Build the param
			params = append(params, mocksie.Param{
				Name: name,
				Type: i.Name + "." + typ.Sel.Name,
			})

		default:
			// Not supported
			continue
		}
	}
	return params
}
