# ----------------------------------------
# Std library
# ----------------------------------------

Length = trait T
    len (Fn T Int)

Wrap = trait T V
    wrap (Fn V T)

Unwrap = trait T V
    unwrap (Fn T (Fn V V))

UnwrapThenWrap = trait T V
    unwrap-then-wrap (Fn T (Fn V T))

Bool = enum False True

Result = generic (K E)
    enum
	Ok K
	Error E

implement Unwrap Result
    unwrap (value cb)
        case value
	    Ok value
	        cb value
	    Err err
	        err

implement UnwrapThenWrap Result
    bare-unwrap (value cb)
        case value
	    Ok value
                Ok (cb value)
	    Err err
	        err

# ----------------------------------------
# Actual code
# ----------------------------------------

Token = enum
    LeftParen
    RightParen
    Identifier Text
    String
    IntTok
    FloatTok
    Indent Int32
    Dedent Int32
    Operator String
    

TokenPos = struct
    token Token
    line Int32
    col Int32


Lexer = struct
    line Int32
    col Int32
    cur-rune Rune
    skip-next-read Bool
    indent-stack (List Int)
    is-eof Bool


operators = list
    "+" "-" "*" "/" "**" "%"
    "and" "or" "not"
    "==" "!=" ">" "<" "<=" ">="
    "=" "v="


:: Fn io.Reader Lexer
new-lexer = fn (reader)
    Lexer
        reader ((bufio new-reader) reader)
	line 1
	col 0
	skip-next-read False
	is-eof False


:: Fn Lexer Token String TokenPos
new-token-pos = fn (lexer token value)
    Token
        token token
	line lexer.line
	col lexer.col


sum = fn (lst) (reduce (fn (a b) (a + b)) lst) 

get-indent-level = fn (indent-stack) (sum indent-stack)

calculate-dedent = fn (indent-stack level dedents)
    if (get-indent-level indent-stack <= level)
        return (tuple indent-stack level)
    calculate-dedent (indent-stack 0 -1) level (dedents + 1)


detect-indent = fn (lexer level)
    lexer v= (read-rune lexer) unwrap

    # Empty line. Ignore.
    if (lexer cur-rune == '\n')
        lexer = lexer (line (lexer line + 1))
        return (Ok (tuple None lexer))

    # Check if we reached the end of the indentation
    if (lexer cur-rune != ' ' or lexer is-eof)

    	lexer = lexer (skip-next-read true)
	indent-level = get-indent-level (lexer indent-stack)

	if (level > indent-level)
	    lexer indent-stack = append 
	    	lexer indent-stack 
		level - indent-level
	    return (Ok (tuple (Indent 1) lexer))

	if (level < indent-level)
	    dedents = calculate-dedent lexer.indent-stack 0
	    return (Ok (tuple (Dedent dedents) lexer))

    detect-indent lexer (level + 1)
	    
    
process-character = fn (lexer)
    rune v= unwrap (lexer reader read-rune)
    if rune == '\n' then
    	lexer = lexer
	    col 0
	    line (lexer line + 1)

:: Fn Lexer (Result Token LexerError)
next = fn (lexer)
    token v= (process-character lexer) unwrap
    switch token
	case (Ok t)
	    Ok token
    	case None
	    next lexer
    	

LexerNoIndent = struct
    lexer Lexer
    token-buffer (List Token)
    next-is-open-paren Bool


next-no-indent = fn (lexer-ni)
    if ((lexer-ni token-buffer) len > 0)

