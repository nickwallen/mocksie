package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nickwallen/mocksie/internal/generator"
	"github.com/nickwallen/mocksie/internal/parser"
	"github.com/spf13/cobra"
)

var generateArgs = struct {
	inFile  string
	outFile string
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

			// Which file to parse?
			if len(generateArgs.inFile) == 0 {
				return fmt.Errorf("no input defined")
			}

			// Find all interfaces
			p, err := parser.New(generateArgs.inFile)
			if err != nil {
				return err
			}
			ifaces, err := p.FindInterfaces()
			if err != nil {
				return err
			}
			if len(ifaces) == 0 {
				return nil // No interfaces found. Nothing to do.
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

			// Generate a mock for each interface found
			gen, err := generator.New(out)
			if err != nil {
				return err
			}

			for _, iface := range ifaces {
				err := gen.GenerateMock(iface)
				if err != nil {
					return err
				}
			}
			return nil // Success
		},
	}

	// Define the accepted flags
	cmd.PersistentFlags().StringVarP(&generateArgs.inFile, "in", "i", "",
		"The input file containing the interface definition.")
	cmd.PersistentFlags().StringVarP(&generateArgs.outFile, "out", "o", "",
		"The output file to write the generated mocks to.")
	return cmd
}
