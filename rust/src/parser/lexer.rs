use regex::Regex;
use regex_macro::regex;

#[derive(Debug, Copy, Clone, PartialEq)]
pub struct Loc {
    line: i32,
    col: i32
}

#[derive(Debug, PartialEq)]
pub enum TokenType {
    OpenParen,
    CloseParen,
    Identifier(String),
    StringLiteral(String),
    IntLiteral(i32),
    FloatLiteral(f32),
    Indent(i32),
    Dedent(i32),
    Operator(String),
    LineEnd
}

#[derive(Debug, PartialEq)]
pub struct Token {
    pub token_type: TokenType,
    pub loc: Loc
}

pub struct Lexer {
    stream: Box<dyn Iterator<Item = char>>,
    loc: Loc,
    cur_char: Option<char>,
    skip_next_read: bool,
    indent_stack: Vec<i32>,
    operators: Vec<String>,
    int_regex: Regex,
    float_regex: Regex
}

#[derive(Debug)]
pub enum LexerErrorType {
    InvalidIntLiteral(String),
    InvalidFloatLiteral(String)
}

#[derive(Debug)]
pub struct LexerError {
    pub error_type: LexerErrorType,
    pub loc: Loc
}

pub fn new_lexer(stream: Box<dyn Iterator<Item = char>>) -> Lexer {
    let ops = vec![
        "+", "-", "*", "/", "^", "%", 
        "and", "or", "not",
        "==", "!=", ">", "<", "<=", ">=",
        "=", "v="
    ];
    let int_regex = regex!(r"^-?\d+$");
    let float_regex = regex!("[+-]?([0-9]*[.])?[0-9]+");
    Lexer {
        stream,
        loc: Loc { line: 1, col: 0},
        skip_next_read: false,
        cur_char: None,
        indent_stack: Vec::new(),
        operators: ops.iter().map(|i| i.to_string()).collect(),
        int_regex: int_regex.clone(),
        float_regex: float_regex.clone()
    }
}

fn new_token(lexer: &Lexer, token_type: TokenType) -> Token {
    Token {
        token_type,
        loc: lexer.loc 
    }
}

fn get_indent_level(indent_stack: &Vec<i32>) -> i32 {
    indent_stack.iter().fold(0, |result, &i| result + i)
}

fn next_char(lexer: &mut Lexer) {
    if lexer.skip_next_read {
        lexer.skip_next_read = false;
        return;
    }

    lexer.cur_char = lexer.stream.next();
    if lexer.cur_char == None {
        return;
    }
    lexer.loc.col += 1;
    return;
}

fn dedent_level(indent_stack: &mut Vec<i32>, level: i32, dedents: i32) -> i32 {
    if get_indent_level(indent_stack) <= level {
        return dedents;
    }
    indent_stack.pop();
    dedent_level(indent_stack, level, dedents + 1)
}

fn detect_indent(lexer: &mut Lexer, level: i32) -> Option<Token> {
    next_char(lexer);
    let is_empty_line = match lexer.cur_char {
        Some(c) => c == '\n',
        None => false
    };

    // Empty line. Ignore.
    if is_empty_line {
        lexer.loc.line += 1;
        return None;
    }

    let is_end = match lexer.cur_char {
        Some(c) => c != ' ',
        None => true
    };

    // If we reached the end of the indentation, update the lexer
    // with the new indentation stack and emit the approriate token.
    if is_end {
        lexer.skip_next_read = true;
        let indent_level = get_indent_level(&lexer.indent_stack);

        // Indent detection
        if level > indent_level {
            lexer.indent_stack.push(level - indent_level);
            return Some(new_token(lexer, TokenType::Indent(1)));
        }

        // Dedent detection
        if level < indent_level {
            let mut stack = &mut lexer.indent_stack;
            let dedents = dedent_level(&mut stack, level, 0);
            return Some(new_token(lexer, TokenType::Dedent(dedents)))
        }

        // Same indent as previous line. Do nothing.
        return None;
    }

    detect_indent(lexer, level + 1)
}

fn detect_string(lexer: &mut Lexer, val: String) -> Option<Token> {
    next_char(lexer);
    let c = lexer.cur_char?;
    if c == '"' {
        return Some(new_token(lexer, TokenType::StringLiteral(val)));
    } 
    detect_string(lexer, val + &c.to_string())
}

fn new_error(lexer: &Lexer, error_type: LexerErrorType) -> LexerError {
    LexerError {
        loc: lexer.loc,
        error_type
    }
}

