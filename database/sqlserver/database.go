package sqlserver

import (
	"database/sql"
	"fmt"
	"globetrot/database/common"

	_ "github.com/denisenkom/go-mssqldb"
)

type SqlServerDatabase struct {
	connection *sql.DB
	database   string
}

func (sqlserver *SqlServerDatabase) Connect(username string, password string, host string, port int, database string) {
	sqlserver.database = database
	cs := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", username, password, host, port, database)
	//cs := "sqlserver://sa:Onicepo1!@172.21.240.1:1433/sqlexpress?database=Globetrot"
	db, err := sql.Open("sqlserver", cs)

	fmt.Println(cs)

	if err != nil {
		panic(err.Error())
	}

	sqlserver.connection = db
}

func (sqlserver *SqlServerDatabase) CreateMigrationsTable() {

	sql := `
	IF NOT EXISTS ( SELECT 1 FROM sys.tables )
	BEGIN

		CREATE TABLE scripts_run
		(
			script_name VARCHAR(255) NOT NULL PRIMARY KEY,
			[hash] VARCHAR(512) NOT NULL,
			[date] DATETIME NOT NULL DEFAULT GETUTCDATE()
		);

	END;
	`

	_, err := sqlserver.connection.Exec(sql)

	if err != nil {
		panic(err.Error())
	}
}

func (sqlserver *SqlServerDatabase) ApplyScript(sql string, script_name string, sha string) {

	_, err := sqlserver.connection.Exec(sql)

	if err != nil {
		panic(err.Error())
	}

	script_sql := `
	UPDATE scripts_run
		SET hash='%s',
			date=GETUTCDATE()
		WHERE script_name='%s';
	
	IF @@ROWCOUNT = 0
		INSERT INTO scripts_run (script_name, hash) VALUES ( '%s', '%s' );`

	_, err = sqlserver.connection.Exec(fmt.Sprintf(script_sql, sha, script_name, script_name, sha))
	if err != nil {
		panic(err.Error())
	}
}

func (sqlserver *SqlServerDatabase) GetScriptRun(scriptName string) *common.ScriptRunRow {

	sql := fmt.Sprintf("SELECT script_name AS ScriptName, hash AS Hash FROM scripts_run WHERE script_name = '%s'", scriptName)
	rows, err := sqlserver.connection.Query(sql)
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

func (sqlserver *SqlServerDatabase) Close() {
	sqlserver.connection.Close()
}
