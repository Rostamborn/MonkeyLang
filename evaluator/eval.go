package evaluator

import (
    "monkey/ast"
    "monkey/object"
)

var (
    NULL = &object.Null{}
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        case *ast.Program:
            return evalStatements(node.Statements)
        case *ast.ExpressionStatement:
            return Eval(node.Expression)
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.Boolean:
            return nativeBoolToBooleanObject(node.Value)
        case *ast.PrefixExpression:
            right := Eval(node.Right)
            return evalPrefixExpression(node.Operator, right)
        case *ast.InfixExpression:
            left := Eval(node.Left)
            right := Eval(node.Right)
            return evalInfixExpression(node.Operator, left, right)
    }

    return NULL
}

func evalStatements(stmts []ast.Statement) object.Object {
    var res object.Object

    for _, stmt := range stmts {
        res = Eval(stmt)
    }

    return res
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
    if input {
        return TRUE
    }
    return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
    switch operator {
        case "!":
            return evalBangOperatorExpression(right)
        case "-":
            return evalMinusPrefixOperatorExpression(right)
        default:
            return NULL
    }
}

func evalBangOperatorExpression(right object.Object) object.Object {
    switch right {
    case TRUE:
        return FALSE
    case FALSE:
        return TRUE
    case NULL:
        return TRUE
    default:
        if right.Type() == object.INTEGER_OBJ && right.(*object.Integer).Value == 0 {
                return TRUE
        } else {
            return FALSE
        }
    }
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
    if right.Type() != object.INTEGER_OBJ {
        return NULL
    }

    value := right.(*object.Integer).Value
    return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ :
        return evalIntegerInfixExpression(operator, left, right)
    case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
        return evalBooleanInfixExpression(operator, left, right)
    default:
        return NULL
    }
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.Integer).Value
    rightVal := right.(*object.Integer).Value

    switch operator {
    case "+":
        return &object.Integer{Value: leftVal + rightVal}
    case "-":
        return &object.Integer{Value: leftVal - rightVal}
    case "*":
        return &object.Integer{Value: leftVal * rightVal}
    case "/":
        return &object.Integer{Value: leftVal / rightVal}
    case "<":
        return nativeBoolToBooleanObject(leftVal < rightVal) // pointer comparison
    case ">":
        return nativeBoolToBooleanObject(leftVal > rightVal) // pointer comparison
    case "==":
        return nativeBoolToBooleanObject(leftVal == rightVal) // pointer comparison
    case "!=":
        return nativeBoolToBooleanObject(leftVal != rightVal) // pointer comparison
    default:
        return NULL
    }
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
    leftVal := left.(*object.Boolean).Value
    rightVal := right.(*object.Boolean).Value

    switch operator {
    case "==":
        return nativeBoolToBooleanObject(leftVal == rightVal) // pointer comparison
    case "!=":
        return nativeBoolToBooleanObject(leftVal != rightVal) // pointer comparison
    default:
        return NULL
    }
}
