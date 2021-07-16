package mysql

import (
	"testing"
)

func TestParserSimpleSelectStatementReturnsCorrectStatementCount(t *testing.T) {
	sql := "SELECT * FROM Course WHERE City = 'DEFAULT_CITY';"

	parser := new(Parser)
	parser.Init(sql)

	statements := parser.Parse()
	expected := 1

	if len(statements) != expected {
		t.Errorf("Incorrect number of statements parsed. Expected %d, received %d", expected, len(statements))
	}

	if statements[0].value != "SELECT * FROM Course WHERE City = 'DEFAULT_CITY'" {
		t.Errorf("Single statement parsed incorrectly.")
	}
}

func TestParseMultiStatementReturnsCorrectStatementCount(t *testing.T) {
	sql := `SELECT * FROM Course WHERE City = 'DEFAULT_CITY';
	SELECT * FROM Course WHERE City = 'DEFAULT_CITY';`

	parser := new(Parser)
	parser.Init(sql)

	statements := parser.Parse()
	expected := 2

	if len(statements) != expected {
		t.Errorf("Incorrect number of statements parsed. Expected %d, received %d", expected, len(statements))
	}
}

func TestPraseStoredProcedureWithDelimiterDeclaration(t *testing.T) {
	sql := `/*
	* Author:		Dusty Hoppe
	* Create date: 07-05-2021
	* Description: Inserts a record into the not_users table
	*/
   
   DROP PROCEDURE IF EXISTS Users_Insert;
   
   DELIMITER $$
   
   CREATE PROCEDURE Users_Insert (
	   pUsername VARCHAR(50), 
	   pPassword VARCHAR(100)
   )
   BEGIN
	   INSERT INTO users ( Username, Password )
		   VALUES ( pUsername, pPassword );
   
   END$$
   
   DELIMITER ;`

	parser := new(Parser)
	parser.Init(sql)

	statements := parser.Parse()
	expected := 2

	if len(statements) != expected {
		t.Errorf("Incorrect number of statements parsed. Expected %d, received %d", expected, len(statements))
	}

}
