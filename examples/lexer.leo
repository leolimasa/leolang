import "../test.leo" tst

type-alias Name String

union TokenType
    LeftParen
    RightParen
    Identifier
    String
    Number
    Indent
    Dedent
    BinOperator
    

union Token
    Token
        type TokenType
        line Int
        col Int
        value Int


union Lexer
    Lexer
        line Int
	col Int
	rune-size Int
	indent-level Int


:: Fn io.Reader Lexer
new-lexer = fn (reader)
    Lexer
        reader (bufio.new-reader reader)
	line 1
	col 0


:: Fn Lexer TokenType String Token
new-token = fn (lexer token-type value)
    Token
        type token-type
	line lexer.line
	col lexer.col
	value value

:: Fn Lexer 
    Result 
        Optional Token
        LexerError
process-character = fn (lexer)
    rune = use (try (lexer.reader.read-rune))
    if (rune == '\n')
    	lexer = set lexer.col 0
	lexer = set lexer.line (lexer.line + 1)
    

:: Fn Lexer (Result Token LexerError)
next = fn (lexer)
    token = use (try (process-character lexer))
    switch token
	case (Ok t)
	    Ok token
    	case None
	    next lexer
    	
     


