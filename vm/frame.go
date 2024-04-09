package vm

import (
    "monkey/code"
    "monkey/object"
)

type Frame struct {
    fun *object.CompiledFunction
    ip int
    basePtr int // for storing the start of the calling frame on the stack
}

func New_Frame(fun *object.CompiledFunction, basePtr int) *Frame {
    return &Frame{fun: fun, ip: -1, basePtr: basePtr}
}

func (f *Frame) Instructions() code.Instructions {
    return f.fun.Instructions
}
