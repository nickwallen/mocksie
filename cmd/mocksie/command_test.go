package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateCmd_OneInterface(t *testing.T) {
	var out bytes.Buffer

	// Read the expected mock output
	expectedMock, err := ioutil.ReadFile("../../internal/testdata/mockGreeter.go")
	require.NoError(t, err)

	// Run the command
	cmd := NewGenerateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--in", "../../internal/testdata/greeter.go"})
	err = cmd.Execute()
	require.NoError(t, err)
	require.Equal(t, string(expectedMock), out.String())
}

func Test_GenerateCmd_TwoInterfaces(t *testing.T) {
	t.SkipNow() // Skip until we deal with multiple interfaces
	var out bytes.Buffer

	// Read the expected mock output
	expectedMock, err := ioutil.ReadFile("../../internal/testdata/mockGreeter.go")
	require.NoError(t, err)

	cmd := NewGenerateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--in", "../../internal/testdata/greeters.go"})
	err = cmd.Execute()
	require.NoError(t, err)
	require.Equal(t, string(expectedMock), out.String())
}

func Test_GenerateCmd_OutFile(t *testing.T) {
	// Create a temp file for the output
	outFile, err := ioutil.TempFile("", "mockGreeter.go")
	require.NoError(t, err)
	defer outFile.Close()

	// Read the expected mock output
	expectedMock, err := ioutil.ReadFile("../../internal/testdata/mockGreeter.go")
	require.NoError(t, err)

	// Generate a mock
	cmd := NewGenerateCmd()
	cmd.SetArgs([]string{
		"--in", "../../internal/testdata/greeter.go",
		"--out", outFile.Name(),
	})
	err = cmd.Execute()
	require.NoError(t, err)

	// Validate the generated mock
	generatedMock, err := ioutil.ReadFile(outFile.Name())
	require.NoError(t, err)
	require.Equal(t, string(expectedMock), string(generatedMock))
}

func Test_GenerateCmd_NoArgs(t *testing.T) {
	cmd := NewGenerateCmd()
	err := cmd.Execute()
	require.Error(t, err, "Passing no flags is an error.")
}