fn detect_ident_or_literal(lexer: &mut Lexer, mut val: String) -> Result<Option<Token>, LexerError> {

    // Check whether we reached the end of the identifier
    let mut is_eof = false;
    let is_end = match lexer.cur_char {
        Some(c) => c == ' ' || c == '\n' || c == '\r' || c == ')',
        None => {
            is_eof = true;
            true
        }
    };
    
    if is_end {
        // We consumed the end character for the identifier, so make sure
        // it gets processed for the following next call
        lexer.skip_next_read = true;

        // remove last character from value (as long as it's not eof - eof is not a char)
        if !is_eof {
            val.pop();
        }
        
        // return a binary op token if the identifier is a binary operator
        if lexer.operators.contains(&val) {
            return Ok(Some(new_token(lexer, TokenType::Operator(val))));
        }

        // check for int literal
        if lexer.int_regex.is_match(&val) {
            let number = val.parse();
            return
                match number {
                    Ok(n) => Ok(Some(new_token(lexer, TokenType::IntLiteral(n)))),
                    Err(_) => Err(new_error(lexer, LexerErrorType::InvalidIntLiteral(val)))
                } 
        }

        // check for float literal
        if lexer.float_regex.is_match(&val) {
            let number = val.parse();
            return
                match number {
                    Ok(n) => Ok(Some(new_token(lexer, TokenType::FloatLiteral(n)))),
                    Err(_) => Err(new_error(lexer, LexerErrorType::InvalidFloatLiteral(val)))
                } 
        }

        // regular identifier
        return Ok(Some(new_token(lexer, TokenType::Identifier(val))));
    }

    // continue reading
    next_char(lexer);
    let mut v = val.clone();
    if let Some(c) = lexer.cur_char {
        v = val + &c.to_string();
    }
    detect_ident_or_literal(lexer, v)
}

impl Iterator for Lexer {
    type Item = Result<Token, LexerError>;
    fn next(&mut self) -> Option<Self::Item> {
        // Grab next character
        next_char(self);
        let c = self.cur_char?;

        // new line 
        if c == '\n' {
            self.loc.col = 0;
            self.loc.line += 1;
            let indent_token = detect_indent(self, 0);
            if let Some(tok) = indent_token {
                return Some(Ok(tok));
            }
            return Some(Ok(new_token(self, TokenType::LineEnd)));
        } 

        // detect parens
        if c == '(' {
            return Some(Ok(new_token(self, TokenType::OpenParen)));
        }
        if c == ')' {
            return Some(Ok(new_token(self, TokenType::CloseParen)));
        }

        // detect string
        if c == '"' {
            let string = detect_string(self, "".to_string());
            if let Some(t) = string {
                return Some(Ok(t));
            }
        }

        // detect identifier or literals - anything that's not a space and none 
        // of the tokens above.
        if c != ' ' && c != '\t' && c != '\n' {
            let ident = detect_ident_or_literal(self, c.to_string());
            match ident {
                Ok(opt) => match opt {
                    Some(t) => return Some(Ok(t)),
                    None => return None
                }
                Err(err) => return Some(Err(err))
            };
        }

        // continue reading
        self.next()
    }
}

struct LexerNoIndent {
    lexer: Lexer,
    token_buffer: Vec<Token>,
    next_is_open_paren: bool
}

fn new_lexer_no_indent(mut lexer: Lexer) -> LexerNoIndent {
    let first_token = Token {
        token_type: TokenType::OpenParen,
        loc: Loc { col:0, line:0},
    };
    LexerNoIndent { lexer, token_buffer: vec![first_token], next_is_open_paren: false }
}

impl Iterator for LexerNoIndent {
    type Item = Result<Token, LexerError>;
    fn next(&mut self) -> Option<Self::Item> {
        // If there is a token in the token buffer, return that.
        if let Some(t) = self.token_buffer.pop() {
            return Some(Ok(t));
        }

        // Advance next token. Return if none or error.
        let tok_res = self.lexer.next()?;
        if let Err(tok_err) = tok_res {
            return Some(Err(tok_err));
        }
        let Ok(tok) = tok_res;

        // emit a open paren if the previous iteration told it to
        if self.next_is_open_paren {
            self.next_is_open_paren = false;
            self.token_buffer.push(tok);
            return Some(Ok(new_token(&self.lexer, TokenType::OpenParen)));
        }

        match tok.token_type {
            TokenType::LineEnd => {
                self.next_is_open_paren = true;
                Some(Ok(new_token(&self.lexer, TokenType::CloseParen)))
            },
            TokenType::Indent(_) => {
                self.next_is_open_paren = true;
                self.next()
            },
            TokenType::Dedent(level) => {
                // emit one dedent token right away and add any remaining
                // dedents to the token buffer
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_lexer() {
        let program = r#"my-fun = fn (a b)
                print "hello world"

                map
                    123 43.74
                some-call
                line-ends-here
            dedented-all-the-way
                indent-one-level
                    indent-two-levels
            dedent-again"#;
        let mut lexer = new_lexer(Box::new(program.chars()));

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("my-fun".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Operator("=".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("fn".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::OpenParen, tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("a".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("b".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::CloseParen, tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Indent(1), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("print".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::StringLiteral("hello world".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::LineEnd, tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("map".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Indent(1), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::IntLiteral(123), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::FloatLiteral(43.74), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Dedent(1), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("some-call".to_string()), tok.token_type); 
        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::LineEnd, tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("line-ends-here".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Dedent(1), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("dedented-all-the-way".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Indent(1), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("indent-one-level".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Indent(1), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("indent-two-levels".to_string()), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Dedent(2), tok.token_type); 

        let tok = lexer.next().unwrap().unwrap();
        assert_eq!(TokenType::Identifier("dedent-again".to_string()), tok.token_type); 

        let tok = lexer.next();
        assert!(tok.is_none());
    }
}
