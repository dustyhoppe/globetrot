package postgres

import (
	"database/sql"
	"fmt"
	"globetrot/database/common"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresDatabase struct {
	connection *sql.DB
	database   string
}

func (postgres *PostgresDatabase) Connect(username string, password string, host string, port int, database string) {
	postgres.database = database
	cs := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", username, password, host, port, database)
	db, err := sql.Open("pgx", cs)

	if err != nil {
		panic(err.Error())
	}

	postgres.connection = db
}

func (postgres *PostgresDatabase) CreateMigrationsTable() {

	sql := `
	CREATE TABLE IF NOT EXISTS scripts_run
	(
		script_name VARCHAR(255) NOT NULL PRIMARY KEY,
		hash VARCHAR(512) NOT NULL,
		date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	GRANT SELECT ON scripts_run TO public;
	`

	_, err := postgres.connection.Exec(sql)

	if err != nil {
		panic(err.Error())
	}
}

func (postgres *PostgresDatabase) ApplyScript(sql string, script_name string, sha string) {

	_, err := postgres.connection.Exec(sql)

	if err != nil {
		panic(err.Error())
	}

	script_sql := `INSERT INTO scripts_run (script_name, hash) VALUES ( '%s', '%s' )
		ON CONFLICT (script_name) DO UPDATE SET hash='%s', date=CURRENT_TIMESTAMP;`

	_, err = postgres.connection.Exec(fmt.Sprintf(script_sql, script_name, sha, sha))
	if err != nil {
		panic(err.Error())
	}
}

func (postgres *PostgresDatabase) GetScriptRun(scriptName string) *common.ScriptRunRow {

	sql := fmt.Sprintf("SELECT script_name AS ScriptName, hash AS Hash FROM scripts_run WHERE script_name = '%s'", scriptName)
	rows, err := postgres.connection.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	if rows.Next() {
		row := common.ScriptRunRow{}
		err = rows.Scan(&row.ScriptName, &row.Hash)
		if err != nil {
			panic(err.Error())
		}

		return &row
	}

	return nil
}

func (postgres *PostgresDatabase) Close() {
	postgres.connection.Close()
}
