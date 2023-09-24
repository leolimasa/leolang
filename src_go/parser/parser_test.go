package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	program := `print "hello world"
map (a (b c) d)
(start-with-parens) c
	`
	
	raw_lexer := NewLexer(strings.NewReader(program))
	lexer := NewLexerNoIndent(raw_lexer)

	parsed, err := Parse(&lexer) 
	assert.Nil(t, err)
	//fmt.Printf("%s\n", parsed.List[0].Atom.Value)
	fmt.Printf("%+v\n", parsed)
	assert.Equal(t, "print", parsed.List[0].List[0].Atom.Value)
	assert.Equal(t, "hello world", parsed.List[0].List[1].Atom.Value)
	assert.Equal(t, "map", parsed.List[1].List[0].Atom.Value)
	assert.Equal(t, "a", parsed.List[1].List[1].List[0].Atom.Value)
	assert.Equal(t, "b", parsed.List[1].List[1].List[1].List[0].Atom.Value)
	assert.Equal(t, "c", parsed.List[1].List[1].List[1].List[1].Atom.Value)
	assert.Equal(t, "d", parsed.List[1].List[1].List[2].Atom.Value)
	assert.Equal(t, "start-with-parens", parsed.List[2].List[0].List[0].Atom.Value)
	assert.Equal(t, "c", parsed.List[2].List[1].Atom.Value)
}
