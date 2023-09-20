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
type SymbolExpr struct {
	Atom *Token
	List []SymbolExpr
}

func Parse(lexer ILexer) (*SymbolExpr, *ParserError) {
	curExpr := SymbolExpr {}
	for {
		token, err := lexer.Next()
		if err != nil {
			return nil, err
		}

		// EOF
		if token == nil {
			return &curExpr, nil
		}
	
		// Start of a new sub expression
		if token.Type == LeftParen {
			subExpr, err := Parse(lexer)
			if err != nil {
				return nil, err
			}
			curExpr.List = append(curExpr.List, *subExpr)
		}

		// End of expression
		if token.Type == RightParen {
			return &curExpr, nil
		}

		// Atom
		if token.Type == Identifier || token.Type == String || token.Type == Int || token.Type == Float || token.Type == Operator  {
			curExpr.List = append(curExpr.List, SymbolExpr {Atom: token})
		}
	}
}

// func NormalizeIfs(expr SymbolExpr) SymbolExpr {
// }
//
// func GroupOperators(expr SymbolExpr) SymbolExpr {
// }
//
// func DesugarWrapOperator(expr SymbolExpr) SymbolExpr {
// }

// func RemoveSingleAtomLists(expr SymbolExpr)
// 
