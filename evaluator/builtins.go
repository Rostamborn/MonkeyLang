package evaluator

import "monkey/object"

var builtins = map[string]*object.Builtin{
    "len": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return newError("wrong number of arguments. got=%d, want=1", len(args))
            }

            switch arg := args[0].(type) {
            case *object.String:
                return &object.Integer{Value: int64(len(arg.Value))}
            case *object.Array:
                return &object.Integer{Value: int64(len(arg.Elements))}
            default:
                return newError("argument to `len` not supported, got %s", args[0].Type())
            }
        },
    },
    "ordered_remove": { // it is done in O(n) time
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 2 {
                return newError("wrong number of arguments. got=%d, want=2", len(args))
            }

            if args[0].Type() != object.ARRAY_OBJ {
                return newError("first argument to `ordered_remove` must be ARRAY, got %s", args[0].Type())
            }

            if args[1].Type() != object.INTEGER_OBJ {
                return newError("second argument to `ordered_remove` must be INTEGER, got %s", args[1].Type())
            }

            arr := args[0].(*object.Array)
            idx := args[1].(*object.Integer).Value

            if idx >= int64(len(arr.Elements)) || idx < 0 {
                return newError("index out of bounds")
            }
            
            arr.Elements = append(arr.Elements[:idx], arr.Elements[idx+1:]...)

            return arr
        },
    },
    "unordered_remove": { // it is done in O(1) time
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 2 {
                return newError("wrong number of arguments. got=%d, want=2", len(args))
            }

            if args[0].Type() != object.ARRAY_OBJ {
                return newError("first argument to `ordered_remove` must be ARRAY, got %s", args[0].Type())
            }

            if args[1].Type() != object.INTEGER_OBJ {
                return newError("second argument to `ordered_remove` must be INTEGER, got %s", args[1].Type())
            }

            arr := args[0].(*object.Array)
            idx := args[1].(*object.Integer).Value

            if idx >= int64(len(arr.Elements)) || idx < 0 {
                return newError("index out of bounds")
            }
            
            arr.Elements[idx] = arr.Elements[len(arr.Elements)-1]
            arr.Elements = arr.Elements[:len(arr.Elements)-1]

            return arr
        },
    },
    "append": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 2 {
                return newError("wrong number of arguments. got=%d, want=2", len(args))
            }

            if args[0].Type() != object.ARRAY_OBJ {
                return newError("first argument to `append` must be ARRAY, got %s", args[0].Type())
            }

            arr := args[0].(*object.Array)
            arr.Elements = append(arr.Elements, args[1])

            return arr
        },
    },
    "insert": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 3 {
                return newError("wrong number of arguments. got=%d, want=3", len(args))
            }

            if args[0].Type() != object.ARRAY_OBJ {
                return newError("first argument to `insert` must be ARRAY, got %s", args[0].Type())
            }

            if args[1].Type() != object.INTEGER_OBJ {
                return newError("second argument to `insert` must be INTEGER, got %s", args[1].Type())
            }

            arr := args[0].(*object.Array)
            idx := args[1].(*object.Integer).Value

            if idx >= int64(len(arr.Elements)) || idx < 0 {
                return newError("index out of bounds")
            }

            arr.Elements = append(arr.Elements[:idx], append([]object.Object{args[2]}, arr.Elements[idx:]...)...)

            return arr
        },
    },
}
