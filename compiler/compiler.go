package compiler

import (
    "fmt"
    "monkey/ast"
    "monkey/object"
    "monkey/code"
)

type EmittedInstruction struct {
    Opcode code.Opcode
    Pos int
}

type Bytecode struct {
    Instructions code.Instructions
    Constants []object.Object
}

type Compiler struct {
    Instructions code.Instructions
    Constants []object.Object
    lastIns EmittedInstruction
    prevIns EmittedInstruction
    symTable *SymTable
}


func New_Compiler() *Compiler {
    return &Compiler{
        Instructions: code.Instructions{},
        Constants: []object.Object{},
        lastIns: EmittedInstruction{},
        prevIns: EmittedInstruction{},
        symTable: NewSymTable(),
    }
}

func New_Compiler_With_States(constants []object.Object, symTable *SymTable) *Compiler {
    return &Compiler{
        Instructions: code.Instructions{},
        Constants: constants,
        lastIns: EmittedInstruction{},
        prevIns: EmittedInstruction{},
        symTable: symTable,
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
    case *ast.BlockStatement:
        for _, s := range node.Statements {
            err := c.Compile(s)
            if err != nil {
                return err
            }
        }
    case *ast.LetStatement:
        err := c.Compile(node.Value)
        if err != nil {
            return err
        }
        
        symbol := c.symTable.Define(node.Name.Value)
        c.emit(code.OpSetGlobal, symbol.Index)

    case *ast.ExpressionStatement:
        err := c.Compile(node.Expression)
        if err != nil {
            return err
        }
        c.emit(code.OpPop)

    case *ast.InfixExpression:
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
        case "==":
            c.emit(code.OpEqual)
        case "!=":
            c.emit(code.OpNotEqual)
        default:
            return fmt.Errorf("unknown operator: %s", node.Operator)
        }
    case *ast.PrefixExpression:
        err := c.Compile(node.Right)
        if err != nil {
            return err
        }
        
        switch node.Operator {
        case "!":
            c.emit(code.OpBang)
        case "-":
            c.emit(code.OpMinus)
        default:
            return fmt.Errorf("unkown operator %s", node.Operator)
        }
    case *ast.IfExpression:
        err := c.Compile(node.Condition)
        if err != nil {
            return err
        }

        jneInsPos := c.emit(code.OpJNE, 6969)

        err = c.Compile(node.Consequence)
        if err != nil {
            return err
        }
        if c.lastIns.Opcode == code.OpPop { // remove the pop created by consequence
            c.removeLastPop()
        }

        // TODO: this is deeply wrong(meaning the ast is generated wrongly)
        // and hence the generated code is wrong. fix them all
        for _, alt := range node.Alternative {
            err := c.Compile(alt)
            if err != nil {
                return err
            }
        }

        jmpPos := c.emit(code.OpJmp, 6969)

        c.changeOperand(jneInsPos, len(c.Instructions))

        if node.Default == nil {
            c.emit(code.OpNull)
        } else {
            err := c.Compile(node.Default)
            if err != nil {
                return err
            }
            if c.lastIns.Opcode == code.OpPop { // remove the pop created by consequence
                c.removeLastPop()
            }
        }

        c.changeOperand(jmpPos, len(c.Instructions))

    case *ast.IntegerLiteral:
        integer := &object.Integer{Value: node.Value}
        c.emit(code.OpConstant, c.addConstant(integer)) // the index in constant pool
    case *ast.Boolean:
        if node.Value {
            c.emit(code.OpTrue)
        } else {
            c.emit(code.OpFalse)
        }
    case *ast.StringLiteral:
        str := &object.String{Value: node.Value}
        c.emit(code.OpConstant, c.addConstant(str))
    case *ast.Identifier:
        symbol, ok := c.symTable.Resolve(node.Value)
        if !ok {
            return fmt.Errorf("undefined variable %s", node.Value)
        }

        c.emit(code.OpGetGlobal, symbol.Index)
    case *ast.ArrayLiteral:
        for _, elem := range node.Elements {
            err := c.Compile(elem)
            if err != nil {
                return err
            }
        }

        c.emit(code.OpArray, len(node.Elements))
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

// NOTE: the new Instruction should be the same width as the old one
func (c *Compiler) replaceInstruction(newIns []byte, pos int) {
    for i := 0; i < len(newIns); i++ {
        c.Instructions[pos + i] = newIns[i]
    }
}

func (c *Compiler) changeOperand(opPos int, operand int) {
    op := code.Opcode(c.Instructions[opPos])
    newIns := code.Make(op, operand)
    c.replaceInstruction(newIns, opPos)
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
    ins := code.Make(op, operands...)
    pos := c.addInstruction(ins)

    c.setLastInstruction(op, pos)

    return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
    last := EmittedInstruction{Opcode: op, Pos: pos}
    c.prevIns = c.lastIns
    c.lastIns = last
}

func (c *Compiler) removeLastPop() {
    c.Instructions = c.Instructions[:c.lastIns.Pos]
    c.lastIns = c.prevIns
}
