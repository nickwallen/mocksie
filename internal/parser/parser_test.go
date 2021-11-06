package parser

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nickwallen/mocksie/internal"
)

func Test_FileParser_FindInterfaces_OK(t *testing.T) {
	tests := []struct {
		name     string
		code     []byte
		expected []*parser.Interface
	}{
		{
			name: "multiple-methods",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(name string) (string, error)
					SayGoodbye(name string) (string, error)
				}
			`),
			expected: []*parser.Interface{
				{
					Name: "greeter",
					Methods: []parser.Method{
						{
							Name: "SayHello",
							Params: []parser.Param{
								{Name: "name", Type: "string"},
							},
							Results: []parser.Result{
								{Name: "", Type: "string"},
								{Name: "", Type: "error"},
							},
						},
						{
							Name: "SayGoodbye",
							Params: []parser.Param{
								{Name: "name", Type: "string"},
							},
							Results: []parser.Result{
								{Name: "", Type: "string"},
								{Name: "", Type: "error"},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple-interfaces",
			code: []byte(`
				package main
				type thisOne interface {
					DoThisThing() (string, error)
				}
				type thatOne interface {
					DoThatThing() (string, error)
				}
			`),
			expected: []*parser.Interface{
				{
					Name: "thisOne",
					Methods: []parser.Method{
						{
							Name:   "DoThisThing",
							Params: []parser.Param{},
							Results: []parser.Result{
								{Name: "", Type: "string"},
								{Name: "", Type: "error"},
							},
						},
					},
				},
				{
					Name: "thatOne",
					Methods: []parser.Method{
						{
							Name:   "DoThatThing",
							Params: []parser.Param{},
							Results: []parser.Result{
								{Name: "", Type: "string"},
								{Name: "", Type: "error"},
							},
						},
					},
				},
			},
		},
		{
			name: "no-interfaces",
			code: []byte(`
				package main
				// No interfaces defined here
			`),
			expected: []*parser.Interface{},
		},
		{
			name: "named-results",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(name string) (greeting string, err error)
				}
			`),
			expected: []*parser.Interface{
				{
					Name: "greeter",
					Methods: []parser.Method{
						{
							Name: "SayHello",
							Params: []parser.Param{
								{Name: "name", Type: "string"},
							},
							Results: []parser.Result{
								{Name: "greeting", Type: "string"},
								{Name: "err", Type: "error"},
							},
						},
					},
				},
			},
		},
		{
			name: "unnamed-params",
			code: []byte(`
				package main
				type greeter interface {
					SayHello(string) (string, error)
				}
			`),
			expected: []*parser.Interface{
				{
					Name: "greeter",
					Methods: []parser.Method{
						{
							Name: "SayHello",
							Params: []parser.Param{
								{Name: "", Type: "string"},
							},
							Results: []parser.Result{
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
			parser, err := NewFileParser(file.Name())
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
	_, err = NewFileParser(file.Name())
	require.NoError(t, err)
}

func Test_FileParser_NewFileParser_FileDoesNotExist(t *testing.T) {
	_, err := NewFileParser("/this/file/does/not/exist.go")
	require.Error(t, err)
}
