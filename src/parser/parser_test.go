package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	program := `my-fun = fn (a b)
	print "hello world"
	map (a (b c) d)
	(start-with-parens) c
	`	
}
