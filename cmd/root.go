package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	username       string
	password       string
	host           string
	port           int
	database       string
	filePath       string
	commandTimeout int
	version        string
	environment    string
	dryRun         bool
	databaseType   string
)

var rootCommand = &cobra.Command{
	Use:   "globetrot [subcommand]",
	Short: "A database change management tool",
	Long:  "Globetrot is a CLI tool used for managing database changes/migrations",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Hello CLI") },
}

func init() {

	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "root", "The username to use when connecting to the database")
	rootCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "The password to use when connecting to the database")
	rootCommand.PersistentFlags().StringVarP(&host, "server", "s", "localhost", "The host server the database is located at")
	rootCommand.PersistentFlags().IntVarP(&port, "port", "P", 3306, "The port the database server the database is listening at")
	rootCommand.PersistentFlags().StringVarP(&database, "database", "d", "", "The name of the database")
	rootCommand.PersistentFlags().StringVarP(&filePath, "filePath", "f", ".\\", "The directory where your SQL scripts are located")
	rootCommand.PersistentFlags().BoolVar(&dryRun, "dryRun", false, "Indicates whether the command should perform a dry run")
	rootCommand.PersistentFlags().StringVarP(&databaseType, "databaseType", "t", "mysql", "Indicates the database type/platform")
	rootCommand.PersistentFlags().StringVarP(&environment, "env", "e", "", "The environment the migration is targeting")

	rootCommand.MarkPersistentFlagRequired("username")
	rootCommand.MarkPersistentFlagRequired("password")
	rootCommand.MarkPersistentFlagRequired("host")
	rootCommand.MarkPersistentFlagRequired("port")
	rootCommand.MarkPersistentFlagRequired("database")
	rootCommand.MarkPersistentFlagRequired("filePath")
	rootCommand.MarkPersistentFlagRequired("databaseType")
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
