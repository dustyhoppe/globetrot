package postgres

import (
	"database/sql"
	"fmt"
	"globetrot/database/common"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresDatabase struct {
	database   string
	username   string
	password   string
	host       string
	port       int
	connection *sql.DB
}

func (postgres *PostgresDatabase) Open() {
	cs := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", postgres.username, postgres.password, postgres.host, postgres.port, postgres.database)
	db, err := sql.Open("pgx", cs)

	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(time.Minute)

	postgres.connection = db
}

func (postgres *PostgresDatabase) Init(username string, password string, host string, port int, database string) {
	postgres.username = username
	postgres.password = password
	postgres.host = host
	postgres.port = port
	postgres.database = database
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
        defer rows.Close()

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
	defer postgres.connection.Close()
}
