package compiler

import (
    "fmt"
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
        c.emit(code.OpPop)
    case *ast.InfixExpression:
        fmt.Println("infix expression")
        if node.Operator == ">" {
            err := c.Compile(node.Right)
            if err != nil {
                return err
            }

            err = c.Compile(node.Left)
            if err != nil {
                return err
            }

            c.emit(code.OpLessThan)
            return nil
        }

        err := c.Compile(node.Left)
        if err != nil {
            return err
        }

        err = c.Compile(node.Right)
        if err != nil {
            return err
        }

        switch node.Operator {
        case "+":
            c.emit(code.OpAdd)
        case "-":
            c.emit(code.OpSub)
        case "*":
            c.emit(code.OpMul)
        case "/":
            c.emit(code.OpDiv)
        case "<":
            c.emit(code.OpLessThan)
            fmt.Println("lessthan emitted")
        case "==":
            c.emit(code.OpEqual)
        case "!=":
            c.emit(code.OpNotEqual)
        default:
            return fmt.Errorf("unknown operator: %s", node.Operator)
        }
    case *ast.IntegerLiteral:
        integer := &object.Integer{Value: node.Value}
        c.emit(code.OpConstant, c.addConstant(integer)) // the index in constant pool
    case *ast.Boolean:
        if node.Value {
            c.emit(code.OpTrue)
        } else {
            c.emit(code.OpFalse)
        }
    }                                                   // is the constant identifier
    fmt.Println(c.Instructions)
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
