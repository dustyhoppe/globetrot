package mysql

import (
	"fmt"
	"strings"
	"unicode"
)

const DELIMITER_DECLARE = "delimiter"
const DEFAULT_DELIMITER = ";"

type Scanner struct {
	script     string
	delimiter  string
	ansiQuotes bool
	start      int
	current    int
	line       int
	tokens     []*Token
}

func (s *Scanner) Init(script string) {
	s.script = script
	s.delimiter = DEFAULT_DELIMITER
}

func (s *Scanner) Scan() []*Token {

	s.start = 0
	s.current = 0
	s.line = 1

	for !s.IsDone() {
		s.start = s.current
		s.ScanToken()
	}

	s.EndOfFile()
	return s.tokens
}

func (s *Scanner) ScanToken() {

	c := s.NextChar()

	if len(s.delimiter) == 1 && c == rune(s.delimiter[0]) {
		s.SingleCharacterDelimiter()
		return
	} else if len(s.delimiter) == 2 && c == rune(s.delimiter[0]) && s.Peek() == rune(s.delimiter[1]) {
		s.TwoCharacterDelimiter()
		return
	} else if len(s.delimiter) > 2 && c == rune(s.delimiter[0]) && s.Peek() == rune(s.delimiter[1]) {
		if s.MultiCharacterDelimiter() {
			return
		}
	}

	switch c {
	case '\t':
	case ' ':
		s.Whitespace()
		break
	case '\r':
		if s.Peek() == '\n' {
			s.NextChar()
			s.EndOfLine()
		}
		break
	case '\n':
		s.EndOfLine()
		break
	case '"':
		if s.ansiQuotes {
			s.AnsiQuoted()
			break
		}
		break
	case '\'':
		s.SingleQuoted()
		break
	case '`':
		s.Quoted()
		break
	case '#':
		s.Comment()
		break
	case '-':
		if s.Match('-') {
			s.NextChar()
			s.Comment()
			break
		}
		break
	case '/':
		if s.Match('*') {
			s.NextChar()
			s.MultiLineComment()
		} else {
			s.General()
		}
		break
	default:
		s.General()
		break
	}
}

func (s *Scanner) SingleCharacterDelimiter() {
	s.AddToken(TOKEN_DELIMITER)
}

func (s *Scanner) TwoCharacterDelimiter() {
	s.NextChar() // consume second character
	s.AddToken(TOKEN_DELIMITER)
}

func (s *Scanner) MultiCharacterDelimiter() bool {
	if s.start+len(s.delimiter) <= len(s.script) {
		possible := s.script[s.start : s.start+len(s.delimiter)]
		if possible == s.delimiter {
			s.current = s.start + len(s.delimiter) // consume the rest of the delimiter
			s.AddToken(TOKEN_DELIMITER)
			return true
		}
	}

	return false
}

/// Returns the current char (byte) and advances the current index
func (s *Scanner) NextChar() rune {

	s.current++
	return rune(s.script[s.current-1])
}

/// Returns a flag indicating whether the entire script has been scanned
func (s *Scanner) IsDone() bool {
	return s.current >= len(s.script)
}

/// Returns the next current char without advancing the index
/// Returns a null rune in the event the script has been completely scanned
func (s *Scanner) Peek() rune {
	if s.IsDone() {
		return rune(0)
	}

	return rune(s.script[s.current])
}

/// Returns the next, next current char without advancing the index
/// Returns a null rune in the event the script has been completely scanned
func (s *Scanner) PeekPeek() rune {
	if s.IsDone() || (s.current+1 >= len(s.script)) {
		return rune(0)
	}

	return rune(s.script[s.current+1])
}

func (s *Scanner) PeekMultiCharacterDelimiter() bool {
	if s.current+len(s.delimiter) <= len(s.script) {
		possible := s.script[s.current : len(s.delimiter)+s.current]
		if possible == s.delimiter {
			return true
		}
	}
	return false
}

/// Consumes whitespace characters until a non-whitespace character is encountered
func (s *Scanner) Whitespace() {
	for !s.IsDone() && s.Peek() != '\n' && unicode.IsSpace(s.Peek()) {
		s.NextChar()
	}

	s.AddToken(TOKEN_WHITESPACE)
}

func (s *Scanner) General() {
	for !s.IsDone() && s.IsGeneralCharacter(s.Peek()) {
		if s.IsGeneralCharacter(rune(s.delimiter[0])) {

			if len(s.delimiter) == 1 && s.Peek() == rune(s.delimiter[0]) {
				break
			} else if len(s.delimiter) == 2 && s.Peek() == rune(s.delimiter[0]) && s.PeekPeek() == rune(s.delimiter[1]) {
				break
			} else if len(s.delimiter) > 2 && s.PeekMultiCharacterDelimiter() {
				break
			}
		}

		s.NextChar()
	}

	s.AddToken((TOKEN_TEXT))
}

