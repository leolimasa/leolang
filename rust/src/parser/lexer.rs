use regex::Regex;
use regex_macro::regex;

#[derive(Debug, Copy, Clone)]
pub struct Loc {
    line: i32,
    col: i32
}

#[derive(Debug)]
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

#[derive(Debug)]
pub struct Token {
    pub token_type: TokenType,
    loc: Loc
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
    error_type: LexerErrorType,
    loc: Loc
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
    let c = lexer.cur_char?;

    // Empty line. Ignore.
    if c == '\n' {
        lexer.loc.line += 1;
        return None;
    }

    // If we reached the end of the indentation, update the lexer
    // with the new indentation stack and emit the approriate token.
    if c != ' ' {
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

fn detect_identifier(lexer: &mut Lexer, mut val: String) -> Result<Option<Token>, LexerError> {
    let c = match lexer.cur_char {
        Some(c) => c,
        None => return Ok(None)
    };

    // Check whether we reached the end of the identifier
    if c == ' ' || c == '\n' || c == '\r' || c == ')' {
        // We consumed the end character for the identifier, so make sure
        // it gets processed for the following next call
        lexer.skip_next_read = true;

        // remove last character from value
        val.pop();
        
        // return a binary op token if the identifier is a binary operator
        if lexer.operators.contains(&val) {
            return Ok(Some(new_token(lexer, TokenType::Operator(val))));
        }

        // check for int
        if lexer.int_regex.is_match(&val) {
            let number = val.parse();
            return
                match number {
                    Ok(n) => Ok(Some(new_token(lexer, TokenType::IntLiteral(n)))),
                    Err(_) => Err(new_error(lexer, LexerErrorType::InvalidIntLiteral(val)))
                } 
        }

        // check for float
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
    let c = match lexer.cur_char {
        Some(c) => c,
        None => return Ok(None)
    };
    detect_identifier(lexer, val + &c.to_string())
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

        // detect identifier - anything that's not a space and none of the
        // tokens above.
        if c != ' ' && c != '\t' && c != '\n' {
            let ident = detect_identifier(self, "".to_string());
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

