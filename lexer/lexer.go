package lexer

import "monkey/token"

type Lexer struct {
	input           string // The input string
	currentPosition int    // The current position of the input
	readPosition    int    // currentPosition + 1
	ch              byte   // The current byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.currentPosition = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()
	switch l.ch {
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			tok.Literal = "!="
			tok.Type = token.NOT_EQ
			return tok
		}
		tok = newToken(token.BANG, l.ch)
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			tok.Literal = "=="
			tok.Type = token.EQ
			return tok
		}
		tok = newToken(token.ASSIGN, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookUpIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}

// readIdentifier is a helper function that reads all the characters until it encounters a character that is not a letter
func (l *Lexer) readIdentifier() string {
	position := l.currentPosition
	// Continue reading until you encounter a non-letter character
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.currentPosition]
}

// skipWhiteSpace is a helper function to remove all the whiteSpaces
func (l *Lexer) skipWhiteSpace() {
	for l.ch == '\t' || l.ch == '\n' || l.ch == ' ' || l.ch == '\r' {
		l.readChar()
	}
}

// isLetter is a helper function that checks if the byte is a letter
func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

// readNumber is a helper function that reads all the bytes until it encounters a non number
func (l *Lexer) readNumber() string {
	position := l.currentPosition
	// read all the digits
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.currentPosition]
}

// isDigit is a helper function to check if the byte is a digit
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// newToken is a helper function to return a  new token.Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// peekChar is a helper function that returns the next char
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
