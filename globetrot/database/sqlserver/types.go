package sqlserver

import (
	"fmt"
)

const print_debug_enabled = false

type TokenType int
type StatementType int

const (
	// Token Types
	TOKEN_BATCH_SEPARATER = iota
	TOKEN_DELIMITER
	TOKEN_DASH_COMMENT
	TOKEN_STAR_COMMENT
	TOKEN_WHITESPACE
	TOKEN_END_OF_LINE
	TOKEN_END_OF_FILE
	TOKEN_QUOTE
	TOKEN_TEXT

	// Statement Types
	STATEMENT_SQL
	STATEMENT_BATCH_SEPARATOR
)

func PrintToken(t *Token) {

	tokenType := ""

	switch t.tokenType {
	case TOKEN_BATCH_SEPARATER:
		tokenType = "BATCH_SEPARATOR"
		break
	case TOKEN_DELIMITER:
		tokenType = "DELIMITER"
		break
	case TOKEN_DASH_COMMENT:
		tokenType = "DASH_COMMENT"
		break
	case TOKEN_STAR_COMMENT:
		tokenType = "STAR_COMMENT"
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
	case TOKEN_QUOTE:
		tokenType = "QUOTE"
		break
	case TOKEN_TEXT:
		tokenType = "TEXT"
		break
	default:
		tokenType = "UNKNOWN"
	}

	if print_debug_enabled {
		fmt.Printf("TOKEN: %s\nTYPE:  %s\n\n", t.value, tokenType)
	}
}

func PrintStatement(s *ParsedBatch) {

	statementType := ""

	switch s.statementType {
	case STATEMENT_BATCH_SEPARATOR:
		statementType = "BATCH_SEPARATOR"
		break
	case STATEMENT_SQL:
		statementType = "SQL"
		break
	default:
		statementType = "UNKNOWN"
		break
	}

	if print_debug_enabled {
		fmt.Printf("STATEMENT: %s\nTYPE:  %s\n\n", s.value, statementType)
	}
}