/// Consumes a multi-line comment and adds the token to the collection
func (s *Scanner) MultiLineComment() {
	for !(s.Match('*') && s.Peek() == '/') {
		if s.Peek() == '\n' {
			s.line++
		}

		s.NextChar()

		if s.IsDone() {
			panic(fmt.Sprintf("Unterminated comment on line %d", s.line))
		}
	}

	s.NextChar() // consume *
	s.NextChar() // consume /

	s.AddToken(TOKEN_COMMENT)
}

func (s *Scanner) Comment() {
	for !(s.Peek() == '\r' && s.PeekPeek() == '\n') && s.Peek() != '\n' && !s.IsDone() {
		s.NextChar()
	}

	s.AddToken(TOKEN_COMMENT)
}

func (s *Scanner) Quoted() {
	for !s.IsDone() && !s.IsQuote(s.Peek()) {
		if s.Peek() == '\n' {
			s.line++
		}

		s.NextChar()

		if s.IsDone() {
			panic(fmt.Sprintf("Unterminated quoted value on line %d", s.line))
		}
	}

	s.NextChar()
	s.AddToken(TOKEN_QUOTE)
}

func (s *Scanner) SingleQuoted() {
	for !s.IsDone() && !s.IsSingleQuote(s.Peek()) {
		if s.Peek() == '\n' {
			s.line++
		}

		s.NextChar()

		if s.IsDone() {
			panic(fmt.Sprintf("Unterminated single quoted value on line %d", s.line))
		}
	}

	s.NextChar()
	s.AddToken(TOKEN_SINGLE_QUOTE)
}

func (s *Scanner) AnsiQuoted() {
	for !s.IsDone() && !s.IsAnsiQuote(s.Peek()) {
		if s.Peek() == '\n' {
			s.line++
		}

		s.NextChar()

		if s.IsDone() {
			panic(fmt.Sprintf("Unterminated double quoted value on line %d", s.line))
		}
	}

	s.NextChar()
	s.AddToken(TOKEN_ANSI_QUOTE)
}

/// Determines whether a given rune is a quote
func (s Scanner) IsQuote(r rune) bool {
	return r == '`'
}

/// Determines whether a given rune is an ANSI quote
func (s *Scanner) IsAnsiQuote(r rune) bool {
	return r == '"'
}

/// Determines whether a given rune is a single quote
func (s Scanner) IsSingleQuote(r rune) bool {
	return r == '\''
}

/// Returns a flag indicating whether the current rune is
/// equals to the passed in rune
func (s *Scanner) Match(r rune) bool {
	if s.IsDone() {
		return false
	}

	if rune(s.script[s.current]) != r {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) EndOfLine() {
	s.line++
	s.AddEmptyToken(TOKEN_END_OF_LINE)
}

func (s *Scanner) EndOfFile() {
	s.AddEmptyToken(TOKEN_END_OF_FILE)
}

func (s *Scanner) AddToken(tokenType TokenType) {
	end := s.current - s.start
	if end < 0 {
		end = 0
	}

	if s.current+end < 1 {
		return
	}

	value := s.script[s.start:s.current]
	if strings.ToLower(value) == DELIMITER_DECLARE {
		s.DelimiterDeclaration()
		return
	}

	token := new(Token)
	token.Init(value, s.line, s.start, tokenType)

	PrintToken(token)

	s.tokens = append(s.tokens, token)
}

func (s *Scanner) DelimiterDeclaration() {
	value := s.script[s.start:s.current]
	token := new(Token)
	token.Init(value, s.line, s.start, TOKEN_DELIMITER_DECLARE)

	PrintToken(token)

	s.tokens = append(s.tokens, token)

	// move our start forward for next scan
	s.start = s.current

	// Consume any whitespace
	s.Whitespace()

	if s.IsDone() {
		panic("Delmiter keyword used but no delimiter was provided")
	}

	// consume the new delimiter
	for s.Peek() != '\n' && !s.IsDone() {
		s.NextChar()
	}

	// add the new delimiter token
	value = strings.TrimSpace(s.script[s.start:s.current])
	delimiter_token := new(Token)
	delimiter_token.Init(value, s.line, s.start, TOKEN_DELIMITER)
	s.tokens = append(s.tokens, delimiter_token)

	PrintToken(delimiter_token)

	// set the new delimiter
	s.delimiter = value
}

///
func (s *Scanner) AddEmptyToken(tokenType TokenType) {
	token := new(Token)
	token.Init("", s.line, s.start, tokenType)

	s.tokens = append(s.tokens, token)
}

/// Returns a flag indicating whether a given rune is a letter or digit
func (s Scanner) IsGeneralCharacter(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
