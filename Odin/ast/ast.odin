package ast

import "../token"
import "core:bytes"
import "core:strconv"
import "core:strings"
import "core:fmt"
import "core:intrinsics"

has_field :: intrinsics.type_has_field

new_node :: proc($T: typeid, allocator := context.temp_allocator) -> ^T where has_field(T, "derived") {
    node := new(T)
    node.derived = node

    return node
}

Any_Node :: union {
    ^Program,
    // Stmts
    ^Expr_Stmt,
    ^Let_Stmt,
    ^Return_Stmt,
    ^Block_Stmt,

    // Exprs
    ^Ident,
    ^Int_Literal,
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

// Any_Stmt :: union {
//     ^Expr_Stmt,
//     ^Let_Stmt,
//     ^Return_Stmt,
//     ^Block_Stmt,
// }
//
// Any_Expr :: union {
//     ^Ident,
//     ^Int_Literal,
//     ^String_Literal,
//     ^Bool_Literal,
//     ^Prefix_Expr,
//     ^Infix_Expr,
//     ^If_Expr,
//     ^Function_Literal,
//     ^Call_Expr,
//     ^Array_Literal,
//     ^Index_Expr,
//     ^Hash_Expr,
// }

Node :: struct {
    derived: Any_Node,
}

Expr :: struct {
    using expr_base: Node,
    // derived_expr: Any_Node,
}

Stmt :: struct {
    using expr_base: Node,
    // derived_stmt: Any_Node,
}

// Node_inst :: struct($T: typeid) {
//     using node: ^T,
// }
//
// Expr_inst :: struct($T: typeid) {
//     using expr: ^T,
// }
//
// Stmt_inst :: struct($T: typeid) {
//     using stmt: ^T,
// }

// Program

Program :: struct {
    using node: Node,
    statements: [dynamic]^Stmt,
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
    token: token.Token, // { token
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
    index: ^Expr,
}

Hash_Expr :: struct {
    using node: Expr,
    token: token.Token,
    pairs: map[^Expr]^Expr,
}

// always use this proc to turn the AST into a string
to_string :: proc(node: Node) -> string {
    switch v in node.derived {
        case ^Program: return program_string(v)
        case ^Expr_Stmt: return expr_stmt_string(v)
        case ^Let_Stmt: return let_stmt_string(v)
        case ^Return_Stmt: return return_stmt_string(v)
        case ^Block_Stmt: return block_stmt_string(v)
        case ^Ident: return ident_string(v)
        case ^Int_Literal: return int_literal_string(v)
        case ^String_Literal: return string_literal_string(v)
        case ^Bool_Literal: return bool_literal_string(v)
        case ^Prefix_Expr: return prefix_expr_string(v)
        case ^Infix_Expr: return infix_expr_string(v)
        case ^If_Expr: return if_expr_string(v)
        case ^Function_Literal: return function_expr_string(v)
        case ^Call_Expr: return call_expr_string(v)
        case ^Array_Literal: return array_expr_string(v)
        case ^Index_Expr: return index_expr_string(v)
        case ^Hash_Expr: return hash_expr_string(v)
        case:
            panic("unknown node type")
    }
}

program_string :: proc(p: ^Program) -> string {
    out: bytes.Buffer

    for stmt in p.statements {
        bytes.buffer_write(&out, transmute([]u8)to_string(stmt))
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

    bytes.buffer_write(&out, transmute([]u8)tok_literal(s.token))
    bytes.buffer_write(&out, transmute([]u8)string(" "))
    bytes.buffer_write(&out, transmute([]u8)to_string(s.name))
    bytes.buffer_write(&out, transmute([]u8)string(" = ")) // we must convert the untyped string literal to typed string

    if s.value != nil {
        bytes.buffer_write(&out, transmute([]u8)to_string(s.value))
    }

    bytes.buffer_write(&out, transmute([]u8)string(";"))

    return bytes.buffer_to_string(&out)
}

return_stmt_string :: proc(s: ^Return_Stmt) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)tok_literal(s.token))
    bytes.buffer_write(&out, transmute([]u8)string(" "))

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

    return bytes.buffer_to_string(&out)
}

int_literal_string :: proc(e: ^Int_Literal) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)e.token.literal)

    return bytes.buffer_to_string(&out)
}

string_literal_string :: proc(e: ^String_Literal) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)string("\""))
    bytes.buffer_write(&out, transmute([]u8)e.value)
    bytes.buffer_write(&out, transmute([]u8)string("\""))

    return bytes.buffer_to_string(&out)
}

bool_literal_string :: proc(e: ^Bool_Literal) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)e.token.literal)

    return bytes.buffer_to_string(&out)
}


