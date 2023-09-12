package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	program := `my-fun = fn (a b)
    print "hello world"

    map
        a b
    some-call
    line-ends-here
dedented-all-the-way
    indent-one-level
      indent-two-levels
dedent-again`

	lexer := NewLexer(strings.NewReader(program))

	token, _ := lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "my-fun", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, BinOperator, token.Type)
	assert.Equal(t, "=", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, token.Value, "fn", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, token.Type, LeftParen,  token.Type)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "a", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "b",token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, RightParen, token.Type)

	token, _ = lexer.Next()
	assert.Equal(t, Indent, token.Type)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "print", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, String, token.Type)
	assert.Equal(t, "hello world", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, LineEnd, token.Type)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "map", token.Value)
	assert.Equal(t, 4, token.Line)

	token, _ = lexer.Next()
	assert.Equal(t, Indent, token.Type)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "a", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "b", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, Dedent, token.Type)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "some-call", token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, LineEnd, token.Type)

	lexer.Next() // line-ends-here

	token, _ = lexer.Next() 
	assert.Equal(t, Dedent, token.Type)
	assert.Equal(t, 1, token.Value)

	lexer.Next() // dedent-all-the-way
	token, _ = lexer.Next()
	assert.Equal(t, Indent, token.Type)

	lexer.Next() // indent-one-level
	token, _ = lexer.Next()
	assert.Equal(t, Indent, token.Type)

	token, _ = lexer.Next() // indent-two-levels
	token, _ = lexer.Next()
	assert.Equal(t, Dedent, token.Type)
	assert.Equal(t, 2, token.Value)

	token, _ = lexer.Next()
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "dedent-again", token.Value)

	token, _ = lexer.Next()
	assert.Nil(t, token)
}
