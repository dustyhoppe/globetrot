package sqlserver

import (
	"testing"
)

func TestScanSimpleSelectStatementReturnsCorrectTokens(t *testing.T) {
	sql := "SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';"

	scanner := new(Scanner)
	scanner.Init(sql)

	tokens := scanner.Scan()

	if len(tokens) != 17 {
		t.Errorf("Incorrect number of tokens returned.")
	}
}

func TestScanMultiBatchStatementReturnsCorrectTokenCount(t *testing.T) {
	sql := `SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';
	GO
	SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';
	`

	scanner := new(Scanner)
	scanner.Init(sql)

	tokens := scanner.Scan()
	expected := 37

	if len(tokens) != expected {
		t.Errorf("Incorrect number of tokens returns. Expected %d, received %d", expected, len(tokens))
	}
}

func TestScanStoredProcedureDropAndCreateReturnsCorrectTokenCount(t *testing.T) {
	sql := `-- =============================================
	-- Author:		Dusty Hoppe
	-- Create date: 07-05-2021
	-- Description: Inserts a record into the not_users table
	-- =============================================
	
	IF EXISTS (SELECT * FROM sys.objects WHERE type = 'P' AND name = 'Users_Insert')
		DROP PROCEDURE Users_Insert
	GO
	
	CREATE PROCEDURE Users_Insert
		-- Add the parameters for the stored procedure here
		@pUsername VARCHAR(50),
		@pPassword VARCHAR(100)
	AS
	BEGIN
		-- SET NOCOUNT ON added to prevent extra result sets from
		-- interfering with SELECT statements.
		SET NOCOUNT ON;
	
		-- Insert statements for procedure here
		INSERT INTO not_users ( Username, [Password] )
			VALUES ( @pUsername, @pPassword );
	
	END`

	scanner := new(Scanner)
	scanner.Init(sql)

	tokens := scanner.Scan()
	expected := 119

	if len(tokens) != expected {
		t.Errorf("Incorrect number of tokens returns. Expected %d, received %d", expected, len(tokens))
	}
}

func TestScanStatementWithStarCommentReturnsCorrectTokenCount(t *testing.T) {
	sql := `/*
	* Author:		Dusty Hoppe
	* Create date: 07-05-2021
	* Description: Inserts a record into the not_users table
	*/
	SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';
	`

	scanner := new(Scanner)
	scanner.Init(sql)

	tokens := scanner.Scan()
	expected := 19

	if len(tokens) != expected {
		t.Errorf("Incorrect number of tokens returns. Expected %d, received %d", expected, len(tokens))
	}
}
