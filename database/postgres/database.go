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

const create_scripts_run_table_script = `
CREATE OR REPLACE FUNCTION CreateGlobetrotScriptsRunTable(in schName varchar, in tblName varchar) RETURNS void AS $$
DECLARE 
    t_exists integer;
    t_user varchar(255);
    t_table varchar(255);
BEGIN
	SELECT INTO t_exists COUNT(*) FROM pg_tables WHERE schemaname = lower(schName) and tablename = lower(tblName);
	SELECT current_user into t_user;
	SELECT lower(schName) || '.' || lower(tblName) into t_table;
		
	IF t_exists = 0 THEN
		EXECUTE 'CREATE TABLE ' || t_table || '
		( 
			script_name		varchar(255)		NULL
			,hash		varchar(512)		NULL
			,date		timestamp		NOT NULL default current_timestamp
		);
		alter table ' || t_table || ' add constraint ' || replace(t_table, '.', '_') || '_pk' || ' primary key (script_name);
		GRANT SELECT ON TABLE ' || t_table || ' TO public;';
	END IF;
END;
$$ LANGUAGE 'plpgsql';
SELECT CreateGlobetrotScriptsRunTable('{0}', '{1}');
DROP FUNCTION CreateGlobetrotScriptsRunTable(in schName varchar, in tblName varchar);`

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
