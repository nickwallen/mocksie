package generator

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"

	"github.com/nickwallen/mocksie/internal"
	"github.com/stretchr/testify/require"
)

func Test_Generator_GenerateMock_OK(t *testing.T) {
	tests := []struct {
		name     string
		iface    *mocksie.Interface
		expected string
	}{
		{
			name: "basic-greeter",
			iface: &mocksie.Interface{
				Name:    "greeter",
				Package: "testdata",
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
				},
			},
			expected: `
package testdata



// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (name string) (string, error)
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(name string) (string, error) {
    return m.DoSayHello(name)
}

`,
		},
		{
			name: "methods-multiple",
			iface: &mocksie.Interface{
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
			expected: `
package main



// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (name string) (string, error)
    DoSayGoodbye func (name string) (string, error)
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(name string) (string, error) {
    return m.DoSayHello(name)
}

// SayGoodbye relies on DoSayGoodbye for defining it's behavior. If this is causing a panic,
// define DoSayGoodbye within your test case.
func (m *mockGreeter) SayGoodbye(name string) (string, error) {
    return m.DoSayGoodbye(name)
}

`,
		},
		{
			name: "results-one",
			iface: &mocksie.Interface{
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
						},
					},
				},
			},
			expected: `
package main



// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (name string) string
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(name string) string {
    return m.DoSayHello(name)
}

`,
		},
		{
			name: "results-none",
			iface: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "name", Type: "string"},
						},
						Results: []mocksie.Result{},
					},
				},
			},
			expected: `
package main



// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (name string) 
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(name string)  {
    m.DoSayHello(name)
}

`,
		},
		{
			name: "results-named",
			iface: &mocksie.Interface{
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
			expected: `
package main



// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (name string) (greeting string, err error)
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(name string) (greeting string, err error) {
    return m.DoSayHello(name)
}

`,
		},
		{
			name: "params-multiple",
			iface: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "first", Type: "string"},
							{Name: "last", Type: "string"},
						},
						Results: []mocksie.Result{
							{Name: "", Type: "string"},
							{Name: "", Type: "error"},
						},
					},
				},
			},
			expected: `
package main



// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (first string, last string) (string, error)
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(first string, last string) (string, error) {
    return m.DoSayHello(first, last)
}

`,
		},
		{
			name: "params-none",
			iface: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Methods: []mocksie.Method{
					{
						Name:   "SayHello",
						Params: []mocksie.Param{},
						Results: []mocksie.Result{
							{Name: "", Type: "string"},
							{Name: "", Type: "error"},
						},
					},
				},
			},
			expected: `
package main



// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func () (string, error)
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello() (string, error) {
    return m.DoSayHello()
}

`,
		},
		{
			name: "params-unnamed",
			iface: &mocksie.Interface{
				Name:    "greeter",
				Package: "main",
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "name", Type: "string"},
						},
						Results: []mocksie.Result{},
					},
				},
			},
			expected: `
package main



// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (name string) 
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(name string)  {
    m.DoSayHello(name)
}

`,
		},
		{
			name: "imports",
			iface: &mocksie.Interface{
				Name:    "greeter",
				Package: "testdata",
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
			expected: `
package testdata

import (
    "io"
)

// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (in io.Reader, out io.Writer) error
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(in io.Reader, out io.Writer) error {
    return m.DoSayHello(in, out)
}

`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a generator
			out := bytes.NewBufferString("")
			gen, err := New(out)
			require.NoError(t, err)

			// Generate the mock
			err = gen.GenerateMock(test.iface)
			require.NoError(t, err)
			require.Equal(t, test.expected, out.String())

			// Validate the generated code
			_, err = parser.ParseFile(token.NewFileSet(), "", out.String(), parser.AllErrors)
			require.NoError(t, err)
		})
	}
}
