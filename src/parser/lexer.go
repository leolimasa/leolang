package parser

/**
TODO
- Add hex and more complex primitives parsing
- Add string escaping
- Multiline strings
**/

import (
	"bufio"
	"io"
	"regexp"
)

type TokenType int

var floatRegex *regexp.Regexp = regexp.MustCompile("[+-]?([0-9]*[.])?[0-9]+")

var integerRegex *regexp.Regexp = regexp.MustCompile(`^-?\d+$`)

var operators = []string{
	"+", "-", "*", "/", "**", "%",
	"and", "or", "not",
	"==", "!=", ">", "<", "<=", ">=",
	"=", "v=",
}

const (
	LeftParen  TokenType = 1
	RightParen TokenType = 2
	Identifier TokenType = 3
	String     TokenType = 4
	Int        TokenType = 5
	Float      TokenType = 6
	Indent     TokenType = 7
	Dedent     TokenType = 8
	Operator   TokenType = 9
	LineEnd    TokenType = 10
)

type Token struct {
	Type  TokenType
	Line  int
	Col   int // the column here is NOT the byte position, but rune pos
	Value any
}

type Lexer struct {
	reader       bufio.Reader
	line         int
	col          int
	curRune      rune
	runeSize     int
	skipNextRead bool
	indentStack  []int
	isEof        bool
}

func NewLexer(reader io.Reader) Lexer {
	return Lexer{
		reader:       *bufio.NewReader(reader),
		line:         1,
		col:          0,
		skipNextRead: false,
		isEof:        false,
	}
}

func (l *Lexer) newToken(tokentType TokenType, value any) Token {
	return Token{
		Type:  tokentType,
		Col:   l.col,
		Line:  l.line,
		Value: value,
	}
}

func (l *Lexer) getIndentLevel() int {
	level := 0
	for _, curLevel := range l.indentStack {
		level += curLevel
	}
	return level
}

// detectIndent returns a indent or dedent token.
// will return the same amount of dedents as there are levels dedented.
func (l *Lexer) detectIndent() (*Token, *ParserError) {
	level := 0
	for {
		err := l.readRune()
		if err != nil {
			return nil, err
		}

		// Emtpy line. Ignore.
		if l.curRune == '\n' {
			l.line++
			return nil, nil
		}

		// Check if we reached the end of the indentation
		if l.curRune != ' ' || l.isEof {
			l.skipNextRead = true
			indentLevel := l.getIndentLevel()
			if level > indentLevel {
				tok := l.newToken(Indent, 1)
				l.indentStack = append(l.indentStack, level-indentLevel)
				return &tok, nil
			}
			if level < indentLevel {
				// calculates dedent levels
				dedents := 0
				for l.getIndentLevel() > level {
					dedents++
					l.indentStack = l.indentStack[:len(l.indentStack)-1]
				}

				tok := l.newToken(Dedent, dedents)
				return &tok, nil
			}
			return nil, nil
		}
		level++
	}
}

func (l *Lexer) detectString() (*Token, *ParserError) {
	value := ""
	for {
		err := l.readRune()
		if err != nil {
			return nil, err
		}
		// TODO parse escape characters
		if l.curRune == '"' {
			tok := l.newToken(String, value)
			return &tok, nil
		}
		value += string(l.curRune)
	}
}

func (l *Lexer) detectIdentifier() (*Token, *ParserError) {
	value := string(l.curRune)
	for {
		// Check if we're done reading the identifier
		if l.curRune == ' ' || l.curRune == '\n' || l.curRune == '\r' || l.curRune == ')' || l.isEof {

			// We consumed the end character for the identifier, so make sure
			// it gets processed on the following Next() call
			l.skipNextRead = true

			// Remove last character from value
			value = value[:len(value)-1]

			// Return a binary op token if the identifier is a binary operator
			for _, op := range operators {
				if op == value {
					tok := l.newToken(Operator, value)
					return &tok, nil
				}
			}

			// Check if it's an int
			if integerRegex.MatchString(value) {
				tok := l.newToken(Int, value)
				return &tok, nil
			}

			// Check if it's a float
			if floatRegex.MatchString(value) {
				tok := l.newToken(Float, value)
				return &tok, nil
			}

			tok := l.newToken(Identifier, value)
			return &tok, nil

		}
		// continue reading
		err := l.readRune()
		if err != nil {
			return nil, err
		}
		value += string(l.curRune)
	}
}

