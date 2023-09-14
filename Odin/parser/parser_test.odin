package parser

import "core:fmt"
import "core:testing"
import "../lexer"
import "../ast"

@(test)
test_let_stmt :: proc(t: ^testing.T) {
    tests := []struct {
               input: string,
               expected_ident: string,
               expected_val: int,
           }{
               {"let x = 5;", "x", 5},
           }

    for tt in tests {
        lex := lexer.new_lexer(tt.input)
        p := new_parser(lex)
        program := parse_program(p)

        stmt := program.statements[0].derived.(^ast.Let_Stmt)
        name := stmt.name
        value := stmt.value.derived.(^ast.Int_Literal)

        if name.value != tt.expected_ident {
            testing.errorf(t, "expected %v, got %v", tt.expected_ident, stmt.name.value)
        }

        if value.value != tt.expected_val {
            testing.errorf(t, "expected %v, got %v", tt.expected_val, value.value)
        }
    }
}
