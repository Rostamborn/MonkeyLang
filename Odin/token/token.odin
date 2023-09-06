package token

ILLEGAL :: "ILLEGAL"
EOF :: "EOF"
// Identifiers and lierals
IDENT :: "IDENT"
INT :: "INT"
STRING :: "STRING"
// Operators
ASSIGN :: "="
PLUS :: "+" 
MINUS :: "-"
ASTERISK :: "*"
SLASH :: "/"
BANG :: "!"
LT :: "<"
GT :: ">"
EQ :: "=="
NOT_EQ :: "!="
// Delimiters
COMMA :: ","
COLON :: ":"
SEMICOLON :: ";"
LPAREN :: "("
RPAREN :: ")"
LBRACE :: "{"
RBRACE :: "}"
LBRACKET :: "["
RBRACKET :: "]"
// Keywords
FUNCTION :: "FUNCTION"
LET :: "LET"
IF :: "IF"
ELSE :: "ELSE"
RETURN :: "RETURN"
TRUE :: "TRUE"
FALSE :: "FALSE"

TokenType :: distinct string
//
keywords := map[string]TokenType{ // is this heap??
    "fn" = FUNCTION,
    "let" = LET,
    "if" = IF,
    "else" = ELSE,
    "true" = TRUE,
    "false" = FALSE,
    "return" = RETURN,
}

Token :: struct {
    type: TokenType,
    literal: string,
}

new_token :: proc(token_type: TokenType, ch: byte, allocator := context.temp_allocator) -> Token {
    buf := make([]u8, 1)
    buf[0] = ch
    return Token{type=token_type, literal=transmute(string)buf}
}

lookup_identifier :: proc(ind: string) -> TokenType {
    if tok, ok := keywords[ind]; ok {
        return tok
    }
    return IDENT
}
