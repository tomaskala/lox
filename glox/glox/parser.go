package glox

type Parser struct {
	tokens  []Token // List of tokens to be parsed.
	current int     // Index of the next token to be parsed.
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.expression()
	return expr, err
}

func (p *Parser) expression() (Expr, error) {
	expr, err := p.equality()
	return expr, err
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return expr, err
	}
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return expr, err
		}
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return expr, err
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return expr, err
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
		return expr, err
	}
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return expr, err
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
		return expr, err
	}
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return expr, err
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
		return Unary{
			operator: operator,
			right:    right,
		}, err
	} else {
		return p.primary()
	}
}

func (p *Parser) primary() (Expr, error) {
	switch {
	case p.match(FALSE):
		return Literal{value: false}, nil
	case p.match(TRUE):
		return Literal{value: true}, nil
	case p.match(NIL):
		return Literal{value: nil}, nil
	case p.match(NUMBER, STRING):
		return Literal{p.previous().literal}, nil
	case p.match(LEFT_PAREN):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after an expression.")
		if err != nil {
			return nil, err
		}
		return Grouping{expression: expr}, nil
	default:
		return nil, parseError(p.peek(), "Expect an expression.")
	}
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
	if p.atEnd() {
		return false
	} else {
		return p.peek().tokenType == tokenType
	}
}

func (p *Parser) advance() Token {
	if !p.atEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) atEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	} else {
		dummy := Token{
			tokenType: EOF,
			lexeme: "",
			literal: nil,
			line: 0,
		}
		return dummy, parseError(p.peek(), message)
	}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.atEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}
		switch p.peek().tokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		default:
			p.advance()
		}
	}
}
