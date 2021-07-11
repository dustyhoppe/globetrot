package cmd

import (
	"github.com/spf13/cobra"

	"globetrot/utils"
)

var (
	configFile   string
	username     string
	password     string
	host         string
	port         int
	database     string
	filePath     string
	version      string
	environment  string
	dryRun       bool
	databaseType string
)

var migrateCommand = &cobra.Command{
	Use:  "migrate",
	Long: "Performs a database migration",
	Run: func(cmd *cobra.Command, args []string) {

		config := &utils.Config{
			Username:    username,
			Password:    password,
			Host:        host,
			Database:    database,
			Port:        port,
			Type:        databaseType,
			FilePath:    filePath,
			DryRun:      dryRun,
			Environment: environment,
		}

		runner := new(utils.Runner)
		runner.Init(*config)

		runner.Migrate()
	},
}

func init() {
	migrateCommand.PersistentFlags().StringVarP(&configFile, "configPath", "c", "", "The path to the configuration to use in-place of passing command line arguments")
	migrateCommand.PersistentFlags().StringVarP(&username, "username", "u", "root", "The username to use when connecting to the database")
	migrateCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "The password to use when connecting to the database")
	migrateCommand.PersistentFlags().StringVarP(&host, "server", "s", "localhost", "The host server the database is located at")
	migrateCommand.PersistentFlags().IntVarP(&port, "port", "P", 3306, "The port the database server the database is listening at")
	migrateCommand.PersistentFlags().StringVarP(&database, "database", "d", "", "The name of the database")
	migrateCommand.PersistentFlags().StringVarP(&filePath, "filePath", "f", ".\\", "The directory where your SQL scripts are located")
	migrateCommand.PersistentFlags().BoolVar(&dryRun, "dryRun", false, "Indicates whether the command should perform a dry run")
	migrateCommand.PersistentFlags().StringVarP(&databaseType, "databaseType", "t", "mysql", "Indicates the database type/platform (mysql, postgres, sqlserver)")
	migrateCommand.PersistentFlags().StringVarP(&environment, "env", "e", "", "The environment the migration is targeting")

	migrateCommand.MarkPersistentFlagRequired("username")
	migrateCommand.MarkPersistentFlagRequired("password")
	migrateCommand.MarkPersistentFlagRequired("server")
	migrateCommand.MarkPersistentFlagRequired("port")
	migrateCommand.MarkPersistentFlagRequired("database")
	migrateCommand.MarkPersistentFlagRequired("filePath")
	migrateCommand.MarkPersistentFlagRequired("databaseType")

	rootCommand.AddCommand(migrateCommand)
}
