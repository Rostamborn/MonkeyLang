package ast

import "../token"
import "core:bytes"

Any_Node :: union {
    // Stmts
    ^Expr_Stmt,
    ^Let_Stmt,
    ^Return_Stmt,
    ^Block_Stmt,

    // Exprs
    ^Ident,
    ^Int_Literal,
    ^Float_Literal,
    ^String_Literal,
    ^Bool_Literal,
    ^Prefix_Expr,
    ^Infix_Expr,
    ^If_Expr,
    ^Function_Literal,
    ^Call_Expr,
    ^Array_Literal,
    ^Index_Expr,
    ^Hash_Expr,
}

Any_Stmt :: union {
    ^Expr_Stmt,
    ^Let_Stmt,
    ^Return_Stmt,
    ^Block_Stmt,
}

Any_Expr :: union {
    ^Ident,
    ^Int_Literal,
    ^Float_Literal,
    ^String_Literal,
    ^Bool_Literal,
    ^Prefix_Expr,
    ^Infix_Expr,
    ^If_Expr,
    ^Function_Literal,
    ^Call_Expr,
    ^Array_Literal,
    ^Index_Expr,
    ^Hash_Expr,
}


Node :: struct {
    derived: Any_Node,
}

Expr :: struct {
    using expr_base: Node,
    derived_expr: Any_Expr,
}

Stmt :: struct {
    using expr_base: Node,
    derived_stmt: Any_Stmt,
}

Expr_inst :: struct($T: typeid) {
    using expr: ^T,
}

Stmt_inst :: struct($T: typeid) {
    using stmt: ^T,
}

// Program

Program :: struct {
    statemetns: []Stmt,
}

// Stmts

Expr_Stmt :: struct {
    using node: Stmt,
    expr: ^Expr,
}

Let_Stmt :: struct {
    using node: Stmt,
    token: token.Token,
    name: ^Ident,
    value: ^Expr,
}

Return_Stmt :: struct {
    using node: Stmt,
    token: token.Token,
    return_value: ^Expr,
}

Block_Stmt :: struct {
    using node: Stmt,
    token: token.Token,
    statements: []^Stmt,
}

// Exprs

Ident :: struct {
    using node: Expr,
    token: token.Token,
    value: string,
}

Int_Literal :: struct {
    using node: Expr,
    token: token.Token,
    value: int,
}

Float_Literal :: struct {
    using node: Expr,
    token: token.Token,
    value: f64,
}

String_Literal :: struct {
    using node: Expr,
    token: token.Token,
    value: string,
}

Bool_Literal :: struct {
    using node: Expr,
    token: token.Token,
    value: bool,
}

Prefix_Expr :: struct {
    using node: Expr,
    token: token.Token,
    operator: string, // -, !
    right: ^Expr,
}

Infix_Expr :: struct {
    using node: Expr,
    token: token.Token,
    operator: string, // +, -, <, &, etc.
    left: ^Expr,
    right: ^Expr,
}

If_Expr :: struct {
    using node: Expr,
    token: token.Token,
    condition: ^Expr,
    consequence: ^Block_Stmt,
    alternatives: []^If_Expr,
    default: ^Block_Stmt,
}

Function_Literal :: struct {
    using node: Expr,
    token: token.Token,
    params: []^Ident,
    body: ^Block_Stmt,
}

Call_Expr :: struct {
    using node: Expr,
    token: token.Token,
    func: ^Expr, // Ident or Function_Expr
    args: []^Expr,
}

Array_Literal :: struct {
    using node: Expr,
    token: token.Token,
    elems: []^Expr,
}

Index_Expr :: struct {
    using node: Expr,
    token: token.Token,
    left: ^Expr, // Ident or Array_Expr
    Index: ^Expr,
}

Hash_Expr :: struct {
    using node: Expr,
    token: token.Token,
    pairs: map[^Expr]^Expr,
}

// String

to_string :: proc {
    program_string,

    // Stmts
    expr_stmt_string,
    let_stmt_string,
    return_stmt_string,
    block_stmt_string,

    // Exprs
    ident_string,
    int_literal_string,
    float_literal_string,
    string_literal_string,
    bool_literal_string,
    prefix_expr_string,
    infix_expr_string,
    if_expr_string,
    function_expr_string,
    call_expr_string,
    array_expr_string,
    index_expr_string,
    hash_expr_string,

}

program_string :: proc(p: ^Program) -> string {
    out: bytes.Buffer

    for str in p.statemetns {
        bytes.buffer_write(&out, transmute([]u8)to_string(p))
    }

    return bytes.buffer_to_string(&out)
}

// Stmts

expr_stmt_string :: proc(s: ^Expr_Stmt) -> string {
    out: bytes.Buffer

    if s.expr != nil {
        return to_string(s.expr)
    }

    return ""
}

let_stmt_string :: proc(s: ^Let_Stmt) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)tok_literal(s.token)+" ")
    bytes.buffer_write(&out, transmute([]u8)to_string(s.name))
    bytes.buffer_write(&out, transmute([]u8)string(" = ")) // we must convert the untyped string literal to typed string

    if s.value != nil {
        bytes.buffer_write(&out, transmute([]u8)to_string(s.value))
    }

    bytes.buffer_write(&out, transmute([]u8)string(";"))

    return bytes.buffer_to_string(&out)
}

return_stmt_sting :: proc(s: ^Return_Stmt) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)tok_literal(s.token) + " ")

    if s.return_value != nil {
        bytes.buffer_write(&out, transmute([]u8)to_string(s.return_value))
    }

    bytes.buffer_write(&out, transmute([]u8)string(";"))

    return bytes.buffer_to_string(&out)
}

block_stmt_string :: proc(s: ^Block_Stmt) -> string {
    out: bytes.Buffer

    for stmt in s.statements {
        bytes.buffer_write(&out, transmute([]u8)to_string(stmt)) // what should I do?
    }

    return bytes.buffer_to_string(&out)
}

// Exprs

ident_string :: proc(e: ^Ident) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)e.value)
}
