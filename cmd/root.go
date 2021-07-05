package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const envPrefix = "GLOBETROT"

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

var rootCommand = &cobra.Command{
	Use:   "globetrot [subcommand]",
	Short: "A database change management tool",
	Long:  "Globetrot is a CLI tool used for managing database changes/migrations",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Globetrot") },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
}

func init() {

	rootCommand.PersistentFlags().StringVarP(&configFile, "configPath", "c", "", "The path to the configuration to use in-place of passing command line arguments")
	rootCommand.PersistentFlags().StringVarP(&username, "username", "u", "root", "The username to use when connecting to the database")
	rootCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "The password to use when connecting to the database")
	rootCommand.PersistentFlags().StringVarP(&host, "server", "s", "localhost", "The host server the database is located at")
	rootCommand.PersistentFlags().IntVarP(&port, "port", "P", 3306, "The port the database server the database is listening at")
	rootCommand.PersistentFlags().StringVarP(&database, "database", "d", "", "The name of the database")
	rootCommand.PersistentFlags().StringVarP(&filePath, "filePath", "f", ".\\", "The directory where your SQL scripts are located")
	rootCommand.PersistentFlags().BoolVar(&dryRun, "dryRun", false, "Indicates whether the command should perform a dry run")
	rootCommand.PersistentFlags().StringVarP(&databaseType, "databaseType", "t", "mysql", "Indicates the database type/platform (mysql, postgres)")
	rootCommand.PersistentFlags().StringVarP(&environment, "env", "e", "", "The environment the migration is targeting")

	rootCommand.MarkPersistentFlagRequired("username")
	rootCommand.MarkPersistentFlagRequired("password")
	rootCommand.MarkPersistentFlagRequired("server")
	rootCommand.MarkPersistentFlagRequired("port")
	rootCommand.MarkPersistentFlagRequired("database")
	rootCommand.MarkPersistentFlagRequired("filePath")
	rootCommand.MarkPersistentFlagRequired("databaseType")
}

func initializeConfig(cmd *cobra.Command) error {

	v := viper.New()

	v.SetConfigName("globetrot")

	v.AddConfigPath(configFile)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
