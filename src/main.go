package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type TokenType int

var numberRegex *regexp.Regexp = regexp.MustCompile("[+-]?([0-9]*[.])?[0-9]+")
var binOperators = []string{
	"+", "-", "*", "/", "**", "%",
	"and", "or",
	"==", "!=", ">", "<", "<=", ">=",
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
	Col int
	Value string
}

type Lexer struct {
	reader bufio.Reader
	line int
	col int
	runeSize int
	indentLevel int
	numberRegex regexp.Regexp
}

type LexerError struct {
	Line int
	Col int
	Error error
}

func NewLexer(reader io.Reader) Lexer {
	return Lexer {
		reader: *bufio.NewReader(reader),
		line: 1,
		col: 0,
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

func (t *Lexer) detectIndent() (*Token, *LexerError) {
	level := 0
	var tok Token
	for {
		curRune, err := t.readRune()
		if err != nil {
			return nil, err
		}

		if curRune != ' ' {
			t.rewind()
			if level > t.indentLevel {
				tok = t.newToken(Indent, "")
				t.indentLevel = level
				return &tok, nil
			}
			if level < t.indentLevel {
				tok = t.newToken(Dedent, "")
				t.indentLevel = level
				return &tok, nil
			}
			return nil, nil
		}
		level++
	}
}

func (t *Lexer) detectString() (*Token, *LexerError) {
	value := ""
	for {
		curRune, err := t.readRune()
		if err != nil {
			return nil, err
		}
		// TODO parse escape characters
		if curRune == '"' {
			tok := t.newToken(String, value)
			return &tok, nil
		}
		value += string(curRune)
	}
}

func (t *Lexer) detectIdentifier() (*Token, *LexerError) {
	value := ""
	for {
		curRune, err := t.readRune()
		if err != nil {
			return nil, err
		}
		if curRune == ' ' || curRune == '\n' || curRune == '\r' || curRune == ')' {
			t.rewind()
			for _, op := range binOperators {
				if op == value {
					tok := t.newToken(BinOperator, value)
					return &tok, nil
				}
			}
			tok := t.newToken(Identifier, value)
			return &tok, nil
				
		}
		value += string(curRune)
	}
}

func (t *Lexer) newError(err error) LexerError {
	return LexerError {
		Line: t.line,
		Col: t.col,
		Error: err,
	}
}

func (t *Lexer) readRune() (rune, *LexerError) {
	curRune, size, err := t.reader.ReadRune()
	t.col += size
	if err != nil {
		tokError := t.newError(err)
		return '0', &tokError
	}
	return curRune, nil
}

func (t *Lexer) rewind() *LexerError {
	err := t.reader.UnreadRune()
	if err != nil {
		tokError := t.newError(err)
		return &tokError
	}
	t.col -= t.runeSize
	return nil
}

// Next returns the next token, or nil for end of stream.
func (t *Lexer) Next() (*Token, *LexerError) {
	for {
		curRune, err := t.readRune()
		if err != nil {
			if err.Error == io.EOF {
				return nil, nil
			}
			return nil, err
		}

		// Do not allow tabs as indentation
		// TODO

		// Detect indent at the beginning of the line
		// TODO indent is broken
		if (t.col == 1 && curRune == ' ') {
			indentToken, err := t.detectIndent()
			if err != nil {
				return nil, err
			}
			if indentToken != nil {
				return indentToken, nil
			}
		}

		// Detect parens
		if curRune == '(' {
			token := t.newToken(LeftParen, "")
			return &token, nil
		}
		if curRune == ')' {
			token := t.newToken(RightParen, "")
			return &token, nil
		}

		// Detect string
		if curRune == '"' {
			stringToken, err := t.detectString()
			if err != nil {
				return nil, err
			}
			return stringToken, nil
		}

		// Detect identifier (anything that is not a space and is not one of the tokens above). 
		// Convert it to either number or operator, depending on its value.
		if curRune != ' ' {
			identifierToken, err := t.detectIdentifier()
			if err != nil {
				return nil, err
			}
			// TODO this can probably be handled with a rewind
			identifierToken.Value = string(curRune) + identifierToken.Value
			return identifierToken, nil
		}
	}
}


func main() {
	program := `my-fun = fn (a b)
	        print "hello world"
	        if (a > b) True False`	

	lexer := NewLexer(strings.NewReader(program))
	for {
		t, err := lexer.Next()
		if err != nil {
			panic(err.Error)
		}
		fmt.Printf("%d:%d %d %s \n", t.Line, t.Col, t.Type, t.Value)	
	}
}
