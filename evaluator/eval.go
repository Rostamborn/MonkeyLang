package evaluator

import (
    "fmt"
    "monkey/ast"
    "monkey/object"
)

var (
    NULL = &object.Null{}
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {
        case *ast.Program:
            return evalProgram(node, env)
        case *ast.BlockStatement:
            return evalBlockStatement(node, env)
        case *ast.ExpressionStatement:
            return Eval(node.Expression, env)
        case *ast.LetStatement:
            value := Eval(node.Value,env)
            if isError(value) {
                return value
            }
            env.Set(node.Name.Value, value)
            return value
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.StringLiteral:
            return &object.String{Value: node.Value}
        case *ast.Boolean:
            return nativeBoolToBooleanObject(node.Value)
        case *ast.PrefixExpression:
            right := Eval(node.Right, env)
            if isError(right) {
                return right
            }
            return evalPrefixExpression(node.Operator, right)
        case *ast.InfixExpression:
            left := Eval(node.Left, env)
            if isError(left) {
                return left
            }

            right := Eval(node.Right,env)
            if isError(right) {
                return right
            }
            return evalInfixExpression(node.Operator, left, right)
        case *ast.IfExpression:
            return evalIfExpression(node, env)
        case *ast.ReturnStatement:
            value := Eval(node.ReturnValue, env)
            if isError(value) {
                return value
            }
            return &object.ReturnValue{Value: value}
        case *ast.Identifier:
            return evalIdentifier(node, env)
        case *ast.FunctionLiteral:
            params := node.Parameters
            body := node.Body
            return &object.Function{Parameters: params, Env: env, Body: body}
        case *ast.CallExpression:
            function := Eval(node.Function, env)
            if isError(function) {
                return function
            }

            args := evalExpressions(node.Arguments, env)
            if len(args) == 1 && isError(args[0]) { // wtf is this
                return args[0]
            }

            return applyFunction(function, args)
        case *ast.ArrayLiteral:
            elements := evalExpressions(node.Elements, env)
            if len(elements) == 1 && isError(elements[0]) {
                return elements[0]
            }

            return &object.Array{Elements: elements}
        case *ast.IndexExpression:
            left := Eval(node.Left, env)
            if isError(left) {
                return left
            }

            index := Eval(node.Index, env)
            if isError(index) {
                return index
            }

            return evalIndexExpression(left, index)
        case *ast.HashLiteral:
            return evalHashLiteral(node, env)
    }

    return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
    var result object.Object

    for _, stmt := range program.Statements {
        result = Eval(stmt, env)

        switch res := result.(type) {
        case *object.ReturnValue:
            return res.Value
        case *object.Error:
            return res
        }
    }

    return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
    var res object.Object

    for _, stmt := range block.Statements {
        res = Eval(stmt, env)

        if res != nil {
            resType := res.Type()
            if resType == object.RETURN_VALUE_OBJ || resType == object.ERROR_OBJ {
                return res
            }
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
            return newError("unknown operator: %s%s", operator, right.Type())
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
        return newError("unknown operator: -%s", right.Type())
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
    case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
        return evalStringInfixExpression(operator, left, right)
    case left.Type() != right.Type():
        return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
    default:
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
    if operator != "+" {
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }

    leftVal := left.(*object.String).Value
    rightVal := right.(*object.String).Value

    return &object.String{Value: leftVal + rightVal}
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
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
    condition := Eval(ie.Condition, env)
    if isError(condition) {
        return condition
    }
    if isTruthy(condition) {
        cons := Eval(ie.Consequence, env)
        if isError(cons) {
            return cons
        }
        return cons
    }
    for _, alt := range ie.Alternative {
        return evalIfExpression(alt, env)
    }
    if ie.Default != nil {
        defEvalueated := Eval(ie.Default, env)
        if isError(defEvalueated) {
            return defEvalueated
        }
        return defEvalueated
    } else {
        return NULL
    }
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    if val, ok := env.Get(node.Value); ok {
        return val
    }

    if builtin, ok := builtins[node.Value]; ok {
        return builtin
    }

    return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
    var result []object.Object

    for _, e := range exps {
        evaluated := Eval(e, env)
        if isError(evaluated) {
            return []object.Object{evaluated}
        }

        result = append(result, evaluated)
    }

    return result
}

func evalIndexExpression(left, index object.Object) object.Object {
    switch {
    case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
        return evalArrayIndexExpression(left, index)
    case left.Type() == object.HASHMAP_OBJ:
        return evalHashMapIndexExpression(left, index)
    default:
        return newError("index operator not supported: %s", left.Type())
    }
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
    arrayObject := array.(*object.Array)
    idx := index.(*object.Integer).Value
    maximum := int64(len(arrayObject.Elements) - 1)

    if idx < 0 || idx > maximum {
        return NULL
    }

    return arrayObject.Elements[idx]
}

func evalHashMapIndexExpression(hash, index object.Object) object.Object {
    hashMapObject := hash.(*object.HashMap)

    key, ok := index.(object.Hashable)
    if !ok {
        return newError("unusable as hash key: %s", index.Type())
    }

    pair, ok := hashMapObject.Pairs[key.HashKey()]
    if !ok {
        return NULL
    }

    return pair.Value
}

func evalHashLiteral(hashLitral *ast.HashLiteral, env *object.Environment) object.Object {
    pairs := make(map[object.HashKey]object.HashPair)

    for keyNode, valueNode := range hashLitral.Pairs {
        key := Eval(keyNode, env)
        if isError(key) {
            return key
        }

        hashKey, ok := key.(object.Hashable)
        if !ok {
            return newError("unusable as hash key: %s", key.Type())
        }

        value := Eval(valueNode, env)
        if isError(value) {
            return value
        }

        hashed := hashKey.HashKey()

        pairs[hashed] = object.HashPair{Key: key, Value: value}
    }

    return &object.HashMap{Pairs: pairs}
}

func applyFunction(function object.Object, args []object.Object) object.Object {
    switch fn := function.(type) {
    case *object.Function:
        extendedEnv := extendFunctionEnv(fn, args)
        evaluated := Eval(fn.Body, extendedEnv)
        return unwrapReturnValue(evaluated)
    case *object.Builtin:
        return fn.Fn(args...)
    default:
        return newError("not a function: %s", function.Type())
    }
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
    env := object.NewEnclosedEnvironment(fn.Env)

    for index, param := range fn.Parameters {
        env.Set(param.Value, args[index])
    }

    return env
}

func unwrapReturnValue(obj object.Object) object.Object {
    if returnValue, ok := obj.(*object.ReturnValue); ok {
        return returnValue.Value
    }

    return obj
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

func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
    if obj != nil {
        return obj.Type() == object.ERROR_OBJ
    }
    return false
}
