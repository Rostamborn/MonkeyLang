package lexer

import "../token"
import "core:fmt"

Lexer :: struct {
    input: string,
    pos: int,
    next_pos: int,
    ch: byte,
}

new_lexer :: proc(input: string) -> ^Lexer {
    lex := new(Lexer)
    lex.input = input
    next_token(lex)

    return lex
}

next_char :: proc(l: ^Lexer) {
    if l.next_pos >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.next_pos]
    }
    l.pos = l.next_pos
    l.next_pos += 1
}

next_token :: proc(l: ^Lexer) -> token.Token {
    tok: token.Token

    skip_whitespace(l)

    switch l.ch {
        case '=':
        if peek_char(l) == '=' {
            tok.type = token.EQ
            tok.literal = "=="
            next_char(l)
        } else {
            tok = token.new_token(token.ASSIGN, l.ch)
        }
        case ':':
        tok = token.new_token(token.COLON, l.ch)
        case ';':
        tok = token.new_token(token.SEMICOLON, l.ch)
        case '(':
        tok = token.new_token(token.LPAREN, l.ch)
        case ')':
        tok = token.new_token(token.RPAREN, l.ch)
        case ',':
        tok = token.new_token(token.COMMA, l.ch)
        case '{':
        tok = token.new_token(token.LBRACE, l.ch)
        case '}':
        tok = token.new_token(token.RBRACE, l.ch)
        case '[':
        tok = token.new_token(token.LBRACKET, l.ch)
        case ']':
        tok = token.new_token(token.RBRACKET, l.ch)
        case '+':
        tok = token.new_token(token.PLUS, l.ch)
        case '-':
        tok = token.new_token(token.MINUS, l.ch)
        case '!':
        if peek_char(l) == '=' {
            next_char(l)
            tok.type = token.NOT_EQ
            tok.literal = "!="
        } else {
            tok = token.new_token(token.BANG, l.ch)
        }
        case '/':
        tok = token.new_token(token.SLASH, l.ch)
        case '*':
        tok = token.new_token(token.ASTERISK, l.ch)
        case '<':
        tok = token.new_token(token.LT, l.ch)
        case '>':
        tok = token.new_token(token.GT, l.ch)
        case '"':
        next_char(l)
        tok.type = token.STRING
        tok.literal = read_string(l)
        case 0:
        tok.type = token.EOF
        case:
        if is_digit(l.ch) {
            tok.type = token.INT
            tok.literal = read_number(l)
            return tok
        } else if is_letter(l.ch) {
            tok.literal = read_identifier(l)
            tok.type = token.lookup_identifier(tok.literal)
            return tok
        } else {
            tok = token.new_token(token.ILLEGAL, l.ch)
        }
    }
    next_char(l)

    return tok
}

read_number :: proc(l: ^Lexer) -> string {
    pos := l.pos
    for is_digit(l.ch) {
        next_char(l)
    }

    return l.input[pos:l.pos]
}

read_identifier :: proc(l: ^Lexer) -> string {
    pos := l.pos
    for is_letter(l.ch) {
        next_char(l)
    }

    return l.input[pos:l.pos]
}

read_string :: proc(l: ^Lexer) -> string {
    pos := l.pos
    for l.ch != '"' && l.ch != 0 {
        next_char(l)
    }
    
    return l.input[pos:l.pos]
}

is_digit :: proc(ch: byte) -> bool {
    return '0' <= ch && ch <= '9'
}

is_letter :: proc(ch: byte) -> bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

skip_whitespace :: proc(l: ^Lexer) {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        next_char(l)
    }
}

peek_char :: proc(l: ^Lexer) -> byte {
    if l.next_pos >= len(l.input) {
        return 0
    } else {
        return l.input[l.next_pos]
    }
}
