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
        }
    }

    return nil;
}