func (l *Lexer) newError(err error) ParserError {
	return ParserError{
		Line:  l.line,
		Col:   l.col,
		Error: err,
	}
}

func (l *Lexer) readRune() *ParserError {
	if l.isEof {
		return nil
	}
	
	if l.skipNextRead {
		l.skipNextRead = false
		return nil
	}
	curRune, _, err := l.reader.ReadRune()
	l.curRune = curRune
	l.col++
	if err != nil {
		if err == io.EOF {
			l.isEof = true
		} else {
			tokError := l.newError(err)
			return &tokError
		}
	}
	return nil
}

// Next returns the next token, or nil for end of stream.
func (l *Lexer) Next() (*Token, *ParserError) {
	for {
		err := l.readRune()

		if l.isEof {
			return nil, nil
		}

		// Check if we reached end of file
		if err != nil {
			return nil, err
		}

		// New line
		if l.curRune == '\n' {
			l.col = 0
			l.line++
			indentToken, err := l.detectIndent()
			if err != nil {
				return nil, err
			}
			if indentToken != nil {
				return indentToken, nil
			}
			token := l.newToken(LineEnd, nil)
			return &token, nil
			//continue
		}

		// Detect parens
		if l.curRune == '(' {
			token := l.newToken(LeftParen, "")
			return &token, nil
		}
		if l.curRune == ')' {
			token := l.newToken(RightParen, "")
			return &token, nil
		}

		// Detect string
		if l.curRune == '"' {
			stringToken, err := l.detectString()
			if err != nil {
				return nil, err
			}
			return stringToken, nil
		}

		// Detect identifier (anything that is not a space and is not one of the tokens above).
		// Convert it to either number or operator, depending on its value.
		if l.curRune != ' ' && l.curRune != '\t' && l.curRune != '\n' {
			identifierToken, err := l.detectIdentifier()
			if err != nil {
				return nil, err
			} 
			return identifierToken, nil
		}
	}
}

type LexerNoIndent struct {
	lexer Lexer
	tokenBuffer []*Token
	nextIsOpenParen bool
}

func NewLexerNoIndent(lexer Lexer) LexerNoIndent {
	return LexerNoIndent {
		lexer: lexer,
		nextIsOpenParen: false,
	}
}

func (ln *LexerNoIndent) Next() (*Token, *ParserError) {
	if len(ln.tokenBuffer) > 0 {
		tok := ln.tokenBuffer[len(ln.tokenBuffer)-1]
		ln.tokenBuffer = ln.tokenBuffer[:len(ln.tokenBuffer)-1]
		return tok, nil
	}
	token, err := ln.lexer.Next()
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, nil
	}
	if ln.nextIsOpenParen {
		tok := ln.lexer.newToken(LeftParen, nil)
		ln.nextIsOpenParen = false
		ln.tokenBuffer = append(ln.tokenBuffer, token)
		return &tok, nil
	}
	if token.Type == LineEnd {
		tok := ln.lexer.newToken(RightParen, nil)
		ln.nextIsOpenParen = true
		return &tok, nil
	}

	if token.Type == Indent {
		ln.nextIsOpenParen = true
		return ln.Next()
	}
	if token.Type == Dedent {
		// we remove one dedent because the first one will
		// be emitted right away.
		dedents := token.Value.(int) - 1
		for i := 0; i < dedents; i++ {
			tok := ln.lexer.newToken(RightParen, nil)
			ln.tokenBuffer = append(ln.tokenBuffer, &tok)
		}
		tok := ln.lexer.newToken(RightParen, nil)
		return &tok, nil
	}
	return token, nil
}


