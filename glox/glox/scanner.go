package glox

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	source  string  // The source code being scanned.
	tokens  []Token // List of scanned tokens.
	start   int     // Index of the first character in the lexeme being scanned.
	current int     // Index of the character currently being considered.
	line    int     // Which line `current` is on, one-based.
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]Token, []error) {
	var errors []error
	for !s.atEnd() {
		s.start = s.current
		if err := s.scanToken(); err != nil {
			errors = append(errors, err)
		}
	}
	eofToken := Token{
		tokenType: EOF,
		lexeme:    "",
		literal:   nil,
		line:      s.line,
	}
	s.tokens = append(s.tokens, eofToken)
	return s.tokens, errors
}

func (s *Scanner) atEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	r := s.advance()
	switch r {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.atEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		// Skip whitespace.
	case '\n':
		s.line++
	case '"':
		if err := s.string(); err != nil {
			return err
		}
	default:
		if unicode.IsDigit(r) {
			if err := s.number(); err != nil {
				return err
			}
		} else if unicode.IsLetter(r) || r == '_' {
			s.identifier()
		} else {
			return gloxError(s.line, fmt.Sprintf("Unexpected character: %q.", r))
		}
	}
	return nil
}

func (s *Scanner) advance() rune {
	r, w := utf8.DecodeRuneInString(s.source[s.current:])
	s.current += w
	return r
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addLiteral(tokenType, nil)
}

func (s *Scanner) addLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	token := Token{
		tokenType: tokenType,
		lexeme:    text,
		literal:   literal,
		line:      s.line,
	}
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) match(expected rune) bool {
	if s.atEnd() {
		return false
	}
	r, w := utf8.DecodeRuneInString(s.source[s.current:])
	if r != expected {
		return false
	}
	s.current += w
	return true
}

func (s *Scanner) peek() rune {
	if s.atEnd() {
		return rune(0)
	}
	r, _ := utf8.DecodeRuneInString(s.source[s.current:])
	return r
}

func (s *Scanner) peekNext() rune {
	if s.atEnd() {
		return rune(0)
	}
	_, w := utf8.DecodeRuneInString(s.source[s.current:])
	if s.current+w >= len(s.source) {
		return rune(0)
	}
	r, _ := utf8.DecodeRuneInString(s.source[s.current+w:])
	return r
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.atEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.atEnd() {
		return gloxError(s.line, "Unterminated string literal.")
	}
	s.advance()                                // Read the closing quote.
	value := s.source[s.start+1 : s.current-1] // Trim the surrounding quotes.
	s.addLiteral(STRING, value)
	return nil
}

func (s *Scanner) number() error {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance() // Consume the decimal dot.
		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}
	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return gloxError(s.line, "Invalid number literal.")
	}
	s.addLiteral(NUMBER, num)
	return nil
}

func (s *Scanner) identifier() {
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) || s.peek() == '_' {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}
