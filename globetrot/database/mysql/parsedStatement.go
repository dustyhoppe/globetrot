package mysql

type ParsedStatement struct {
	statementType StatementType
	value         string
	delimiter     string
}

func (p *ParsedStatement) Init(statementType StatementType, value string, delimiter string) {
	p.statementType = statementType
	p.value = value
	p.delimiter = delimiter
}
