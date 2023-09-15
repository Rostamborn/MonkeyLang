package parser

import "core:strconv"
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

prefix_parse_func :: proc(^Parser) -> ^ast.Expr
infix_parse_func :: proc(^Parser, ^ast.Expr) -> ^ast.Expr

Parser :: struct {
    lex: ^lexer.Lexer,
    curr_token: token.Token,
    next_token: token.Token,

    errors: [dynamic]string,

    prefix_parse_funcs: map[token.TokenType]prefix_parse_func,
    infix_parse_funcs: map[token.TokenType]infix_parse_func,
}

new_parser :: proc(lex: ^lexer.Lexer) -> ^Parser {
    p := new(Parser, context.temp_allocator)
    p.lex = lex
    p.errors = make([dynamic]string, context.temp_allocator)
    p_next_token(p)
    p_next_token(p)
    
    p.prefix_parse_funcs = make(map[token.TokenType]prefix_parse_func, 2, context.temp_allocator)
    p.prefix_parse_funcs[token.IDENT] = parse_ident
    p.prefix_parse_funcs[token.INT] = parse_int_literal
    p.prefix_parse_funcs[token.STRING] = parse_string_literal
    p.prefix_parse_funcs[token.TRUE] = parse_bool_literal
    p.prefix_parse_funcs[token.FALSE] = parse_bool_literal
    p.prefix_parse_funcs[token.BANG] = parse_prefix_expr
    p.prefix_parse_funcs[token.MINUS] = parse_prefix_expr
    p.prefix_parse_funcs[token.LPAREN] = parse_grouped_expr
    p.prefix_parse_funcs[token.IF] = parse_if_expr
    p.prefix_parse_funcs[token.FUNCTION] = parse_function_literal
    p.prefix_parse_funcs[token.LBRACKET] = parse_array_literal
    p.prefix_parse_funcs[token.LBRACE] = parse_hash_expr

    p.infix_parse_funcs = make(map[token.TokenType]infix_parse_func, 2, context.temp_allocator)
    p.infix_parse_funcs[token.PLUS] = parse_infix_expr
    p.infix_parse_funcs[token.MINUS] = parse_infix_expr
    p.infix_parse_funcs[token.ASTERISK] = parse_infix_expr
    p.infix_parse_funcs[token.SLASH] = parse_infix_expr
    p.infix_parse_funcs[token.EQ] = parse_infix_expr
    p.infix_parse_funcs[token.NOT_EQ] = parse_infix_expr
    p.infix_parse_funcs[token.LT] = parse_infix_expr
    p.infix_parse_funcs[token.GT] = parse_infix_expr
    p.infix_parse_funcs[token.LPAREN] = parse_call_expr
    p.infix_parse_funcs[token.LBRACKET] = parse_index_expr

    return p
}

// Helpers

p_next_token :: proc(p: ^Parser) {
    p.curr_token = p.next_token
    p.next_token = lexer.next_token(p.lex)
}

p_peek_travarse :: proc(p: ^Parser, t: token.TokenType) -> bool {
    if p.next_token.type == t {
        p_next_token(p)
        return true
    } else {
        p_peek_error(p, t)
        return false
    }
}

p_peek_error :: proc(p: ^Parser, t: token.TokenType) {
    message := fmt.tprintf("expected next token to be (%s), got (%s) instead",
    t, p.next_token.type)
    append(&p.errors, message)
}

no_prefix_func_error :: proc(p: ^Parser, t: token.TokenType) {
    buf: []u8
    message := fmt.bprintf(buf, "no prefix_parse_func found for {%s}", t)
    append(&p.errors, message)
}

