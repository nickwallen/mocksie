package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateCmd_OK(t *testing.T) {
	var out bytes.Buffer

	// Read the expected mock output
	expectedMock, err := ioutil.ReadFile("../../internal/testdata/mockGreeter.go")
	require.NoError(t, err)

	// Run the command
	cmd := NewGenerateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{
		"--name", "greeter",
		"--in", "../../internal/testdata/greeter.go",
	})
	err = cmd.Execute()
	require.NoError(t, err)
	require.Equal(t, string(expectedMock), out.String())
}

func Test_GenerateCmd_InterfaceNotFound(t *testing.T) {
	var out bytes.Buffer

	// Run the command
	cmd := NewGenerateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{
		"--name", "doesNotExist", // An interface by this name does not exist
		"--in", "../../internal/testdata/greeter.go",
	})
	err := cmd.Execute()
	require.Error(t, err)
}

func Test_GenerateCmd_ChooseOneInterface(t *testing.T) {
	var out bytes.Buffer

	// Read the expected mock output
	expectedMock, err := ioutil.ReadFile("../../internal/testdata/mockHelloGreeter.go")
	require.NoError(t, err)

	// The input file has multiple interfaces defined
	cmd := NewGenerateCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{
		"--name", "helloGreeter",
		"--in", "../../internal/testdata/greeters.go",
	})

	// If there is more than one interface found, the user must choose
	// which interface to generate a mock for using the --name flag.
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
		"--name", "greeter",
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
