package vm

import (
    "monkey/code"
    "monkey/object"
)

type Frame struct {
    fun *object.CompiledFunction
    ip int
}

func New_Frame(fun *object.CompiledFunction) *Frame {
    return &Frame{fun: fun, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
    return f.fun.Instructions
}
