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
    switch node := node.(type) {
    case *ast.Program:
        for _, s := range node.Statements {
            err := c.Compile(s)
            if err != nil {
                return err
            }
        }
    case *ast.ExpressionStatement:
        err := c.Compile(node.Expression)
        if err != nil {
            return err
        }
    case *ast.InfixExpression:
        err := c.Compile(node.Left)
        if err != nil {
            return err
        }

        err = c.Compile(node.Right)
        if err != nil {
            return err
        }
    case *ast.IntegerLiteral:
        integer := &object.Integer{Value: node.Value}
        c.emit(code.OpConstant, c.addConstant(integer)) // the index in constant pool
    }                                                   // is the constant identifier
                                                        // adn the VM knows to load what
    return nil
}

func (c *Compiler) Bytecode() *Bytecode {
    return &Bytecode{
        Instructions: c.Instructions,
        Constants: c.Constants,
    }
}

func (c *Compiler) addConstant(obj object.Object) int {
    c.Constants = append(c.Constants, obj)
    return len(c.Constants) - 1
}

func (c *Compiler) addInstruction(ins []byte) int {
    pos_new_ins := len(c.Instructions)
    c.Instructions = append(c.Instructions, ins...)
    return pos_new_ins
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
    ins := code.Make(op, operands...)
    pos := c.addInstruction(ins)
    return pos
}
