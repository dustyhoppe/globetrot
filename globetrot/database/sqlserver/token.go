package sqlserver

type Token struct {
	value     string
	line      int
	column    int
	tokenType TokenType
}

func (t *Token) Init(value string, line int, column int, tokenType TokenType) {
	t.value = value
	t.line = line
	t.column = column
	t.tokenType = tokenType
}
