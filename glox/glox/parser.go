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

func (p *Parser) Parse() ([]Stmt, []error) {
	var statements []Stmt
	var errors []error
	for !p.atEnd() {
		stmt, err := p.declaration()
		if err != nil {
			errors = append(errors, err)
		} else {
			statements = append(statements, stmt)
		}
	}
	return statements, errors
}

func (p *Parser) declaration() (stmt Stmt, err error) {
	defer func() {
		if r := recover(); r != nil {
			if pr, ok := r.(parserError); ok {
				p.synchronize()
				err = pr
			} else {
				panic(r)
			}
		}
	}()
	if p.match(VAR) {
		return p.varDeclaration(), nil
	} else {
		return p.statement(), nil
	}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect a variable name.")
	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after a variable declaration.")
	return Var{
		name:        name,
		initializer: initializer,
	}
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	} else {
		return p.expressionStatement()
	}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after a value.")
	return Print{
		expression: value,
	}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after an expression.")
	return Expression{
		expression: expr,
	}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.equality()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if v, ok := expr.(Variable); ok {
			return Assign{
				name:  v.name,
				value: value,
			}
		} else {
			panic(parserError{parseError(equals, "Invalid assignment target.")})
		}
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (p *Parser) unary() Expr {
	switch {
	case p.match(BANG, MINUS):
		operator := p.previous()
		right := p.unary()
		return Unary{
			operator: operator,
			right:    right,
		}
	case p.match(BANG_EQUAL, EQUAL_EQUAL):
		p.comparison()
		panic(parserError{parseError(p.previous(), "Missing left operand.")})
	case p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL):
		p.term()
		panic(parserError{parseError(p.previous(), "Missing left operand.")})
	case p.match(PLUS):
		p.factor()
		panic(parserError{parseError(p.previous(), "Missing left operand.")})
	case p.match(SLASH, STAR):
		p.unary()
		panic(parserError{parseError(p.previous(), "Missing left operand.")})
	default:
		return p.primary()
	}
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return Literal{value: false}
	case p.match(TRUE):
		return Literal{value: true}
	case p.match(NIL):
		return Literal{value: nil}
	case p.match(NUMBER, STRING):
		return Literal{value: p.previous().literal}
	case p.match(IDENTIFIER):
		return Variable{name: p.previous()}
	case p.match(LEFT_PAREN):
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after an expression.")
		return Grouping{expression: expr}
	default:
		panic(parserError{parseError(p.peek(), "Expect an expression.")})
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

func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	} else {
		panic(parserError{parseError(p.peek(), message)})
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
