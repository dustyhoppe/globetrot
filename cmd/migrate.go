package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateCommand = &cobra.Command{
	Use:  "migrate",
	Long: "Performs a database migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Migration ran")
	},
}

func init() {
	rootCommand.AddCommand(migrateCommand)
}
