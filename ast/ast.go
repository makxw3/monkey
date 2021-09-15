package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

// Every Node in the ast implements the Node interface
// The TokenLiteral returns the LiteralValue of the associated token
// TokenLiteral() is used for debugging
type Node interface {
	TokenLiteral() string
	String() string
}

// Statements Nodes do not produce any values
// statementNode function is just for debugging
type Statement interface {
	Node
	statementNode()
}

// Expression Nodes produces values
// expressionNode function is just for debugging
type Expression interface {
	Node
	expressionNode()
}

// Every go program is a series of Statements
// The Program struct implements the Node interface
// The Program Node is the root of the AST
type Program struct {
	Statements []Statement
}

// Returns the LiteralValue of the token of the first Statement
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

// The LetStatement contains the token.LET token and the Identifier and the Value
type LetStatement struct {
	Token token.Token // The token.LET token
	Name  *Identifier // The identidier name
	Value Expression  // The expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.Value)
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// The identifier contains the token.IDENT token and the name value
type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string      // The name
}

func (id *Identifier) TokenLiteral() string {
	return id.Token.Literal
}

func (id *Identifier) expressionNode() {}

func (id *Identifier) String() string {
	return id.Value
}

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression  // the return value
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	out.WriteString(rs.ReturnValue.String())
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // The first token in the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token // The prefix token eg !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token eg. +
	Left     Expression  // The expression on the right of the operator
	Operator string      // The operator
	Right    Expression  // The expressin on the right of the operator
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// Boolean Statements
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {

}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.TokenLiteral()
}

// IF Expressions
type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, pr := range fl.Parameters {
		params = append(params, pr.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}
