package glox

import "fmt"

const MAX_ARGS = 255

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
	} else if p.match(FUN) {
		return p.funDeclaration("function"), nil
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

func (p *Parser) funDeclaration(kind string) Stmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect a %s name.", kind))
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	parameters := make([]Token, 0)
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= MAX_ARGS {
				panic(parserError{parseError(p.peek(), fmt.Sprintf("At most %d parameters to a %s are supported.", MAX_ARGS, kind))})
			}
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect a parameter name."))
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, fmt.Sprintf("Expect ')' after %s parameters.", kind))
	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := p.block()
	return Function{
		name:   name,
		params: parameters,
		body:   body,
	}
}

func (p *Parser) statement() Stmt {
	switch {
	case p.match(PRINT):
		return p.printStatement()
	case p.match(LEFT_BRACE):
		return Block{
			statements: p.block(),
		}
	case p.match(IF):
		return p.ifStatement()
	case p.match(WHILE):
		return p.whileStatement()
	case p.match(FOR):
		return p.forStatement()
	case p.match(RETURN):
		return p.returnStatement()
	default:
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

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)
	for !p.check(RIGHT_BRACE) && !p.atEnd() {
		stmt, err := p.declaration()
		if err != nil {
			panic(parserError{err})
		}
		statements = append(statements, stmt)
	}
	p.consume(RIGHT_BRACE, "Expect '}' after a block.")
	return statements
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after an if condition.")
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	} else {
		elseBranch = nil
	}
	return If{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after a while statement.")
	body := p.statement()
	return While{
		condition: condition,
		body:      body,
	}
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}
	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	} else {
		condition = Literal{
			value: true,
		}
	}
	p.consume(SEMICOLON, "Expect ';' after a loop condition.")
	var increment Expr
	if !p.check(SEMICOLON) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after a for statement.")
	body := p.statement()
	if increment != nil {
		body = Block{
			statements: []Stmt{
				body,
				Expression{
					expression: increment,
				},
			},
		}
	}
	body = While{
		condition: condition,
		body:      body,
	}
	if initializer != nil {
		body = Block{
			statements: []Stmt{
				initializer,
				body,
			},
		}
	}
	return body
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	} else {
		value = nil
	}
	p.consume(SEMICOLON, "Expect ';' after return.")
	return Return{
		keyword: keyword,
		value:   value,
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
	expr := p.or()
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

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = Logical{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = Logical{
			left:     expr,
			operator: operator,
			right:    right,
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
		return p.call()
	}
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
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

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := make([]Expr, 0)
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= MAX_ARGS {
				panic(parserError{parseError(p.peek(), fmt.Sprintf("At most %d arguments to a function are supported.", MAX_ARGS))})
			}
			arguments = append(arguments, p.expression())
			if !p.match(COMMA) {
				break
			}
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after a function call.")
	return Call{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}
}
