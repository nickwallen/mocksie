package generator

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nickwallen/mocksie/internal"
)

func Test_Generator_GenerateMock_OK(t *testing.T) {
	expected := `
type mockGreeter struct {
}
`
	writer := bytes.NewBufferString("")
	gen := NewGenerator()
	gen.writer = writer

	err := gen.GenerateMock(&parser.Interface{
		Name:    "greeter",
		Methods: nil,
	})
	require.NoError(t, err)
	require.Equal(t, expected, writer.String())
}

