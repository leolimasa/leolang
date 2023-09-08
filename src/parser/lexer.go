package parser 

/**
TODO
- Add DEDENT for every dedent level
- Add LineEnd for line ends
- Add number parsing
- Add string escaping
- Multiline strings
**/

import (
	"bufio"
	"io"
	"regexp"
)

type TokenType int

var numberRegex *regexp.Regexp = regexp.MustCompile("[+-]?([0-9]*[.])?[0-9]+")

var binOperators = []string{
	"+", "-", "*", "/", "**", "%",
	"and", "or", "not",
	"==", "!=", ">", "<", "<=", ">=", 
	"=", ":=",
}

const (
	LeftParen TokenType = 1
	RightParen TokenType = 2
	Identifier TokenType = 3
	String TokenType = 4
	Number TokenType = 5
	Indent TokenType = 6
	Dedent TokenType = 7
	BinOperator TokenType = 8
)

type Token struct {
	Type TokenType
	Line int
	Col int // the column here is NOT the byte position, but rune pos
	Value string
}

type Lexer struct {
	reader bufio.Reader
	line int
	col int
	curRune rune
	runeSize int
	indentLevel int
	numberRegex regexp.Regexp
	skipNextRead bool
}


func NewLexer(reader io.Reader) Lexer {
	return Lexer {
		reader: *bufio.NewReader(reader),
		line: 1,
		col: 0,
		skipNextRead: false,
	}
}


func (t *Lexer) newToken(tokentType TokenType, value string) Token {
	return Token {
		Type: tokentType,
		Col: t.col,
		Line: t.line,
		Value: value,
	}
}

// detectIndent returns a indent or dedent token.
func (l *Lexer) detectIndent() (*Token, *ParserError) {
	level := 0
	var tok Token
	for {
		err := l.readRune()
		if err != nil {
			return nil, err
		}

		// Emtpy line. Ignore.
		if l.curRune == '\n' {
			return nil, nil
		}

		// Check if we reached the end of the indentation
		// TODO emit DEDENTS according to how many levels were
		// dedented
		// can probably use a indent stack
		if l.curRune != ' ' {
			l.skipNextRead = true
			if level > l.indentLevel {
				tok = l.newToken(Indent, "")
				l.indentLevel = level
				return &tok, nil
			}
			if level < l.indentLevel {
				tok = l.newToken(Dedent, "")
				l.indentLevel = level
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
		if l.curRune == ' ' || l.curRune == '\n' || l.curRune == '\r' || l.curRune == ')' {

			// We consumed the end character for the identifier, so make sure
			// it gets processed on the following Next() call
			l.skipNextRead = true

			// Remove last character from value
			value = value[:len(value)-1]

			// Return a binary op token if the identifier is a binary operator
			for _, op := range binOperators {
				if op == value {
					tok := l.newToken(BinOperator, value)
					return &tok, nil
				}
			}
			// TODO whether the identifier is a number
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
	return ParserError {
		Line: l.line,
		Col: l.col,
		Error: err,
	}
}

func (l *Lexer) readRune() *ParserError {
	if l.skipNextRead {
		l.skipNextRead = false
		return nil
	}
	curRune, _, err := l.reader.ReadRune()
	l.curRune = curRune
	l.col++
	if err != nil {
		tokError := l.newError(err)
		return &tokError
	}
	return nil
}

// Next returns the next token, or nil for end of stream.
func (l *Lexer) Next() (*Token, *ParserError) {
	for {
		err := l.readRune()

		// Check if we reached end of file
		if err != nil {
			if err.Error == io.EOF {
				return nil, nil
			}
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
			continue
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
