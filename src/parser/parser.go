package parser

type ParserError struct {
	Line int
	Col int
	Error error
}

type Node struct {
	Value Token
	Children []Node
}


func TokenList(lexer *Lexer) ([]Token, *ParserError) {
	token, err := lexer.Next()
	if err != nil {
		return []Token{}, err
	}
	if token.Type == String || token.Type == Number || token.Type == Identifier {

	}	
}
