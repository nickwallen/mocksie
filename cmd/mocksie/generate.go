package mocksie

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/nickwallen/mocksie/internal/generator"
	"github.com/nickwallen/mocksie/internal/parser"
)

var generateArgs = struct {
	filename string
}{}

func init() {
	rootCmd.AddCommand(NewGenerateCmd())
}

// NewGenerateCmd a command that generates mock implementations of an interface.
func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate mocks.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()
			log.SetOutput(out)

			// Which file to parse?
			if len(generateArgs.filename) == 0 {
				return fmt.Errorf("no input defined")
			}

			// Find all interfaces
			parser, err := parser.New(generateArgs.filename)
			if err != nil {
				return err
			}
			ifaces, err := parser.FindInterfaces()
			if err != nil {
				return err
			}
			if len(ifaces) == 0 {
				return nil // No interfaces found. Nothing to do.
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
	cmd.PersistentFlags().StringVarP(&generateArgs.filename, "file", "f", "",
		"Generate mocks for all interfaces defined within a file.")

	return cmd
}
