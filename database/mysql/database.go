package mysql

import (
	"database/sql"
	"fmt"
	"globetrot/database/common"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlDatabase struct {
	connection *sql.DB
}

func (mysql *MySqlDatabase) Connect(username string, password string, host string, port int, database string) {
	cs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", cs)

	if err != nil {
		panic(err.Error())
	}

	mysql.connection = db
}

func (mysql *MySqlDatabase) CreateMigrationsTable() {

	sql := `CREATE TABLE IF NOT EXISTS scripts_run (
				script_name VARCHAR(260) NOT NULL PRIMARY KEY,
				hash CHAR(44) NOT NULL,
				date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP()
			)`

	_, err := mysql.connection.Query(sql)

	if err != nil {
		panic(err.Error())
	}
}

func (mysql *MySqlDatabase) ApplyScript(sql string, script_name string, sha string) {

	_, err := mysql.connection.Query(sql)

	if err != nil {
		panic(err.Error())
	}

	_, err = mysql.connection.Query(fmt.Sprintf("INSERT INTO scripts_run (script_name, hash) VALUES ( '%s', '%s' );", script_name, sha))
	if err != nil {
		panic(err.Error())
	}
}

func (mysql *MySqlDatabase) GetScriptRun(scriptName string) *common.ScriptRunRow {

	sql := fmt.Sprintf("SELECT script_name AS ScriptName, hash AS Hash FROM scripts_run WHERE script_name = '%s'", scriptName)
	rows, err := mysql.connection.Query(sql)
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

func (mysql *MySqlDatabase) Close() {
	mysql.connection.Close()
}
