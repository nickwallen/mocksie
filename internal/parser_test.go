package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_FileParser_FindInterfaces_OK(t *testing.T) {
	tests := []struct {
		filename string
		expected []*Interface
	}{
		{
			filename: "testdata/greeter.go",
			expected: []*Interface{
				{
					Name:    "greeter",
					Methods: []string{"SayHello", "SayGoodbye"},
				},
			},
		},
		{
			filename: "testdata/multiple.go",
			expected: []*Interface{
				{
					Name: "thisOne",
					Methods: []string{"DoThisThing"},
				},
				{
					Name:"thatOne",
					Methods: []string{"DoThatThing"},
				},
				{
					Name:"anotherOne",
					Methods: []string{"DoAnotherThing"},
				},
			},
		},
		{
			filename: "testdata/none.go",
			expected: []*Interface{},
		},
	}
	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {

			// Create the file parser
			parser, err := NewFileParser(test.filename)
			require.NoError(t, err)

			// Find all interfaces
			found, err := parser.FindInterfaces()
			require.Equal(t, test.expected, found)
		})
	}
}

func Test_FileParser_NewFileParser_OK(t *testing.T) {
	_, err := NewFileParser("testdata/greeter.go")
	require.NoError(t, err)
}

func Test_FileParser_NewFileParser_FileDoesNotExist(t *testing.T) {
	_, err := NewFileParser("/this/file/does/not/exist.go")
	require.Error(t, err)
}
