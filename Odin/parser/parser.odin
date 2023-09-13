package parser

import "core:fmt"
import "../lexer"
import "../token"
import "../ast"

prec :: enum {
    LOWEST,
    EQUALS, // ==
    LESSGREATER, // < , >
    SUM, // +
    PRODUCT, // *
    PREFIX, // -X, !X
    CALL, // my_func(X)
    INDEX, // array[index]
}

precedences := map[token.TokenType]prec {
    token.EQ = prec.EQUALS,
    token.NOT_EQ = prec.EQUALS,
    token.LT = prec.LESSGREATER,
    token.GT = prec.LESSGREATER,
    token.PLUS = prec.SUM,
    token.MINUS = prec.SUM,
    token.ASTERISK = prec.PRODUCT,
    token.SLASH = prec.PRODUCT,
    token.LPAREN = prec.CALL,
    token.LBRACKET = prec.INDEX,
}

prefix_parse_func :: proc() -> ^ast.Expr
infix_parse_func :: proc(^ast.Expr) -> ^ast.Expr

Parser :: struct {
    lex: ^lexer.Lexer,
    cur_token: token.Token,
    next_token: token.Token,

    errors: [dynamic] string,

    prefix_parse_funcs: map[token.TokenType]prefix_parse_func,
    infix_parse_funcs: map[token.TokenType]infix_parse_func,
}

new_parser :: proc(lex: ^lexer.Lexer) -> ^Parser {
    p := new(Parser, context.temp_allocator)
    p.lex = lex
    p.errors = make([dynamic]string, context.temp_allocator)
    
    p.prefix_parse_funcs = make(map[token.TokenType]prefix_parse_func)

    p.infix_parse_funcs = make(map[token.TokenType]infix_parse_func)

    return p
}

// Helpers

p_next_token :: proc(p: ^Parser) {
    p.cur_token = p.next_token
    p.next_token = lexer.next_token(p.lex)
}

p_peek_travarse :: proc(p: ^Parser, t: token.TokenType) -> bool {
    if p.cur_token.type == t {
        p_next_token(p)
        return true
    } else {
        p_peek_error(p, t)
        return false
    }
}

p_peek_error :: proc(p: ^Parser, t: token.TokenType) {
    buf: []u8
    message := fmt.bprintf(buf, "expected next token to be {%s}, got {%s} instead",
    t, p.next_token.type)
    append(&p.errors, message)
}

curr_precedence :: proc(p: ^Parser) -> prec {
    if p, ok := precedences[p.cur_token.type]; ok {
        return p
    }
    return prec.LOWEST
}

next_precedence :: proc(p: ^Parser) -> prec {
    if p, ok := precedences[p.next_token.type]; ok {
        return p
    }
    return prec.LOWEST
}

// parse functions

parse_program :: proc(p: ^Parser) -> ^ast.Program {
    program := ast.new_node(ast.Program) // default allocator is "context.temp_allocator"
    program.statements = make([dynamic]^ast.Stmt, context.temp_allocator)

    for p.cur_token.type != token.EOF {
        stmt := parse_stmt(p)
        if stmt != nil {
            append(&program.statements, stmt)
        }
        p_next_token(p)
    }

    return program
}

// PARSING STATEMENTS

parse_stmt :: proc(p: ^Parser) -> ^ast.Stmt {
    switch p.cur_token.type {
        case token.LET: return parse_let_stmt(p)
        case token.RETURN: return parse_return_stmt(p)
        case: return parse_expr_stmt(p)
    }
}

parse_expr_stmt :: proc(p: ^Parser) -> ^ast.Stmt {
    expr_stmt := ast.new_node(ast.Expr_Stmt)

    expr := parse_expr(p, prec.LOWEST)
    expr_stmt.expr = expr

    if p_peek_travarse(p, token.SEMICOLON) {
        p_next_token(p)
    }

    return expr_stmt
}

parse_let_stmt :: proc(p: ^Parser) -> ^ast.Stmt {
    let := ast.new_node(ast.Let_Stmt)

    if !p_peek_travarse(p, token.IDENT) {
        return nil
    }

    name := ast.new_node(ast.Ident)
    name.token = p.cur_token
    name.value = p.cur_token.literal
    let.name = name

    if !p_peek_travarse(p, token.ASSIGN) {
        return nil
    }

    p_next_token(p)

    // LOWEST precedence is required at the start of parsing expressions.
    // Also whenever a precedence doesn't exist in the precedence map, we pass LOWEST.
    value := parse_expr(p, prec.LOWEST)
    let.value = value

    // Optional semicolon
    if p_peek_travarse(p, token.SEMICOLON) {
        p_next_token(p)
    }

    return let
}

parse_return_stmt :: proc(p: ^Parser) -> ^ast.Return_Stmt {
    ret := ast.new_node(ast.Return_Stmt)

    // We skip return
    p_next_token(p)

    ret_value := ast.new_node(ast.Expr)
    ret.return_value = parse_expr(p, prec.LOWEST)

    if p_peek_travarse(p, token.SEMICOLON) {
        p_next_token(p)
    }

    return ret
}

// PARSING EXPRESSIONS

// pratt parsing
parse_expr :: proc(p: ^Parser, precedenc: prec) -> ^ast.Expr {
    return nil
}
