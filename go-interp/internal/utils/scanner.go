package utils

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	Source  string
	Tokens  []Token
	Start   int
	Current int
	Line    int
}

func (s *Scanner) ScanTokens() []Token {

	fmt.Println("source", s.Source)
	for true {
		if s.Current >= len(s.Source) {
			break
		}
		s.Start = s.Current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, Token{
		TokenType: EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.Line,
	})

	return s.Tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case "(":
		s.addToken(LEFT_PAREN)
		break
	case ")":
		s.addToken(RIGHT_PAREN)
		break
	case "{":
		s.addToken(LEFT_BRACE)
		break
	case "}":
		s.addToken(RIGHT_BRACE)
		break
	case ",":
		s.addToken(COMMA)
		break
	case ".":
		s.addToken(DOT)
		break
	case "-":
		s.addToken(MINUS)
		break
	case "+":
		s.addToken(PLUS)
		break
	case ";":
		s.addToken(SEMICOLON)
		break
	case "*":
		s.addToken(STAR)
		break
	case "!":
		if s.match("=") {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
		break
	case "=":
		if s.match("=") {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
		break
	case "<":
		if s.match("=") {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
		break
	case ">":
		if s.match("=") {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
		break
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && s.Current < len(s.Source) {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
		break

	case " ":
	case "\r":
	case "\t":
		// Ignore whitespace.
		break

	case "\n":
		s.Line = s.Line + 1
		break
	case "\"":
		s.buildString()
		break
	default:
		if isDigit(c) {
			s.buildNumber()
		} else if isAlpha(c) {
			s.buildIdentifier()
		} else {
			ErrorOut(s.Current, fmt.Sprintf("Unexpected character: %s", c))
			return
		}
	}
}

func (s *Scanner) advance() string {
	char := string(s.Source[s.Current])
	s.Current = s.Current + 1
	return char
}

func (s *Scanner) peek() string {
	if s.Current >= len(s.Source) {
		return ""
	}

	return string(s.Source[s.Current])
}

func (s *Scanner) peekNext() string {
	if s.Current+1 >= len(s.Source) {
		return ""
	}
	return string(s.Source[s.Current+1])
}

func (s *Scanner) match(expected string) bool {
	if s.Current >= len(s.Source) {
		return false
	}

	if string(s.Source[s.Current]) != expected {
		return false
	}

	s.Current = s.Current + 1
	return true

}

func (s *Scanner) buildString() {
	for s.peek() != "\"" && s.Current < len(s.Source) {
		if s.peek() == "\n" {
			s.Line = s.Line + 1
		}
		s.advance()
	}

	if s.Current >= len(s.Source) {
		ErrorOut(s.Line, "Unterminated string 💀")
		return
	}

	s.advance()
	value := s.Source[s.Start+1 : s.Current-1]
	s.addTokenLiteral(STRING, value)

}

func (s *Scanner) buildNumber() {
	for isDigit(s.peek()) {
		s.advance()

	}
	if s.peek() == "." && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}
	float, err := strconv.ParseFloat(string(s.Source[s.Start:s.Current]), 64)
	if err != nil {
		ErrorOut(s.Current, "Error parsing number literal")
	}
	s.addTokenLiteral(NUMBER, float)
}

func (s *Scanner) buildIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.Source[s.Start:s.Current]
	tokenType, ok := keywordMap[text]

	if ok {
		s.addToken(tokenType)
	} else {
		s.addToken(IDENTIFIER)
	}

}

func isDigit(c string) bool {
	return c >= "0" && c <= "9"
}

func isAlpha(c string) bool {
	return (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || c == "_"
}

func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType TokenType, literal any) {
	text := s.Source[s.Start:s.Current]
	s.Tokens = append(s.Tokens, Token{
		TokenType: tokenType,
		Lexeme:    text,
		Literal:   literal,
		Line:      s.Line,
	})
}
