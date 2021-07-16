package database

import (
	"fmt"
	"globetrot/database/common"
	"globetrot/database/mysql"
	"globetrot/database/postgres"
	"globetrot/database/sqlserver"
	"os"
)

type Database interface {
	Init(username string, password string, host string, port int, database string)
	CreateMigrationsTable()
	GetScriptRun(scriptName string) *common.ScriptRunRow
	ApplyScript(sql string, script_name string, sha string)
}

func NewDatabase(databaseType string, username string, password string, host string, port int, database string) Database {
	switch databaseType {
	case "mysql":
		database := new(mysql.MySqlDatabase)
		database.Init(username, password, host, port, databaseType)
		return database
	case "postgres":
		database := new(postgres.PostgresDatabase)
		database.Init(username, password, host, port, databaseType)
		return database
	case "sqlserver":
		database := new(sqlserver.SqlServerDatabase)
		database.Init(username, password, host, port, databaseType)
		return database
	default:
		fmt.Printf("ERROR: %s is not a supported database type\n", databaseType)
		os.Exit(1)
	}

	return nil
}
