package sqlserver

import (
	"fmt"
	"strings"
	"unicode"
)

const DELIMITER = ";"

type Scanner struct {
	script  string
	start   int
	current int
	line    int
	tokens  []*Token
}

func (s *Scanner) Init(script string) {
	s.script = script
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

	// give priority to the batch separator "GO"
	if (c == 'G' || c == 'g') && (s.Peek() == 'O' || s.Peek() == 'o') {
		if s.BatchSeparator() {
			return
		}
	}

	switch c {
	case ';':
		s.Delimiter()
		break
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
	case '\'':
		s.SingleQuoted()
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

func (s *Scanner) BatchSeparator() bool {
	separator := "GO"
	if s.PeekPeek() == '\n' || (s.PeekPeek() == '\r' && s.PeekPeekPeek() == '\n') {
		possible := s.script[s.start : s.start+len(separator)]
		if strings.ToUpper(possible) == separator {
			s.current = s.start + len(separator)
			s.AddToken(TOKEN_BATCH_SEPARATER)
			return true
		}
	}

	return false
}

func (s *Scanner) Delimiter() {
	s.AddToken(TOKEN_DELIMITER)
}

func (s *Scanner) General() {
	for !s.IsDone() && s.IsLetterOrDigit(s.Peek()) {
		if s.IsLetterOrDigit(rune(DELIMITER[0])) {
			if s.Peek() == rune(DELIMITER[0]) {
				break
			}
		}

		s.NextChar()
	}

	s.AddToken(TOKEN_TEXT)
}

// Consumes a single line comment
func (s *Scanner) Comment() {
	for !(s.Peek() == '\r' && s.PeekPeek() == '\n') && s.Peek() != '\n' && !s.IsDone() {
		s.NextChar()
	}

	s.AddToken(TOKEN_DASH_COMMENT)
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

		s.NextChar() // consume *
		s.NextChar() // consume /

		s.AddToken(TOKEN_STAR_COMMENT)
	}
}

/// Returns the current character and advances the current index
func (s *Scanner) NextChar() rune {
	s.current++
	return rune(s.script[s.current-1])
}

/// Consumes whitespace characters until a non-whitespace character is encountered
func (s *Scanner) Whitespace() {
	for !s.IsDone() && s.Peek() != '\n' && unicode.IsSpace(s.Peek()) {
		s.NextChar()
	}
	s.AddToken(TOKEN_WHITESPACE)
}

func (s *Scanner) SingleQuoted() {
	if !s.IsDone() && !s.IsSingleQuote(s.Peek()) {
		if s.Peek() == '\n' {
			s.line++
		}

		s.NextChar()

		if s.IsDone() {
			panic(fmt.Sprintf("Unterminated single quoted value on line %d", s.line))
		}
	}

	s.NextChar()
	s.AddToken(TOKEN_QUOTE)
}

/// Determines whether a given rune is a single quote
func (s Scanner) IsSingleQuote(r rune) bool {
	return r == '\''
}

/// Returns a flag indicating whether the current rune is
// equal to the passed in rune
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

/// Returns a flag indicating whether the entire script has been scanned
func (s *Scanner) IsDone() bool {
	return s.current >= len(s.script)
}

/// Returns the next current character without advancing the index
/// Returns a nil rune in the event the script has been completely scanned
func (s *Scanner) Peek() rune {
	if s.IsDone() {
		return rune(0)
	}
	return rune(s.script[s.current])
}

/// Returns the next, next character without advancing the current index
/// Returns a null rune in the event the script has been completely scanned
func (s *Scanner) PeekPeek() rune {
	if s.IsDone() || (s.current+1 >= len(s.script)) {
		return rune(0)
	}
	return rune(s.script[s.current+1])
}

/// Returns the next, next, next character without advancing the current index
/// Returns a null rune in the event the script has been completely scanned
func (s *Scanner) PeekPeekPeek() rune {
	if s.IsDone() || (s.current+2 >= len(s.script)) {
		return rune(0)
	}
	return rune(s.script[s.current+2])
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
	token := new(Token)
	token.Init(value, s.line, s.start, tokenType)

	PrintToken(token)

	s.tokens = append(s.tokens, token)
}

func (s *Scanner) AddEmptyToken(tokenType TokenType) {
	token := new(Token)
	token.Init("", s.line, s.start, tokenType)

	s.tokens = append(s.tokens, token)
}

/// Returns a flag indicating whether a given rune is a letter or digit
func (s Scanner) IsLetterOrDigit(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
