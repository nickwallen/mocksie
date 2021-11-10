package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateCmd_OK(t *testing.T) {
	out := bytes.NewBufferString("")
	cmd := NewGenerateCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--file", "../../internal/testdata/greeter.go"})
	err := cmd.Execute()
	require.NoError(t, err)
	expected := `
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

// mockGreeter ia a mock implementation of the Greeter interface.
type mockGreeter struct {
    DoSayHello func (in io.Reader, out io.Writer) error
    DoSayGoodbye func (in io.Reader, out io.Writer) error
}

// SayHello relies on DoSayHello for defining it's behavior. If this is causing a panic,
// define DoSayHello within your test case.
func (m *mockGreeter) SayHello(in io.Reader, out io.Writer) error {
    return m.DoSayHello(in, out)
}

// SayGoodbye relies on DoSayGoodbye for defining it's behavior. If this is causing a panic,
// define DoSayGoodbye within your test case.
func (m *mockGreeter) SayGoodbye(in io.Reader, out io.Writer) error {
    return m.DoSayGoodbye(in, out)
}

`
	require.Equal(t, expected, out.String())
}

func Test_GenerateCmd_NoArgs(t *testing.T) {
	cmd := NewGenerateCmd()
	err := cmd.Execute()
	require.Error(t, err, "Passing no flags is an error.")
}
