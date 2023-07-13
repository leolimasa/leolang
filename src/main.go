package main

import (
	"bufio"
	"io"
	"unicode/utf8"
)

type TokenType int

const (
	LeftParen TokenType = 1
	RightParen TokenType = 2
	Identifier TokenType = 3
	String TokenType = 4
	Number TokenType = 5
	Indent TokenType = 6
	Dedent TokenType = 7
	Operator TokenType = 8
)

type Token struct {
	Type TokenType
	Line int
	Char int
	Value string
}

type Tokenizer struct {
	scanner bufio.Scanner
	lineNum int
	curLine []rune 
	runePos int
	charPos int
	curRune rune
	indentLevel int
	isSpaceIndent bool
	stop bool
}

type TokenizerError struct {
	Line int
	Char int
	Value string
}

func NewTokenizer(reader io.Reader) Tokenizer {
	return Tokenizer {
		scanner: *bufio.NewScanner(reader),
		lineNum: 0,
		curLine: []rune{},
		curRune: 0,
	}
}

func (t *Tokenizer) nextLine() error {
	isErr := t.scanner.Scan()
	if isErr {
		t.stop = true
		return t.scanner.Err()
	}
	t.curLine = []rune(t.scanner.Text())
	t.lineNum++
	t.runePos = 0
	t.charPos = 0
	return nil
}

func (t *Tokenizer) nextChar() bool {
	// Reached EOL
	if t.runePos >= len(t.curLine) {
		return true 
	}
	t.curRune = t.curLine[t.runePos]
	t.runePos++
	t.charPos = 
	return false
}

func (t *Tokenizer) rewind() {
	if t.runePos == 0 {
		return
	}
	t.runePos--
}

func (t *Tokenizer) newToken(tokentType TokenType, value string) Token {
	return Token {
		Type: tokentType,
		Char: t.runePos,
		Line: t.lineNum,
		Value: value,
	}
}

func (t *Tokenizer) detectIndent() (*Token, error) {
	level := 0
	var tok Token
	for {
		isEol := t.nextChar()
		if t.curRune != ' ' || isEol {
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

func (t *Tokenizer) detectString() (*Token, error) {
	value := ""
	for {
		isEol := t.nextChar()
		if isEol {
			return nil, "Unterminated string"
		}
		if t.curRune == '"' && t.prevRune() != '\\' {

		}
	}
}

// Next returns the next token, or nil for end of stream.
func (t *Tokenizer) Next() (*Token, error) {
	for !t.stop {
		// Read next line
		if t.runePos >= len(t.curLine) {
			err := t.nextLine()
			if err != nil {
				return nil, err
			}
		}

		// Read next character
		t.nextChar()

		// Do not allow tabs as indentation
		// TODO

		// Detect indent at the beginning of the line
		if (t.runePos == 1 && t.curRune == ' ') {
			indentToken, err := t.detectIndent()
			if err != nil {
				return nil, err
			}
			if indentToken != nil {
				return indentToken, nil
			}
		}

		// Detect parens
		if t.curRune == '(' {
			token := t.newToken(LeftParen, "")
			return &token, nil
		}
		if t.curRune == ')' {
			token := t.newToken(RightParen, "")
			return &token, nil
		}

		// Detect string
		if t.curRune == '"' {
			stringToken := t.detectString()
		}

		// Detect identifier (anything that is not a space and is not one of the tokens above). 
		// Convert it to either number or operator, depending on its value.
		if t.curRune != ' ' {
			identifierToken := t.detectIdentifier()
		}

	}
}


func main() {
	
}
