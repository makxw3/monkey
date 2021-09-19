func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParsingFunctions[p.currentToken.Type]
	leftExpression := prefixFn()
	for precedence < p.peekPrecedence() {
		infixFn := p.infixParsingFunctions[p.peekToken.Type]
		p.nextToken()
		leftExpression = infixFn(leftExpression)
	}
	return leftExpression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

