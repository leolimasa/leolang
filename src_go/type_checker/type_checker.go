package type_checker

import (
	"leolang/parser"
)

type TypeKind int

const (
	Struct TypeKind = 1
	Enum TypeKind = 2
	Int TypeKind = 3
	Float TypeKind = 4
	Function TypeKind = 5
)

type Type struct {
	Name string
	Kind TypeKind
	Value any
}

type TypedExpr struct {
	Atom *parser.Token
	Type *Type
	List []TypedExpr
}

type TypeChecker struct {
}

// func Typecheck(expr parser.SymbolExpr, symbolTable map[string]Type) (*TypedExpr, *parser.ParserError) {
// 	// Extend the symbol table with local definitions
// 	if expr.Atom == nil {
// 		symbolTable, err := SymbolTableFromAssignments(expr.List)
// 		if err != nil {
// 			return err
// 		} 
// 	}
//
// 	if expr.Atom != nil {
// 		if expr.Atom == "=" {
// 			
// 		}		
// 	}
//
// 	for expr, _ := range expr.List {
// 	}
// }
