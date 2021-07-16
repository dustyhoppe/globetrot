package sqlserver

import "strings"

type Parser struct {
	script   string
	scanner  *Scanner
	delmiter string
	tokens   []*Token
	batches  []*ParsedBatch
	start    int
	current  int
}

func (p *Parser) Init(script string) {
	p.script = script
	p.scanner = new(Scanner)
	p.scanner.Init(script)
}

func (p *Parser) Parse() []*ParsedBatch {

	p.start = 0
	p.current = 0

	p.tokens = p.scanner.Scan()

	for !p.IsDone() {
		p.start = p.current
		batch := p.ParseBatch()

		if len(strings.TrimSpace(batch.value)) > 0 && batch.statementType != STATEMENT_BATCH_SEPARATOR {
			p.batches = append(p.batches, batch)
		}
	}

	return p.batches
}

func (p *Parser) ParseBatch() *ParsedBatch {
	batchEnd := false
	statementType := STATEMENT_SQL

	for !p.IsDone() && !batchEnd {
		token := p.NextToken()

		PrintToken(token)

		switch token.tokenType {
		case TOKEN_END_OF_FILE:
			break

		case TOKEN_BATCH_SEPARATER:
			if p.Peek().tokenType == TOKEN_END_OF_LINE {
				p.NextToken()
			}

			batchEnd = true
			break

		default:
			break
		}
	}

	var sb strings.Builder
	for _, token := range p.tokens[p.start:p.current] {
		if token.tokenType == TOKEN_END_OF_LINE {
			sb.WriteRune('\n')
		}
		if token.tokenType == TOKEN_BATCH_SEPARATER {
			// the end of the statement doesn't require a delimiter
		} else {
			sb.WriteString(token.value)
		}
	}

	parsedBatch := &ParsedBatch{
		statementType: StatementType(statementType),
		value:         sb.String(),
	}

	PrintBatch(parsedBatch)

	return parsedBatch
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
