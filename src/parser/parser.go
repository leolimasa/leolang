package parser

type ParserError struct {
	Line int
	Col int
	Error error
}

// Node stores a list of token or a single token.
// This would ideally be a tagged union but go doesn't have that...
type Node struct {
	Value *Token
	List []Node
}

// func TokenList(lexer *Lexer) ([]Token, *ParserError) {
// 	token, err := lexer.Next()
// 	if err != nil {
// 		return []Token{}, err
// 	}
// 	if token.Type == String || token.Type == Number || token.Type == Identifier {
//
// 	}	
// }
