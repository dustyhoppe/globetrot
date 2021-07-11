package sqlserver

type ParsedBatch struct {
	statementType StatementType
	value         string
	delimiter     string
}

func (p *ParsedBatch) Init(statementType StatementType, value string, delimiter string) {
	p.statementType = statementType
	p.value = value
	p.delimiter = delimiter
}
