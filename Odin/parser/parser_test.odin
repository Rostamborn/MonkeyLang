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

@(test)
test_op_prec_parsing :: proc(t: ^testing.T) {
    tests := []struct {
        input: string,
        expected: string,
    }{
        // {"-a * b", "((-a) * b)"},
        // {"!-a", "(!(-a))"},
        // {"a + b + c", "((a + b) + c)"},
        // {"a + b - c", "((a + b) - c)"},
        // {"a * b * c", "((a * b) * c)"},
        // {"a * b / c", "((a * b) / c)"},
        // {"a + b / c", "(a + (b / c))"},
        // // {"(a + b) * c", ""},
        // {"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
        // {"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
        // {"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
        // {"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
        // {"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
        // {"true", "true"},
        // {"false", "false"},
        // {"3 > 5 == false", "((3 > 5) == false)"},
        // {"3 < 5 == true", "((3 < 5) == true)"},
        // {"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
        // {"(5 + 5) * 2", "((5 + 5) * 2)"},
        // {"2 / (5 + 5)", "(2 / (5 + 5))"},
        // {"-(5 + 5)", "(-(5 + 5))"},
        // {"!(true == true)", "(!(true == true))"},
        // {"a + add(b * c) + d", "((a + add((b * c))) + d)"},
        // {"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
        // {"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
        {"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
        // {"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
    }

    for tt in tests {
        lex := lexer.new_lexer(tt.input)
        p := new_parser(lex)
        program := parse_program(p)
        check_parser_errors(t, p)

        actual := ast.to_string(program)
        if actual != tt.expected {
            testing.errorf(t, "expected=%q, got=%q", tt.expected, actual) } 
    }
}

check_parser_errors :: proc(t: ^testing.T, p: ^Parser) {
    errors := p.errors
    if len(errors) == 0 {
        return
    }

    testing.errorf(t, "parser has %d errors", len(errors))
    for msg in errors {
        testing.errorf(t, "parser error: %q", msg)
    }
    testing.fail(t)
}
