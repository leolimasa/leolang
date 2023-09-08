package main

import (
	"fmt"
	"strings"
	"leolang/parser"
)


func main() {
	program := `my-fun = fn (a b)
    print "hello world"
    for 
   	a
	b
    if (a > b) True False`	

	lexer := parser.NewLexer(strings.NewReader(program))
	for {
		t, err := lexer.Next()
		if err != nil {
			panic(err.Error)
		}
		fmt.Printf("%d:%d %d %s \n", t.Line, t.Col, t.Type, t.Value)	
	}
}
