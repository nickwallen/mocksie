package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/nickwallen/mocksie/internal"
)

var (
	errNotFound = errors.New("interface not found")
)

// Parser parses a file containing Go source code.
type Parser struct {
	filename string
}

// New constructs a new Parser.
func New(inFile string) (*Parser, error) {
	// Find the absolute path to the input file
	inFile, err := filepath.Abs(inFile)
	if err != nil {
		return nil, err
	}

	// Ensure the file exists
	if _, err := os.Stat(inFile); err != nil {
		return nil, err
	}

	return &Parser{filename: inFile}, nil
}

// FindInterface returns the interface with the given name.
func (p *Parser) FindInterface(name string) (*mocksie.Interface, error) {
	// Parse the file
	f, err := parser.ParseFile(token.NewFileSet(), p.filename, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	// Find any interfaces
	for _, decl := range f.Decls {
		// Expect a declaration
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// Expect that a type is being declared
		if genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			// Expect a proper type
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// Expect an interface
			ifaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			// Is this the interface that we are looking for?
			if name == typeSpec.Name.String() {
				return buildInterface(typeSpec.Name.String(), ifaceType, f)
			}
		}
	}
	return nil, errNotFound
}

// buildInterface Returns an interface.
func buildInterface(name string, typ *ast.InterfaceType, f *ast.File) (*mocksie.Interface, error) {
	methods, err := buildMethods(typ)
	if err != nil {
		return nil, err
	}
	return &mocksie.Interface{
		Name:    name,
		Package: buildPackage(f),
		Imports: buildImports(f),
		Methods: methods,
	}, nil
}

// buildPackage Returns the package defined within a file.
func buildPackage(f *ast.File) mocksie.Package {
	return mocksie.Package(f.Name.Name)
}

// buildImports Returns the imports defined within a file.
func buildImports(f *ast.File) []mocksie.Import {
	imports := make([]mocksie.Import, 0)
	for _, impSpec := range f.Imports {
		imports = append(imports, mocksie.Import{
			Path: strings.ReplaceAll(impSpec.Path.Value, "\"", ""),
		})
	}
	return imports
}

// buildMethods Returns the methods of an interface.
func buildMethods(typ *ast.InterfaceType) ([]mocksie.Method, error) {
	methods := make([]mocksie.Method, 0)
	for _, field := range typ.Methods.List {
		// Expect the method to be named
		if len(field.Names) == 0 {
			continue
		}

		// Expect a function type
		funcType, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		params, err := buildParams(funcType)
		if err != nil {
			return nil, err
		}

		// Build the method
		methods = append(methods, mocksie.Method{
			Name:    field.Names[0].Name,
			Params:  params,
			Results: buildResults(funcType),
		})
	}
	return methods, nil
}

// buildResults Returns the results (return values) of an interface method.
func buildResults(funcType *ast.FuncType) []mocksie.Result {
	results := make([]mocksie.Result, 0)
	if funcType.Results == nil {
		return results // No function results
	}
	for i := range funcType.Results.List {
		field := funcType.Results.List[i]

		// Expect the field to be an identifier
		typ, ok := field.Type.(*ast.Ident)
		if !ok {
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
			Type: typ.Name,
		})
	}
	return results
}

// buildParams Returns the parameters of an interface method.
func buildParams(funcType *ast.FuncType) ([]mocksie.Param, error) {
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
			ident, ok := typ.X.(*ast.Ident)
			if !ok {
				return nil, fmt.Errorf("expected *ast.Ident, but got something else")
			}

			// Build the param
			params = append(params, mocksie.Param{
				Name: name,
				Type: ident.Name + "." + typ.Sel.Name,
			})

		default:
			// Not supported
			continue
		}
	}
	return params, nil
}
