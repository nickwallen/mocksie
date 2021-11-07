package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/nickwallen/mocksie/internal/generator"
	"github.com/nickwallen/mocksie/internal/parser"
)

var generateArgs = struct {
	filename string
}{}

// generateCmd the command that generates mock implementations of an interface.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate mocks.",
	Run: func(cmd *cobra.Command, _ []string) {
		out := cmd.OutOrStdout()
		log.SetOutput(out)

		// Which file to parse?
		if len(generateArgs.filename) == 0 {
			cobra.CheckErr("undefined filename; use --file")
		}

		// Find all interfaces
		parser, err := parser.New(generateArgs.filename)
		cobra.CheckErr(err)
		ifaces, err := parser.FindInterfaces()
		cobra.CheckErr(err)
		if len(ifaces) == 0 {
			return // No interfaces found. Nothing to do.
		}

		// Generate a mock for each interface found
		gen, err := generator.New(out)
		cobra.CheckErr(err)
		for _, iface := range ifaces {
			err := gen.GenerateMock(iface)
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Define settings for the generate command
	generateCmd.PersistentFlags().StringVarP(&generateArgs.filename, "file", "f", "", "Generate mocks for all interfaces defined within a file.")
}
