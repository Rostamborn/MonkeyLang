package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

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

func (vm *VM) StackTop() object.Object {
    if vm.sp == 0 {
        return nil
    }

    return vm.stack[vm.sp-1]
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
        case code.OpAdd:
            right := vm.pop()
            left := vm.pop()
            if left != nil && right != nil {
                left_val := left.(*object.Integer)
                right_val := right.(*object.Integer)
                vm.push(&object.Integer{Value: left_val.Value + right_val.Value})
            } else {
                return fmt.Errorf("unsupported types for operation %d : %s %s", code.OpAdd, left.Type(), right.Type())
            }
        }
    }

    return nil;
}
