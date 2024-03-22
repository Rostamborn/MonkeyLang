package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
    INT    = "INT"   // 1343456
    STRING = "STRING" // "mate", "mamad"
	// Operators
	ASSIGN = "="
	PLUS   = "+"
    MINUS  = "-"
    BANG   = "!"
    ASTERISK = "*"
    SLASH    = "/"
    LT       = "<"
    GT       = ">"
    EQ       = "=="
    NOT_EQ   = "!="
	// Delimiters
	COMMA     = ","
    COLON     = ":"
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
    LBRACKET  = "["
    RBRACKET  = "]"
	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
    IF       = "IF"
    ELSE     = "ELSE"
    RETURN   = "RETURN"
    TRUE     = "TRUE"
    FALSE    = "FALSE"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
    "if":  IF,
    "else": ELSE,
    "return": RETURN,
    "true": TRUE,
    "false": FALSE,
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func LookupIndentifier(indent string) TokenType {
	if tok, ok := keywords[indent]; ok {
		return tok
	}
	return IDENT
}
