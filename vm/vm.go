package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

const StackSize = 2048
const GlobalSize = 65536
const MaxFrames = 1024

type VM struct {
    constants []object.Object
    stack []object.Object
    sp int // stack pointer which points to the
           // place after the top of the stack
    globals []object.Object
    frames []*Frame
    framesIndex int
}

func New_VM(bytecode *compiler.Bytecode) *VM {
    mainFun := &object.CompiledFunction{Instructions: bytecode.Instructions}
    mainFrame := New_Frame(mainFun)

    frames := make([]*Frame, MaxFrames)
    frames[0] = mainFrame

    return &VM{
        constants: bytecode.Constants,
        stack: make([]object.Object, StackSize),
        sp: 0,
        globals: make([]object.Object, GlobalSize),
        frames: frames,
        framesIndex: 1,
    }
}

func New_VM_With_Global_Store(bytecode *compiler.Bytecode, s []object.Object) *VM {
    vm := New_VM(bytecode)
    vm.globals = s
    return vm
}

func (vm *VM) LastPopped() object.Object {
    return vm.stack[vm.sp]
}

func (vm *VM) push(obj object.Object) error {
    if vm.sp >= StackSize { // StackSize is 2048 and we set limit to 2047 as we start from 0
        return fmt.Errorf("Stack Overflow")
    }
    vm.stack[vm.sp] = obj
    vm.sp++;

    return nil
}

func (vm *VM) pop() object.Object {
    if vm.sp == 0 {
        return Null
    }

    obj := vm.stack[vm.sp-1]
    vm.sp--

    return obj
}

func (vm *VM) Run() error {
    var ip int
    var ins code.Instructions
    var op code.Opcode

    for vm.currFrame().ip < (len(vm.currFrame().Instructions()) - 1) {
        fmt.Printf("ip: %d\n", ip)
        vm.currFrame().ip++

        ip = vm.currFrame().ip
        ins = vm.currFrame().Instructions()
        op = code.Opcode(ins[ip])
        fmt.Println(ins)

        switch op {
        case code.OpConstant:
            const_index := code.ReadUint16(ins[ip+1:])
            vm.currFrame().ip += 2
            err := vm.push(vm.constants[const_index])
            if err != nil {
                return err
            }
        case code.OpJmp:
            pos := int(code.ReadUint16(ins[ip+1:])) // read the operand
            vm.currFrame().ip = pos - 1 // the for loop will increment to pos by itself
        case code.OpJNE:
            pos := int(code.ReadUint16(ins[ip+1:])) // read the operand
            vm.currFrame().ip += 2 // 2 bytes because of the operand of length 2 bytes

            condition := vm.pop()
            if !isTruthy(condition) {
                vm.currFrame().ip = pos - 1
            }
        case code.OpSetGlobal:
            globalIndex := code.ReadUint16(ins[ip+1:])
            vm.currFrame().ip += 2

            vm.globals[globalIndex] = vm.pop()

        case code.OpGetGlobal:
            globalIndex := code.ReadUint16(ins[ip+1:])
            vm.currFrame().ip += 2

            err := vm.push(vm.globals[globalIndex])
            if err != nil {
                return err
            }
        case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
            err := vm.executeBinaryOperation(op)
            if err != nil {
                return err
            }
        case code.OpFalse:
            err := vm.push(False)
            if err != nil {
                return err
            }
        case code.OpTrue:
            err := vm.push(True)
            if err != nil {
                return err
            }
        case code.OpEqual, code.OpLessThan, code.OpNotEqual:
            err := vm.executeComparison(op)
            if err != nil {
                return err
            }
        case code.OpBang:
            err := vm.executeBangOperator()
            if err != nil {
                return err
            }
        case code.OpMinus:
            err := vm.executeMinusOperator()
            if err != nil {
                return err
            }
        case code.OpArray:
            numElements := int(code.ReadUint16(ins[ip+1:]))
            vm.currFrame().ip += 2

            array := vm.buildArray(vm.sp - numElements, vm.sp)
            vm.sp -= numElements

            err := vm.push(array)
            if err != nil {
                return err
            }
        case code.OpHash:
            numElements := int(code.ReadUint16(ins[ip+1:]))
            vm.currFrame().ip += 2

            hashmap, err := vm.buildHashMap(vm.sp - numElements, vm.sp)
            if err != nil {
                return err
            }

            vm.sp -= numElements

            err = vm.push(hashmap)
        case code.OpIndex:
            index := vm.pop()
            left := vm.pop()
            
            err := vm.executeIndexExpression(left, index)
            if err != nil {
                return err
            }
        case code.OpCall:
            fn, ok := vm.stack[vm.sp - 1].(*object.CompiledFunction)
            if !ok {
                return fmt.Errorf("calling non-function")
            }
            frame := New_Frame(fn)
            vm.pushFrame(frame)
        case code.OpReturnValue:
            returnValue := vm.pop()
            vm.popFrame()
            vm.pop() // popping CompiledFunction

            err := vm.push(returnValue)
            if err != nil {
                return err
            }
        case code.OpReturn:
            vm.popFrame()
            vm.pop() // popping CompiledFunction

            err := vm.push(Null)
            if err != nil {
                return err
            }
        case code.OpPop:
            vm.pop()
        case code.OpNull:
            err := vm.push(Null)
            if err != nil {
                return err
            }
        }
            
    }

    return nil;
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
    right := vm.pop()
    left := vm.pop()

    left_type := left.Type()
    right_type := right.Type()

    switch {
    case left_type == object.INTEGER_OBJ && right_type == object.INTEGER_OBJ:
        return vm.executeBinaryIntegerOperation(op, left, right)
    case left_type == object.STRING_OBJ && right_type == object.STRING_OBJ:
        return vm.executeBinaryStringOperation(op, left, right)
    }

    return fmt.Errorf("unsupported types for binary operation: %s %s", left_type, right_type)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
    left_val := left.(*object.Integer).Value
    right_val := right.(*object.Integer).Value
    var result int64

    switch op {
    case code.OpAdd:
        result = left_val + right_val
    case code.OpSub:
        result = left_val - right_val
    case code.OpMul:
        result = left_val * right_val
    case code.OpDiv:
        result = left_val / right_val
    default:
        return fmt.Errorf("unkown integer operator: %d", op)
    }

    return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
    if op != code.OpAdd {
        return fmt.Errorf("unknown string operator: %d", op)
    }

    left_val := left.(*object.String).Value
    right_val := right.(*object.String).Value

    return vm.push(&object.String{Value: left_val + right_val})
}

