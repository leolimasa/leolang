package parser

// NodeType is the type as calculated by the type inference
type NodeType struct {
	TypeVar *string
	ConcreteType *string
	FunType *[]string
}

type ModuleNode struct {
	Name string
	Definitions map[string]DefinitionNode
}

type ApplyNode struct {
	Name string
	Args []ExprNode
	Type *NodeType
}

type ExprNode struct {
	Type *NodeType

	// Union
	Apply *ApplyNode
	Operator *OperatorNode
	Ifexpr *IfNode
}

type BlockItemNode struct {
	Definition *DefinitionNode
	Expr *ExprNode
}

type BlockNode struct {
	Body []BlockItemNode
	Type *NodeType
}

type OperatorNode struct {
	Operator string
	Lhs ExprNode
	Rhs ExprNode
	Type *NodeType
}

type IfNode struct { 
	Condition ExprNode
	IfTrue ExprNode
        IfFalse *ExprNode
	Type *NodeType
}

type FunctionDefNode struct {
	Args []string
	Body []BlockNode
	Type *NodeType
}

type TypeDefNode struct {
	TypeExpr []string	
}

type StructNode struct {
	Name string
	Fields map[string]TypeDefNode
}

type EnumNode struct {
	Name string
}

type DefinitionNode struct {
	Name string

	// Union
	Function *FunctionDefNode
	StructExpr *StructNode
	EnumExpr *EnumNode
	VariableExpr *ExprNode
}
