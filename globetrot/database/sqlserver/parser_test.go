package sqlserver

import (
	"testing"
)

func TestParserSimpleSelectStatementReturnsCorrectStatementCount(t *testing.T) {
	sql := "SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';"

	parser := new(Parser)
	parser.Init(sql)

	statements := parser.Parse()
	expected := 1

	if len(statements) != expected {
		t.Errorf("Incorrect number of statements parsed. Expected %d, received %d", expected, len(statements))
	}

	if statements[0].value != "SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';" {
		t.Errorf("Single statement parsed incorrectly.")
	}
}

func TestParseMultiBatchStatementReturnsCorrectStatementCount(t *testing.T) {
	sql := `SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';
	GO
	SELECT * FROM dbo.Course WHERE City = 'DEFAULT_CITY';`

	parser := new(Parser)
	parser.Init(sql)

	batches := parser.Parse()
	expected := 2

	if len(batches) != expected {
		t.Errorf("Incorrect number of batches parsed. Expected %d, received %d", expected, len(batches))
	}
}

func TestPraseStoredProcedureWithDelimiterDeclaration(t *testing.T) {
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

	parser := new(Parser)
	parser.Init(sql)

	batches := parser.Parse()
	expected := 2

	if len(batches) != expected {
		t.Errorf("Incorrect number of batches parsed. Expected %d, received %d", expected, len(batches))
	}

}
