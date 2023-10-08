package compiler

import (
    "monkey/ast"
    "monkey/object"
    "monkey/code"
)

type Bytecode struct {
    Instructions code.Instructions
    Constants []object.Object
}

type Compiler struct {
    Instructions code.Instructions
    Constants []object.Object
}

func New_Compiler() *Compiler {
    return &Compiler{
        Instructions: code.Instructions{},
        Constants: []object.Object{},
    }
}

func (c *Compiler) Compile(node ast.Node) error {
    return nil
}

func (c *Compiler) Bytecode() *Bytecode {
    return &Bytecode{
        Instructions: c.Instructions,
        Constants: c.Constants,
    }
}