func (vm *VM) executeComparison(op code.Opcode) error {
    right := vm.pop()
    left := vm.pop()

    if left.Type() == object.INTEGER_OBJ || right.Type() == object.INTEGER_OBJ {
        return vm.executeIntegerComparison(op, left, right)
    }

    switch op {
    case code.OpEqual:
        return vm.push(nativeBoolToBooleanObject(right == left))
    case code.OpNotEqual:
        return vm.push(nativeBoolToBooleanObject(right != left))
    default:
        return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
    }
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
    left_val := left.(*object.Integer).Value
    right_val := right.(*object.Integer).Value

    switch op {
    case code.OpEqual:
        return vm.push(nativeBoolToBooleanObject(left_val == right_val))
    case code.OpNotEqual:
        return vm.push(nativeBoolToBooleanObject(left_val != right_val))
    case code.OpLessThan:
        return vm.push(nativeBoolToBooleanObject(left_val < right_val))
    default:
        return fmt.Errorf("unknown operator: %d", op)
    }
}

func (vm *VM) executeBangOperator() error {
    operand := vm.pop()

    switch operand {
    case True:
        return vm.push(False)
    case False:
        return vm.push(True)
    case Null:
        return vm.push(True)
    default:
        return vm.push(False)
    }
}

func (vm *VM) executeMinusOperator() error {
    operand := vm.pop()

    if operand.Type() != object.INTEGER_OBJ {
        return fmt.Errorf("unsupported type for negation: %s", operand.Type())
    }

    val := operand.(*object.Integer).Value
    return vm.push(&object.Integer{Value: -val})
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
    switch {
    case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
        return vm.executeArrayIndex(left, index)
    case left.Type() == object.HASHMAP_OBJ:
        return vm.executeHashIndex(left, index)
    default:
        return fmt.Errorf("index operator not supported: %s", left.Type())
    }
}

func (vm *VM) executeArrayIndex(left, index object.Object) error {
    arrayObj := left.(*object.Array)
    i := index.(*object.Integer).Value
    boundry := int64(len(arrayObj.Elements) - 1)

    if i < 0 || i > boundry {
        return vm.push(Null)
    }

    return vm.push(arrayObj.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index object.Object) error {
    hashMap := hash.(*object.HashMap)

    key, ok := index.(object.Hashable)
    if !ok {
        return fmt.Errorf("unusable as hash key: %s", index.Type())
    }

    pair, ok := hashMap.Pairs[key.HashKey()]
    if !ok {
        return vm.push(Null)
    }

    return vm.push(pair.Value)
}

func (vm *VM) buildArray(start , end int) object.Object {
    elements := make([]object.Object, end - start)

    for i := start; i < end; i++ {
        elements[i - start] = vm.stack[i]
    }

    return &object.Array{Elements: elements}
}

func (vm *VM) buildHashMap(start, end int) (object.Object, error) {
    hashedPairs := make(map[object.HashKey]object.HashPair)

    for i := start; i < end; i += 2 {
        key := vm.stack[i]
        value := vm.stack[i+1]

        pair := object.HashPair{Key: key, Value: value}

        hashKey, ok := key.(object.Hashable)
        if !ok {
            return nil, fmt.Errorf("unsuable as hash key: %s", key.Type())
        }

        hashedPairs[hashKey.HashKey()] = pair
    }

    return &object.HashMap{Pairs: hashedPairs}, nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
    if input {
        return True
    }
    return False
}

func isTruthy(obj object.Object) bool {
    switch obj := obj.(type) {
    case *object.Boolean:
        return obj.Value
    case *object.Null:
        return false
    default:
        return true
    }
}

func (vm *VM) currFrame() *Frame {
    return vm.frames[vm.framesIndex - 1]
}

func (vm *VM) pushFrame(f *Frame) {
    vm.frames[vm.framesIndex] = f
    vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
    vm.framesIndex--
    return vm.frames[vm.framesIndex]
}
