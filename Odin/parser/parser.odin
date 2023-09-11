package parser

import "../lexer"
import "../token"

Parser :: struct {
    lex: ^lexer.Lexer,
    cur_token: token.Token,
    next_token: token.Token,
}

new_parser :: proc(lex: ^lexer.Lexer) -> ^Parser {
    p := new(Parser)
    p.lex = lex

    return p
}

p_next_token :: proc(p: ^Parser) {
    p.cur_token = p.next_token
    p.next_token = lexer.next_token(p.lex)
}

parse_program :: proc(p: ^Parser) -> 
