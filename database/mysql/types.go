package mysql

import (
	"fmt"
)

type TokenType int
type StatementType int

const (
	// TokenTypes
	TOKEN_DELIMITER_DECLARE = iota
	TOKEN_DELIMITER
	TOKEN_COMMENT
	TOKEN_TEXT
	TOKEN_QUOTE
	TOKEN_ANSI_QUOTE
	TOKEN_SINGLE_QUOTE
	TOKEN_WHITESPACE
	TOKEN_END_OF_LINE
	TOKEN_END_OF_FILE

	// StatementTypes
	STATEMENT_SQL
	STATEMENT_DELIMITER
)

func PrintToken(t *Token) {

	tokenType := ""

	switch t.tokenType {
	case TOKEN_DELIMITER_DECLARE:
		tokenType = "DELIMITER_DECLARE"
		break
	case TOKEN_DELIMITER:
		tokenType = "DELIMITER"
		break
	case TOKEN_COMMENT:
		tokenType = "COMMNET"
		break
	case TOKEN_TEXT:
		tokenType = "TEXT"
		break
	case TOKEN_QUOTE:
		tokenType = "QUOTE"
		break
	case TOKEN_ANSI_QUOTE:
		tokenType = "ANSI_QUOTE"
		break
	case TOKEN_SINGLE_QUOTE:
		tokenType = "SINGLE_QUOTE"
		break
	case TOKEN_WHITESPACE:
		tokenType = "WHITESPACE"
		break
	case TOKEN_END_OF_LINE:
		tokenType = "END_OF_LINE"
		break
	case TOKEN_END_OF_FILE:
		tokenType = "END_OF_FILE"
		break
	default:
		tokenType = "UNKNOWN"
	}

	fmt.Printf("TOKEN: %s\nTYPE:  %s\n\n", t.value, tokenType)
}

func PrintStatement(s *ParsedStatement) {

	statementType := ""

	switch s.statementType {
	case STATEMENT_DELIMITER:
		statementType = "DELIMITER"
		break
	case STATEMENT_SQL:
		statementType = "SQL"
		break
	default:
		statementType = "UNKNOWN"
		break
	}

	fmt.Printf("STATEMENT: %s\nTYPE:  %s\n\n", s.value, statementType)
}
