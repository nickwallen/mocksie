package parser

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/nickwallen/mocksie/internal"
	"github.com/stretchr/testify/require"
)

func Test_FileParser_FindInterfaces_OK(t *testing.T) {
	tests := []struct {
		testCase string
		code     []byte
		name     string
		expected *mocksie.Interface
		err      error
	}{
		{
			testCase: "methods-many",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(name string) (string, error)
					SayGoodbye(name string) (string, error)
				}
			`),
			name: "greeter",
			expected: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Imports: []mocksie.Import{},
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "name", Type: "string"},
						},
						Results: []mocksie.Result{
							{Name: "", Type: "string"},
							{Name: "", Type: "error"},
						},
					},
					{
						Name: "SayGoodbye",
						Params: []mocksie.Param{
							{Name: "name", Type: "string"},
						},
						Results: []mocksie.Result{
							{Name: "", Type: "string"},
							{Name: "", Type: "error"},
						},
					},
				},
			},
		},
		{
			testCase: "methods-none",
			code: []byte(`
				package main
				type greeter interface {
				}
			`),
			name: "greeter",
			expected: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Imports: []mocksie.Import{},
				Methods: []mocksie.Method{},
			},
		},
		{
			testCase: "interface-not-defined",
			code: []byte(`
				package main
				// No interfaces defined here
			`),
			name: "greeter",
			err:  errNotFound,
		},
		{
			testCase: "interface-not-found",
			code: []byte(`
				package main
				type greeter interface {
				}
			`),
			name: "doesNotExist",
			err:  errNotFound,
		},
		{
			testCase: "results-named",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(name string) (greeting string, err error)
				}
			`),
			name: "greeter",
			expected: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Imports: []mocksie.Import{},
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "name", Type: "string"},
						},
						Results: []mocksie.Result{
							{Name: "greeting", Type: "string"},
							{Name: "err", Type: "error"},
						},
					},
				},
			},
		},
		{
			testCase: "results-none",
			code: []byte(`
				package main
				import (
					"io"
				)
				type greeter interface {
					SayHello(name string, out io.Writer)
				}
			`),
			name: "greeter",
			expected: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Imports: []mocksie.Import{
					{Path: "io"},
				},
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "name", Type: "string"},
							{Name: "out", Type: "io.Writer"},
						},
						Results: []mocksie.Result{},
					},
				},
			},
		},
		{
			testCase: "params-unnamed",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(string) (string, error)
				}
			`),
			name: "greeter",
			expected: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Imports: []mocksie.Import{},
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "", Type: "string"},
						},
						Results: []mocksie.Result{
							{Name: "", Type: "string"},
							{Name: "", Type: "error"},
						},
					},
				},
			},
		},
		{
			testCase: "imports",
			code: []byte(`
				package main
				import "io"
				type greeter interface {
					SayHello(in io.Reader, out io.Writer) error
				}
			`),
			name: "greeter",
			expected: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Imports: []mocksie.Import{
					{Path: "io"},
				},
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "in", Type: "io.Reader"},
							{Name: "out", Type: "io.Writer"},
						},
						Results: []mocksie.Result{
							{Name: "", Type: "error"},
						},
					},
				},
			},
		},
	}

	// Create a file for the source code
	file, err := ioutil.TempFile("", "interfaces.go")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			// Write the source code to the file
			err = ioutil.WriteFile(file.Name(), test.code, 0700)
			require.NoError(t, err)

			// Ensure the test code is valid
			_, err = parser.ParseFile(token.NewFileSet(), "", test.code, parser.AllErrors)
			require.NoError(t, err)

			// Create the file parser
			p, err := New(file.Name())
			require.NoError(t, err)

			// Find all interfaces
			found, err := p.FindInterface(test.name)
			if test.expected != nil {
				require.Equal(t, test.expected, found)
			}
			require.Equal(t, test.err, err)
		})
	}
}

func Test_FileParser_NewFileParser_OK(t *testing.T) {
	// Create a file for the source code
	file, err := ioutil.TempFile("", "interfaces.go")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	// Write the source code to the file
	code := []byte(`
		package main
		type greeter interface {
			SayHello(name string) (string, error)
			SayGoodbye(name string) (string, error)
		}
	`)
	err = ioutil.WriteFile(file.Name(), code, 0700)
	require.NoError(t, err)

	// Create the parser
	_, err = New(file.Name())
	require.NoError(t, err)
}

func Test_FileParser_NewFileParser_FileDoesNotExist(t *testing.T) {
	_, err := New("/this/file/does/not/exist.go")
	require.Error(t, err)
}
