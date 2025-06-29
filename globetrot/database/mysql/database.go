package mysql

import (
	"database/sql"
	"fmt"
	"globetrot/database/common"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlDatabase struct {
	database   string
	username   string
	password   string
	host       string
	port       int
	connection *sql.DB
}

func (mysql *MySqlDatabase) Open() {
	cs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true&autocommit=true", mysql.username, mysql.password, mysql.host, mysql.port, mysql.database)
	db, err := sql.Open("mysql", cs)

	if err != nil {
		panic(err.Error())
	}

	mysql.connection = db
}

func (mysql *MySqlDatabase) Init(username string, password string, host string, port int, database string) {
	mysql.username = username
	mysql.password = password
	mysql.host = host
	mysql.port = port
	mysql.database = database
}

func (mysql *MySqlDatabase) CreateMigrationsTable() {

	sql := `CREATE TABLE IF NOT EXISTS scripts_run (
				script_name VARCHAR(260) NOT NULL PRIMARY KEY,
				hash CHAR(44) NOT NULL,
				date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP()
			)`

	_, err := mysql.connection.Exec(sql)

	if err != nil {
		panic(err.Error())
	}
}

func (mysql *MySqlDatabase) ApplyScript(sql string, script_name string, sha string) {

	parser := new(Parser)
	parser.Init(sql)
	statements := parser.Parse()

	for _, statement := range statements {
		_, err := mysql.connection.Exec(statement.value)

		if err != nil {
			panic(err.Error())
		}
	}

	script_sql := `INSERT INTO scripts_run (script_name, hash) VALUES ( '%s', '%s' )
		ON DUPLICATE KEY UPDATE hash='%s', date=CURRENT_TIMESTAMP();`

	_, err := mysql.connection.Exec(fmt.Sprintf(script_sql, script_name, sha, sha))
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

func (mysql *MySqlDatabase) Close() {
	defer mysql.connection.Close()
}
