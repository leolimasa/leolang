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
finish-here`

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
	assert.Equal(t, Identifier, token.Type)
	assert.Equal(t, "map", token.Value)

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
	assert.Equal(t, Dedent, token.Type)
}
