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
            return evalProgram(node)
        case *ast.BlockStatement:
            return evalBlockStatement(node)
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
        // case *ast.BlockStatement:
        //     return evalStatements(node.Statements)
        case *ast.IfExpression:
            return evalIfExpression(node)
        case *ast.ReturnStatement:
            value := Eval(node.ReturnValue)
            return &object.ReturnValue{Value: value}
    }

    return NULL
}

func evalProgram(program *ast.Program) object.Object {
    var res object.Object

    for _, stmt := range program.Statements {
        res = Eval(stmt)

        if returnValue, ok := res.(*object.ReturnValue); ok {
            return returnValue.Value
        }
    }

    return res
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
    var res object.Object

    for _, stmt := range block.Statements {
        res = Eval(stmt)

        if res != nil && res.Type() == object.RETURN_VALUE_OBJ {
            return res
        }
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

func evalIfExpression(ie *ast.IfExpression) object.Object {
    condition := Eval(ie.Condition)
    if isTruthy(condition) {
        return Eval(ie.Consequence)
    }
    for _, alt := range ie.Alternative {
        // altCondition := Eval(alt.Condition)
        // if isTruthy(altCondition) {
        //     return Eval(alt.Consequence)
        // }
        return evalIfExpression(alt)
    }
    if ie.Default != nil {
        return Eval(ie.Default)
    } else {
        return NULL
    }
}

func isTruthy(obj object.Object) bool {
    switch obj {
    case NULL:
        return false
    case TRUE:
        return true
    case FALSE:
        return false
    default:
        return true
    }
}
