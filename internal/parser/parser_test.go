package parser

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/nickwallen/mocksie/internal"
	"github.com/stretchr/testify/require"
)

func Test_FileParser_FindInterfaces_OK(t *testing.T) {
	tests := []struct {
		name     string
		code     []byte
		expected []*mocksie.Interface
	}{
		{
			name: "methods-many",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(name string) (string, error)
					SayGoodbye(name string) (string, error)
				}
			`),
			expected: []*mocksie.Interface{
				{
					Name:    "greeter",
					Package: "main",
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
		},
		{
			name: "methods-none",
			code: []byte(`
				package main
				type greeter interface {
				}
			`),
			expected: []*mocksie.Interface{
				{
					Name:    "greeter",
					Package: "main",
					Methods: []mocksie.Method{},
				},
			},
		},
		{
			name: "interfaces-many",
			code: []byte(`
				package main
				type thisOne interface {
					DoThisThing() (string, error)
				}
				type thatOne interface {
					DoThatThing() (string, error)
				}
			`),
			expected: []*mocksie.Interface{
				{
					Name:    "thisOne",
					Package: "main",
					Methods: []mocksie.Method{
						{
							Name:   "DoThisThing",
							Params: []mocksie.Param{},
							Results: []mocksie.Result{
								{Name: "", Type: "string"},
								{Name: "", Type: "error"},
							},
						},
					},
				},
				{
					Name:    "thatOne",
					Package: "main",
					Methods: []mocksie.Method{
						{
							Name:   "DoThatThing",
							Params: []mocksie.Param{},
							Results: []mocksie.Result{
								{Name: "", Type: "string"},
								{Name: "", Type: "error"},
							},
						},
					},
				},
			},
		},
		{
			name: "interfaces-none",
			code: []byte(`
				package main
				// No interfaces defined here
			`),
			expected: []*mocksie.Interface{},
		},
		{
			name: "results-named",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(name string) (greeting string, err error)
				}
			`),
			expected: []*mocksie.Interface{
				{
					Name:    "greeter",
					Package: "main",
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
		},
		{
			name: "params-unnamed",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(string) (string, error)
				}
			`),
			expected: []*mocksie.Interface{
				{
					Name:    "greeter",
					Package: "main",
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
		},
	}

	// Create a file for the source code
	file, err := ioutil.TempFile("", "interfaces.go")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Write the source code to the file
			err = ioutil.WriteFile(file.Name(), test.code, 0700)
			require.NoError(t, err)

			// Create the file parser
			parser, err := New(file.Name())
			require.NoError(t, err)

			// Find all interfaces
			found, err := parser.FindInterfaces()
			require.NoError(t, err)
			require.Equal(t, test.expected, found)
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