curr_precedence :: proc(p: ^Parser) -> prec {
    if p, ok := precedences[p.curr_token.type]; ok {
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

// PARSING

parse_program :: proc(p: ^Parser) -> ^ast.Program {
    program := ast.new_node(ast.Program) // default allocator is "context.temp_allocator"
    program.statements = make([dynamic]^ast.Stmt, context.temp_allocator)

    for p.curr_token.type != token.EOF {
        // fmt.println(p.curr_token)
        stmt := parse_stmt(p) 
        if stmt != nil {
            append(&program.statements, stmt)
        }
        p_next_token(p) // we always stop before a new statement,
    }               // so we must go to the next token.

    return program
}

// PARSING STATEMENTS

parse_stmt :: proc(p: ^Parser) -> ^ast.Stmt {
    switch p.curr_token.type {
        case token.LET: return parse_let_stmt(p)
        case token.RETURN: return parse_return_stmt(p)
        case: return parse_expr_stmt(p)
    }
}

parse_expr_stmt :: proc(p: ^Parser) -> ^ast.Stmt {
    expr_stmt := ast.new_node(ast.Expr_Stmt)
    expr_stmt.token = p.curr_token

    expr := parse_expr(p, prec.LOWEST)
    expr_stmt.expr = expr

    if p.next_token.type == token.SEMICOLON {
        p_next_token(p)
    }

    return expr_stmt
}

parse_let_stmt :: proc(p: ^Parser) -> ^ast.Stmt {
    let := ast.new_node(ast.Let_Stmt)
    let.token = p.curr_token

    if !p_peek_travarse(p, token.IDENT) {
        return nil
    }

    name := ast.new_node(ast.Ident)
    name.token = p.curr_token
    name.value = p.curr_token.literal
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
    if p.next_token.type == token.SEMICOLON {
        p_next_token(p)
    }

    return let
}

parse_return_stmt :: proc(p: ^Parser) -> ^ast.Return_Stmt {
    ret := ast.new_node(ast.Return_Stmt)
    ret.token = p.curr_token

    // We skip return
    p_next_token(p)

    // ret_value := ast.new_node(ast.Expr)
    ret.return_value = parse_expr(p, prec.LOWEST)

    if p.next_token.type == token.SEMICOLON {
        p_next_token(p)
    }

    return ret
}

parse_block_stmts :: proc(p: ^Parser) -> ^ast.Block_Stmt {
    block := ast.new_node(ast.Block_Stmt)
    block.token = p.curr_token
    block.statements = make([dynamic]^ast.Stmt, context.temp_allocator)

    p_next_token(p)

    for p.curr_token.type != token.RBRACE {
        stmt := parse_stmt(p)
        if stmt != nil {
            append(&block.statements, stmt)
        }

        p_next_token(p)
    }

    return block
}

// PARSING EXPRESSIONS

// pratt parsing
parse_expr :: proc(p: ^Parser, precedenc: prec) -> ^ast.Expr {
    prefix := p.prefix_parse_funcs[p.curr_token.type]
    if prefix == nil {
        no_prefix_func_error(p, p.curr_token.type)
        return nil
    }

    left_expr := prefix(p)

    for precedenc < next_precedence(p) {
        infix, ok := p.infix_parse_funcs[p.next_token.type]
        if !ok {
            return left_expr
        }

        p_next_token(p)

        left_expr = infix(p, left_expr)
    }

    return left_expr
}

parse_ident :: proc(p: ^Parser) -> ^ast.Expr {
    ident := ast.new_node(ast.Ident)
    ident.token = p.curr_token
    ident.value = p.curr_token.literal
    
    return ident
}

parse_int_literal :: proc(p: ^Parser) -> ^ast.Expr {
    lit := ast.new_node(ast.Int_Literal)
    lit.token = p.curr_token
    lit.value = strconv.atoi(p.curr_token.literal)

    return lit
}

parse_string_literal :: proc(p: ^Parser) -> ^ast.Expr {
    lit := ast.new_node(ast.String_Literal) // we pass the double-quotes in the lexing
    lit.token = p.curr_token                // phase, so no need to check for them :D
    lit.value = p.curr_token.literal

    return lit
}

parse_bool_literal :: proc(p: ^Parser) -> ^ast.Expr {
    lit := ast.new_node(ast.Bool_Literal)
    lit.token = p.curr_token
    lit.value = p.curr_token.type == token.TRUE // neat!

    return lit
}

parse_prefix_expr :: proc(p: ^Parser) -> ^ast.Expr {
    prefix_expr := ast.new_node(ast.Prefix_Expr) 
    prefix_expr.token = p.curr_token
    prefix_expr.operator = p.curr_token.literal // ! , -

    p_next_token(p)

    prefix_expr.right = parse_expr(p, prec.PREFIX)

    return prefix_expr
}

parse_infix_expr :: proc(p: ^Parser, left: ^ast.Expr) -> ^ast.Expr {
    infix_expr := ast.new_node(ast.Infix_Expr)
    infix_expr.token = p.curr_token
    infix_expr.operator = p.curr_token.literal
    infix_expr.left = left
    precedence := curr_precedence(p)

    p_next_token(p)

    infix_expr.right = parse_expr(p, precedence)

    return infix_expr
}

parse_grouped_expr :: proc(p: ^Parser) -> ^ast.Expr {
    p_next_token(p)

    expr := parse_expr(p, prec.LOWEST)

    if !p_peek_travarse(p, token.RPAREN) {
        return nil
    }

    return expr
}

parse_if_expr :: proc(p: ^Parser) -> ^ast.Expr {
    expr := ast.new_node(ast.If_Expr)
    expr.token = p.curr_token
    expr.alternatives = make([dynamic]^ast.If_Expr, context.temp_allocator)

    if !p_peek_travarse(p, token.LPAREN) {
        return nil
    }
    p_next_token(p)

    expr.condition = parse_expr(p, prec.LOWEST)

    if !p_peek_travarse(p, token.RPAREN) {
        return nil
    }

    if !p_peek_travarse(p, token.LBRACE) {
        return nil
    }

    expr.consequence = parse_block_stmts(p) // we are now on LBRACE

    for p.next_token.type == token.ELSE {
        p_next_token(p)

        if p.next_token.type == token.IF {
            p_next_token(p)

            alt_expr := ast.new_node(ast.If_Expr)
            alt_expr.token = p.curr_token

            if p_peek_travarse(p, token.LPAREN) {
                return nil
            }
            p_next_token(p)

            alt_expr.condition = parse_expr(p, prec.LOWEST)

            if !p_peek_travarse(p, token.RPAREN) {
                return nil
            }

            if !p_peek_travarse(p, token.LBRACE) {
                return nil
            }

            alt_expr.consequence = parse_block_stmts(p)

            append(&expr.alternatives, alt_expr)
        } else {
            if !p_peek_travarse(p, token.LBRACE) {
                return nil
            }
            
            expr.default = parse_block_stmts(p)
        }
    }

    return expr
}

parse_function_literal :: proc(p: ^Parser) -> ^ast.Expr {
    func_lit := ast.new_node(ast.Function_Literal)
    func_lit.token = p.curr_token
    // func_lit.params = make([dynamic]^ast.Ident, context.temp_allocator)

    if !p_peek_travarse(p, token.LPAREN) {
        return nil
    }

    func_lit.params = parse_function_params(p)

    if !p_peek_travarse(p, token.LBRACE) {
        return nil
    }

    func_lit.body = parse_block_stmts(p)

    // we stop at RBRACE, because p_next_token() is called in parse_program()

    return func_lit
}

parse_function_params :: proc(p: ^Parser) -> [dynamic]^ast.Ident {
    params := make([dynamic]^ast.Ident, context.temp_allocator)

    // we are on LBRACE now

    if p.next_token.type == token.RPAREN {
        p_next_token(p)
        return params
    }

    p_next_token(p)

    param := ast.new_node(ast.Ident)
    param.token = p.curr_token
    param.value = p.curr_token.literal

    append(&params, param)

    for p.next_token.type == token.COMMA {
        p_next_token(p)
        p_next_token(p)
        
        param := ast.new_node(ast.Ident)
        param.token = p.curr_token
        param.value = p.curr_token.literal

        append(&params, param)
    }

    if !p_peek_travarse(p, token.RPAREN) {
        return nil
    }

    return params
}

parse_call_expr :: proc(p: ^Parser, func: ^ast.Expr) -> ^ast.Expr {
    expr := ast.new_node(ast.Call_Expr)
    expr.token = p.curr_token
    expr.func = func
    expr.args = parse_expr_list(p, token.RPAREN)

    return expr
}

parse_array_literal :: proc(p: ^Parser) -> ^ast.Expr {
    array := ast.new_node(ast.Array_Literal)
    array.token = p.curr_token
    array.elems = parse_expr_list(p, token.RBRACKET)

    return array
}

parse_expr_list :: proc(p: ^Parser, end: token.TokenType) -> [dynamic]^ast.Expr {
    list := make([dynamic]^ast.Expr, context.temp_allocator)
    
    if p.next_token.type == end {
        p_next_token(p)
        return list
    }

    p_next_token(p)

    append(&list, parse_expr(p, prec.LOWEST))

    for p.next_token.type == token.COMMA {
        p_next_token(p)
        p_next_token(p)

        append(&list, parse_expr(p, prec.LOWEST))
    }

    if !p_peek_travarse(p, end) {
        return nil
    }

    return list
}

parse_index_expr :: proc(p: ^Parser, left: ^ast.Expr) -> ^ast.Expr {
    expr := ast.new_node(ast.Index_Expr)
    expr.left = left
    expr.token = p.curr_token

    p_next_token(p)
    
    expr.index = parse_expr(p, prec.LOWEST)

    if !p_peek_travarse(p, token.RBRACKET) {
        return nil
    }

    return expr
}

parse_hash_expr :: proc(p: ^Parser) -> ^ast.Expr {
    expr := ast.new_node(ast.Hash_Expr)
    expr.token = p.curr_token
    expr.pairs = make(map[^ast.Expr]^ast.Expr, 2, context.temp_allocator)

    for p.next_token.type != token.RBRACE {
        p_next_token(p)

        key := parse_expr(p, prec.LOWEST)

        if !p_peek_travarse(p, token.COLON) {
            return nil
        }

        p_next_token(p)

        value := parse_expr(p, prec.LOWEST)

        expr.pairs[key] = value

        // if we don't hit either RBRACE or COMMA, then there is a syntax error
        if p.next_token.type != token.RBRACE && !p_peek_travarse(p, token.COMMA) {
            return nil
        }
    }

    if !p_peek_travarse(p, token.RBRACE) {
        return nil
    }

    return expr
}


