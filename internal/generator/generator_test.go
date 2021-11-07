package generator

import (
	"bytes"
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
			name: "multiple-parameters",
			iface: &mocksie.Interface{
				Name: "greeter",
				Methods: []mocksie.Method{
					{
						Name: "SayHello",
						Params: []mocksie.Param{
							{Name: "first", Type: "string"},
							{Name: "last", Type: "string"},
						},
						Results: []mocksie.Result{
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

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(first string, last string) error {
    return m.DoSayHello(first, last)
}

`,
		},
		{
			name: "multiple-results",
			iface: &mocksie.Interface{
				Name: "greeter",
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
			name: "named-results",
			iface: &mocksie.Interface{
				Name: "greeter",
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
			name: "multiple-methods",
			iface: &mocksie.Interface{
				Name: "greeter",
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
		})
	}
}
