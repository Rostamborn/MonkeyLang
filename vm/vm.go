package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

const StackSize = 2048

type VM struct {
    instructions code.Instructions
    constants []object.Object
    stack []object.Object
    sp int // stack pointer which points to the
           // place after the top of the stack
}

func New_VM(bytecode *compiler.Bytecode) *VM {
    return &VM{
        instructions: bytecode.Instructions,
        constants: bytecode.Constants,
        stack: make([]object.Object, StackSize),
        sp: 0,
    }
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
        return nil
    }

    obj := vm.stack[vm.sp-1]
    vm.sp--

    return obj
}

func (vm *VM) Run() error {
    for ip := 0; ip < len(vm.instructions); ip++ {
        op := code.Opcode(vm.instructions[ip])

        switch op {
        case code.OpConstant:
            const_index := code.ReadUint16(vm.instructions[ip+1:])
            ip += 2
            err := vm.push(vm.constants[const_index])
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
        case code.OpPop:
            vm.pop()
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
    if input {
        return True
    }
    return False
}
