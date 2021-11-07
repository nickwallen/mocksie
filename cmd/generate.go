package cmd

import (
	"log"

	"github.com/nickwallen/mocksie/internal/generator"
	"github.com/nickwallen/mocksie/internal/parser"
	"github.com/spf13/cobra"
)

var generateArgs = struct {
	filename string
}{}

// generateCmd the command that generates mock implementations of an interface.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate mock implementation.",
	Run: func(cmd *cobra.Command, _ []string) {
		out := cmd.OutOrStdout()
		log.SetOutput(out)

		// Which file to parse?
		if len(generateArgs.filename) == 0 {
			cobra.CheckErr("undefined filename; use --file")
		}

		// Find all interfaces
		parser, err := parser.NewFileParser(generateArgs.filename)
		cobra.CheckErr(err)

		ifaces, err := parser.FindInterfaces()
		cobra.CheckErr(err)
		if len(ifaces) == 0 {
			return // No interfaces found. Nothing to do.
		}

		// Generate the generate
		gen, err := generator.NewGenerator(out)
		cobra.CheckErr(err)

		// Create a mock for each interface
		for _, iface := range ifaces {
			err := gen.GenerateMock(iface)
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Define settings for the generate command
	generateCmd.PersistentFlags().StringVarP(&generateArgs.filename, "file", "f", "", "Generates mocks for all interfaces defined within a file.")
}