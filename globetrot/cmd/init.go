package cmd

import (
	"github.com/spf13/cobra"

	"globetrot/utils"
)

var (
	initPath string
)

var initCommand = &cobra.Command{
	Use:  "init",
	Long: "Initializes a directory as a globetrot script directory",
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

		runner.Initialize(initPath)
	},
}

func init() {
	initCommand.PersistentFlags().StringVarP(&initPath, "path", "p", "", "The path to the directory that should be initialized as a globetrot project directory.")

	initCommand.MarkPersistentFlagRequired("path")

	rootCommand.AddCommand(initCommand)
}
