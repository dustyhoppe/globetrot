package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var baselineCommand = &cobra.Command{
	Use:  "baseline",
	Long: "Generates a baseline for the database and updates the internal tables for change tracking",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Baseline ran")
	},
}

func init() {
	rootCommand.AddCommand(baselineCommand)
}