prefix_expr_string :: proc(e: ^Prefix_Expr) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)string("("))
    bytes.buffer_write(&out, transmute([]u8)e.operator)
    bytes.buffer_write(&out, transmute([]u8)to_string(e.right))
    bytes.buffer_write(&out, transmute([]u8)string(")"))

    return bytes.buffer_to_string(&out)
}

infix_expr_string :: proc(e: ^Infix_Expr) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)string("("))
    bytes.buffer_write(&out, transmute([]u8)to_string(e.left))
    bytes.buffer_write(&out, transmute([]u8)string(" "))
    bytes.buffer_write(&out, transmute([]u8)e.operator)
    bytes.buffer_write(&out, transmute([]u8)string(" "))
    bytes.buffer_write(&out, transmute([]u8)to_string(e.right))
    bytes.buffer_write(&out, transmute([]u8)string(")"))

    return bytes.buffer_to_string(&out)
}

if_expr_string :: proc(e: ^If_Expr) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)e.token.literal)
    bytes.buffer_write(&out, transmute([]u8)string("("))
    bytes.buffer_write(&out, transmute([]u8)to_string(e.condition))
    bytes.buffer_write(&out, transmute([]u8)string(") "))
    bytes.buffer_write(&out, transmute([]u8)to_string(e.consequence))

    for alt in e.alternatives {
        bytes.buffer_write(&out, transmute([]u8)string("else if "))
        bytes.buffer_write(&out, transmute([]u8)if_expr_string(alt))
    }

    if e.default != nil {
        bytes.buffer_write(&out, transmute([]u8)string("else "))
        bytes.buffer_write(&out, transmute([]u8)to_string(e.default))
    }

    return bytes.buffer_to_string(&out)
}

function_expr_string :: proc(e: ^Function_Literal) -> string {
    out: bytes.Buffer

    params: [dynamic]string
    for p in e.params {
        append(&params, to_string(p))
    }
    res := strings.join(params[:], ", ", context.temp_allocator)
    defer delete(res, context.temp_allocator)
    bytes.buffer_write(&out, transmute([]u8)string(e.token.literal))
    bytes.buffer_write(&out, transmute([]u8)string("("))
    bytes.buffer_write(&out, transmute([]u8)res)
    bytes.buffer_write(&out, transmute([]u8)string(") "))
    bytes.buffer_write(&out, transmute([]u8)to_string(e.body))

    return bytes.buffer_to_string(&out)
}

call_expr_string :: proc(e: ^Call_Expr) -> string {
    out: bytes.Buffer

    args: [dynamic]string
    for a in e.args {
        append(&args, to_string(a))
    }
    res := strings.join(args[:], ", ", context.temp_allocator)
    defer delete(res, context.temp_allocator)

    bytes.buffer_write(&out, transmute([]u8)to_string(e.func))
    bytes.buffer_write(&out, transmute([]u8)string("("))
    bytes.buffer_write(&out, transmute([]u8)res)
    bytes.buffer_write(&out, transmute([]u8)string(")"))

    return bytes.buffer_to_string(&out)
}

array_expr_string :: proc(e: ^Array_Literal) -> string {
    out: bytes.Buffer

    elems: [dynamic]string
    for elem in e.elems {
        append(&elems, to_string(elem))
    }
    res := strings.join(elems[:], ", ", context.temp_allocator)
    defer delete(res, context.temp_allocator)

    bytes.buffer_write(&out, transmute([]u8)string("["))
    bytes.buffer_write(&out, transmute([]u8)res)
    bytes.buffer_write(&out, transmute([]u8)string("]"))

    return bytes.buffer_to_string(&out)
}

index_expr_string :: proc(e: ^Index_Expr) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)string("("))
    bytes.buffer_write(&out, transmute([]u8)to_string(e.left))
    bytes.buffer_write(&out, transmute([]u8)string("["))
    bytes.buffer_write(&out, transmute([]u8)to_string(e.index))
    bytes.buffer_write(&out, transmute([]u8)string("])"))

    return bytes.buffer_to_string(&out)
}

hash_expr_string :: proc(e: ^Hash_Expr) -> string {
    out: bytes.Buffer

    bytes.buffer_write(&out, transmute([]u8)string("{"))
    for k, v in e.pairs {
        key := to_string(k)
        value := to_string(v)
        bytes.buffer_write(&out, transmute([]u8)key)
        bytes.buffer_write(&out, transmute([]u8)string(":"))
        bytes.buffer_write(&out, transmute([]u8)value)
        bytes.buffer_write(&out, transmute([]u8)string(", "))
    }
    bytes.buffer_write(&out, transmute([]u8)string("}"))

    return bytes.buffer_to_string(&out)
}

tok_literal :: proc(t: token.Token) -> string {
    return t.literal
}
