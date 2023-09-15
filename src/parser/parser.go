package parser

type ParserError struct {
	Line int
	Col int
	Error error
}

type Parser struct {
}

// Node stores a list of token or a single token.
// This would ideally be a tagged union but go doesn't have that...
type Node struct {
	Atom *Token
	List []Node
}


func (p *Parser) Parse(lexer * o)
