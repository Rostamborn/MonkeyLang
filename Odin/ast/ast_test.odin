package ast

import "core:testing"
import "core:fmt"
import "../token"

@(test)
test_string :: proc(t: ^testing.T) {

    program := new_node(Program) 
    name_ident := new_node(Ident)
    name_ident.token = token.Token{type= token.IDENT, literal= "myVar"}
    name_ident.value = "myVar"
    value_ident := new_node(Ident)
    value_ident.token = token.Token{type= token.IDENT, literal= "anotherVar"}
    value_ident.value = "anotherVar"
    let := new_node(Let_Stmt)
    let.token = token.Token{type= token.LET, literal= "let"}
    let.name = name_ident
    let.value = value_ident
    program.statements = []^Stmt{let}


    if to_string(program) != "let myVar = anotherVar;" {
        testing.errorf(t, "program.String() wrong. got=%q", to_string(program))
    }
}
