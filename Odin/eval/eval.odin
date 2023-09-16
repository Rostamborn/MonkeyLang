package eval

import "../ast"
import "../object"

// reusable objects to avoid unnecessary memory allocation
NULL := object.new_obj(object.Null)
TRUE := object.new_bool_obj(true)
FALSE := object.new_bool_obj(false)

eval :: proc(node: ^ast.Node) -> object.Object {
    #partial switch v in node.derived {
        case ^ast.Int_Literal:

    }

    return NULL
}
