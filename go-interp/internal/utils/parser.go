package utils

import "fmt"

type Parser struct {
	Tokens  []Token
	Current int
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}
		switch p.peek().TokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}

	}

	return expr, nil
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current = p.Current + 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[(p.Current - 1)]
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{
			operator: operator,
			right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {

	if p.match(FALSE) {
		return Literal{
			value: nil,
		}, nil
	}

	if p.match(TRUE) {
		return Literal{
			value: TRUE,
		}, nil
	}

	if p.match(NIL) {
		return Literal{
			value: NIL,
		}, nil
	}

	if p.match(NUMBER, STRING) {
		return Literal{
			value: p.previous().Literal,
		}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return Grouping{
			expression: expr,
		}, nil
	}

	return nil, p.parseError(p.peek(), "Expected expression")
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return Token{}, p.parseError(p.peek(), message)
}

type ParseError struct {
	errorText string
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("%s", pe.errorText)
}

func (p *Parser) parseError(token Token, message string) ParseError {
	errorText := p.error(token, message)
	return ParseError{errorText: errorText}
}

func (p *Parser) error(token Token, message string) string {
	if token.TokenType == EOF {
		return Report(token.Line, " at end", message)
	} else {
		return Report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}
