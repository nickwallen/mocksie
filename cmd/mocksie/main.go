package main

import (
	"github.com/spf13/cobra"
)

func main() {
	cmd := NewGenerateCmd()
	err := cmd.Execute()
	cobra.CheckErr(err)
}
