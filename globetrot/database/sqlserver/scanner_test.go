package sqlserver

import (
	"testing"
)

func TestScanner(t *testing.T) {
	sql := `SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';

	go
	
	SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';
	`

	scanner := new(Scanner)
	scanner.Init(sql)

	tokens := scanner.Scan()

	for _, token := range tokens {
		PrintToken(token)
	}
}

func TestParser(t *testing.T) {
	sql := `SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';

	go
	
	SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';
	`

	parser := new(Parser)
	parser.Init(sql)

	tokens := parser.Parse()

	for _, token := range tokens {
		PrintStatement(token)
	}
}
