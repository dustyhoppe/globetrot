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
	Close()
	Open()
}

func NewDatabase(databaseType string, username string, password string, host string, port int, database string) Database {
	switch databaseType {
	case "mysql":
		db := new(mysql.MySqlDatabase)
		db.Init(username, password, host, port, database)
		return db
	case "postgres":
		db := new(postgres.PostgresDatabase)
		db.Init(username, password, host, port, database)
		return db
	case "sqlserver":
		db := new(sqlserver.SqlServerDatabase)
		db.Init(username, password, host, port, database)
		return db
	default:
		fmt.Printf("ERROR: %s is not a supported database type\n", databaseType)
		os.Exit(1)
	}

	return nil
}
