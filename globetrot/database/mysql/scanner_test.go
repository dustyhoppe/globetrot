package mysql

import (
	"testing"
)

func TestScanSimpleSelectStatementReturnsCorrectTokens(t *testing.T) {
	sql := "SELECT * FROM Course WHERE City = 'DEFAULT_CITY';"

	scanner := new(Scanner)
	scanner.Init(sql)

	tokens := scanner.Scan()

	if len(tokens) != 17 {
		t.Errorf("Incorrect number of tokens returned.")
	}
}

func TestScanMultiStatementReturnsCorrectTokenCount(t *testing.T) {
	sql := `SELECT * FROM Course WHERE City = 'DEFAULT_CITY';
	SELECT * FROM Course WHERE City = 'DEFAULT_CITY';
	`

	scanner := new(Scanner)
	scanner.Init(sql)

	tokens := scanner.Scan()
	expected := 35

	if len(tokens) != expected {
		t.Errorf("Incorrect number of tokens returns. Expected %d, received %d", expected, len(tokens))
	}
}

func TestScanStoredProcedureDropAndCreateReturnsCorrectTokenCount(t *testing.T) {
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
   
   DELIMITER ;
   
   `

	scanner := new(Scanner)
	scanner.Init(sql)

	tokens := scanner.Scan()
	expected := 101

	if len(tokens) != expected {
		t.Errorf("Incorrect number of tokens returns. Expected %d, received %d", expected, len(tokens))
	}
}
