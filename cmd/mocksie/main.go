package main

import (
	"github.com/spf13/cobra"
)

func main() {
	cmd := NewRootCommand()
	cmd.AddCommand(NewGenerateCmd())
	err := cmd.Execute()
	cobra.CheckErr(err)
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mocksie",
		Short: "Mocksie is a no-framework, mocking framework for the Go programming language.",
	}

	return cmd
}
