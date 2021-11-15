package main

import (
	"log"
	"os"

	"github.com/nickwallen/mocksie/internal/generator"
	"github.com/nickwallen/mocksie/internal/parser"
	"github.com/spf13/cobra"
)

var generateArgs = struct {
	inFile  string // Input file containing the interface definition.
	outFile string // Output file to write the generated mock to.
	name    string // Name of the interface to generate a mock for.
}{}

// NewGenerateCmd a command that generates mock implementations of an interface.
func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mocksie",
		Short: "Mocksie will generate mocks for your Golang interfaces ",
		Long: `
Mocksie is a no-framework, mocking framework for the Go programming 
language.

Mocksie will generate mocks for your Golang interfaces that don't require any 
additional mocking or test packages. This allows you to define mock behavior
directly in the test case, does not require learning a new mocking framework,
prevents you from having to maintain boilerplate code, and 
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()
			log.SetOutput(out)

			// Find the interface definition
			p, err := parser.New(generateArgs.inFile)
			if err != nil {
				return err
			}
			found, err := p.FindInterface(generateArgs.name)
			if err != nil {
				return err
			}

			// Open the output file or use stdout if not output file defined
			if len(generateArgs.outFile) > 0 {
				outFile, err := os.Create(generateArgs.outFile)
				if err != nil {
					return err
				}
				defer outFile.Close() // TODO handle the error
				out = outFile
			}

			// Generate the mock
			gen, err := generator.New(out)
			if err != nil {
				return err
			}
			return gen.GenerateMock(found)
		},
	}

	// Define the accepted flags
	cmd.Flags().StringVarP(&generateArgs.inFile, "in", "i", "", "The input file containing the interface definition.")
	cmd.Flags().StringVarP(&generateArgs.outFile, "out", "o", "", "The output file to write the generated mocks to.")
	cmd.Flags().StringVarP(&generateArgs.name, "name", "n", "", "The name of the interface to generate a mock for.")
	err := cmd.MarkFlagRequired("name")
	cobra.CheckErr(err)

	return cmd
}
