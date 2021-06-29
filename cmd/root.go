package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	connectionString string
	filePath         string
	commandTimeout   int
	version          string
	environment      string
	dryRun           bool
	databaseType     string
)

var rootCommand = &cobra.Command{
	Use:   "globetrot [subcommand]",
	Short: "A database change management tool",
	Long:  "Globetrot is a CLI tool used for managing database changes/migrations",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Hello CLI") },
}

func init() {

	rootCommand.PersistentFlags().StringVarP(&connectionString, "connectionString", "c", "", "Connection string to use to connect to database")
	rootCommand.PersistentFlags().StringVarP(&filePath, "filePath", "f", ".\\", "The directory where your SQL scripts are located")
	rootCommand.PersistentFlags().BoolVar(&dryRun, "dryRun", false, "Indicates whether the command should perform a dry run")
	rootCommand.PersistentFlags().StringVarP(&databaseType, "databaseType", "t", "mysql", "Indicates the database type/platform")

	rootCommand.MarkPersistentFlagRequired("connectionString")
	rootCommand.MarkPersistentFlagRequired("filePath")
	rootCommand.MarkPersistentFlagRequired("databaseType")
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
