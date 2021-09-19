package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type Parser struct {
	l                      *lexer.Lexer
	currentToken           token.Token
	peekToken              token.Token
	errors                 []string
	prefixParsingFunctions map[token.TokenType]prefixParsingFunction
	infixParsingFunctions  map[token.TokenType]infixParsingFunction
}

// New is a helper function to create a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Read the next 2 tokens
	p.nextToken()
	p.nextToken()

	p.prefixParsingFunctions = make(map[token.TokenType]prefixParsingFunction)
	p.infixParsingFunctions = make(map[token.TokenType]infixParsingFunction)
	// Register the prefix-parsing-functions
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	// Register the infix-parsing-function which is the same for all infix operators
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	// Creating an empty program object with 0 Statements
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Go through all the tokens until the token.EOF token
	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			// Add the parsed statement to the []Statement
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// Helper function to check if the next token is of type tt
func (p *Parser) expectPeek(tt token.TokenType) bool {
	match := p.peekToken.Type == tt
	if !match {
		msg := fmt.Sprintf("Expected next token to be %s. Got %s instead", tt, p.peekToken.Type)
		p.errors = append(p.errors, msg)
	} else {
		p.nextToken()
	}
	return match
}

// Helper function to check if the current token is of type tt
func (p *Parser) currentTokenIs(tt token.TokenType) bool {
	return p.currentToken.Type == tt
}

// Helper function to check if the peek token is of type tt
func (p *Parser) peekTokenIs(tt token.TokenType) bool {
	return p.peekToken.Type == tt
}

// parseStatement parses Statements
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		// parse the LetStatement if the current token is token.LET
		return p.parseLetStatement()
	case token.RETURN:
		// parse the ReturnStatement if the current token is token.RETURN
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses LET statements
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}
	// Check if the next token is token.IDENT
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// Move to the next token if the next token is token.IDENT
	// p.nextToken()
	// Set the Name variable to a new Identifier using the currentToken
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	// Check if the next token is token.ASSIGN
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// Move to the next token
	// p.nextToken()
	// TODO: Skip all the expressions for now
	// Skip all expressions until you get to the SEMICOLON token
	for p.currentToken.Type != token.SEMICOLON && !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	// move to the next token. past the token.RETURN token
	p.nextToken()
	// TODO:Skip all the expressions until you get to the token.RETURN token
	if !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// Constants used to compare the precedence of different operators
const (
	_ int = iota
	LOWEST
	EQUALS      // == !=
	LESSGREATER // < >
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // ! + -
	CALL        // fn()
)

// The precedence table
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// peekPrecedence is a helper function that returns the precedence of the next token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// currentPrecedence is a helper function that returns the precedence of the current token
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// TODO: I don't understand how this part works
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Parse the next expression using the prefix-parsing-function eg returns integer or identifier
	prefixFn := p.prefixParsingFunctions[p.currentToken.Type]
	leftExpression := prefixFn()
	for precedence < p.peekPrecedence() {
		// returns an ast.InfixExpression for 1 + 2
		infixFn := p.infixParsingFunctions[p.peekToken.Type]
		p.nextToken()
		leftExpression = infixFn(leftExpression)
	}
	return leftExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}
	val, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %s as an integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	// Get the precedence of the current token
	precedence := p.currentPrecedence()
	// Advance to the next token
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	// parse the expression
	expr.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// p.nextToken()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	// p.nextToken()

	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		// p.nextToken()

		expr.Alternative = p.parseBlockStatement()
	}
	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()
	for !p.currentTokenIs(token.LBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParsingFunction) {
	p.prefixParsingFunctions[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParsingFunction) {
	p.infixParsingFunctions[tt] = fn
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// check if there are not parameters
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()

	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.peekTokenIs(token.SEMICOLON) {
		return nil
	}
	return identifiers
}

// The two types of functions needed for a pratt parser ie. a Prefix parsing function and an Infix parsing function
type (
	prefixParsingFunction func() ast.Expression               // Returns an ast.Expression object
	infixParsingFunction  func(ast.Expression) ast.Expression // Takes in an ast.Expression object and returns an ast.Expression object
)
