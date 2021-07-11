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
	Connect(username string, password string, host string, port int, database string)
	CreateMigrationsTable()
	GetScriptRun(scriptName string) *common.ScriptRunRow
	ApplyScript(sql string, script_name string, sha string)
	Close()
}

func NewDatabase(databaseType string) Database {
	switch databaseType {
	case "mysql":
		return new(mysql.MySqlDatabase)
	case "postgres":
		return new(postgres.PostgresDatabase)
	case "sqlserver":
		return new(sqlserver.SqlServerDatabase)
	default:
		fmt.Printf("ERROR: %s is not a supported database type\n", databaseType)
		os.Exit(1)
	}

	return nil
}
