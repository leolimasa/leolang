mod parser;

fn main() {
    let code = "my-fun = fn (a b) (c d)";  
    let mut lexer = parser::lexer::new_lexer(Box::new(code.chars()));
    let tok = lexer.next();
    match tok {
        Some(opt) => match opt {
            Ok(t) => {
                println!("{:?}", t.token_type); 
            }
            Err(_) => {
                println!("Error!");
            } 
        },
        None => println!("no token!")
    } 
}
