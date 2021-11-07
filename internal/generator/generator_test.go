package generator

import (
	"bytes"
	"testing"

	"github.com/nickwallen/mocksie/internal/parser"
	"github.com/stretchr/testify/require"
)

func Test_Generator_GenerateMock_OK(t *testing.T) {
	tests := []struct {
		name     string
		iface    *parser.Interface
		expected string
	}{
		{
			name: "multiple-parameters",
			iface: &parser.Interface{
				Name: "greeter",
				Methods: []parser.Method{
					{
						Name: "SayHello",
						Params: []parser.Param{
							{Name: "first", Type: "string"},
							{Name: "last", Type: "string"},
						},
						Results: []parser.Result{
							{Name: "", Type: "error"},
						},
					},
				},
			},
			expected: `
// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (first string, last string) error
}

// SayHello relies on DoSayHello for defining it's behavior.
func (m *mockGreeter) SayHello(first string, last string) error {
    return m.DoSayHello(first, last)
}
`,
		},
		{
			name: "multiple-results",
			iface: &parser.Interface{
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
				},
			},
			expected: `
// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (name string) (string, error)
}

// SayHello relies on DoSayHello for defining it's behavior.
func (m *mockGreeter) SayHello(name string) (string, error) {
    return m.DoSayHello(name)
}
`,
		},
		{
			name: "named-results",
			iface: &parser.Interface{
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
			expected: `
// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (name string) (greeting string, err error)
}

// SayHello relies on DoSayHello for defining it's behavior.
func (m *mockGreeter) SayHello(name string) (greeting string, err error) {
    return m.DoSayHello(name)
}
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := bytes.NewBufferString("")
			gen := NewGenerator()
			gen.writer = writer

			err := gen.GenerateMock(test.iface)
			require.NoError(t, err)
			require.Equal(t, test.expected, writer.String())
		})
	}
}
