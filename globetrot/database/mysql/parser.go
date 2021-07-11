package mysql

import "strings"

type Parser struct {
	script     string
	scanner    *Scanner
	delimiter  string
	tokens     []*Token
	statements []*ParsedStatement
	ansiQuotes bool
	start      int
	current    int
}

func (p *Parser) Init(script string) {
	p.delimiter = ";"
	p.script = script
	p.scanner = new(Scanner)
	p.scanner.Init(script)
	p.scanner.ansiQuotes = p.ansiQuotes
}

func (p *Parser) Parse() []*ParsedStatement {

	p.start = 0
	p.current = 0

	p.tokens = p.scanner.Scan()

	for !p.IsDone() {
		p.start = p.current

		statement := p.ParseStatement()

		if len(strings.TrimSpace(statement.value)) > 0 && statement.statementType != STATEMENT_DELIMITER {
			p.statements = append(p.statements, statement)
		}
	}

	return p.statements
}

func (p *Parser) ParseStatement() *ParsedStatement {
	setDelimiter := false
	statementEnd := false
	statementType := STATEMENT_SQL

	for !p.IsDone() && !statementEnd {
		token := p.NextToken()

		PrintToken(token)

		switch token.tokenType {

		case TOKEN_END_OF_FILE:
			break

		case TOKEN_DELIMITER_DECLARE:
			statementType = STATEMENT_DELIMITER
			setDelimiter = true
			break

		case TOKEN_DELIMITER:
			if setDelimiter {
				p.delimiter = token.value
				setDelimiter = false
			}

			if p.Peek().tokenType == TOKEN_END_OF_LINE {
				p.NextToken()
			}

			statementEnd = true
			break

		default:
			break
		}
	}

	var sb strings.Builder
	for _, token := range p.tokens[p.start:p.current] {
		if token.tokenType == TOKEN_END_OF_LINE {
			// Append new line
			sb.WriteRune('\n')
		}
		if token.tokenType == TOKEN_DELIMITER {
			// the end of the statement doesn't require a delimiter
		} else {
			sb.WriteString(token.value)
		}
	}

	parsedStatement := &ParsedStatement{
		delimiter:     p.delimiter,
		statementType: StatementType(statementType),
		value:         sb.String(),
	}

	PrintStatement(parsedStatement)

	return parsedStatement
}

func (p *Parser) NextToken() *Token {
	p.current++
	return p.tokens[p.current-1]
}

func (p *Parser) Peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) IsDone() bool {
	return p.current >= len(p.tokens)
}
