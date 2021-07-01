package cmd

import (
	"github.com/spf13/cobra"

	"globetrot/utils"
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
	rootCommand.AddCommand(migrateCommand)
}
