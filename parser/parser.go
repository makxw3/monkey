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
	p.nextToken()
	// Set the Name variable to a new Identifier using the currentToken
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	// Check if the next token is token.ASSIGN
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// Move to the next token
	p.nextToken()
	// TODO: Skip all the expressions for now
	// Skip all expressions until you get to the SEMICOLON token
	for !p.currentTokenIs(token.SEMICOLON) {
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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParsingFunctions[p.currentToken.Type]
	if prefixFn == nil {
		errMsg := fmt.Sprintf("no prefixParseFunction found for %s", p.currentToken.Type)
		p.errors = append(p.errors, errMsg)
		return nil
	}
	leftExpression := prefixFn()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixFn := p.infixParsingFunctions[p.peekToken.Type]
		if infixFn == nil {
			return leftExpression
		}
		// Advances to the next token
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

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParsingFunction) {
	p.prefixParsingFunctions[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParsingFunction) {
	p.infixParsingFunctions[tt] = fn
}

// The two types of functions needed for a pratt parser ie. a Prefix parsing function and an Infix parsing function
type (
	prefixParsingFunction func() ast.Expression               // Returns an ast.Expression object
	infixParsingFunction  func(ast.Expression) ast.Expression // Takes in an ast.Expression object and returns an ast.Expression object
)
