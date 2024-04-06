package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
	"sort"
)

type EmittedInstruction struct {
    Opcode code.Opcode
    Pos int
}

type Bytecode struct {
    Instructions code.Instructions
    Constants []object.Object
}

type CompilationScope struct {
    instructions code.Instructions
    lastIns EmittedInstruction
    prevIns EmittedInstruction
}

type Compiler struct {
    Constants []object.Object
    symTable *SymTable
    scopes []CompilationScope
    scopeIndex int
}


func New_Compiler() *Compiler {
    mainScope := CompilationScope{
        instructions: code.Instructions{},
        lastIns: EmittedInstruction{},
        prevIns: EmittedInstruction{},
    }

    return &Compiler{
        Constants: []object.Object{},
        symTable: NewSymTable(),
        scopes: []CompilationScope{mainScope},
        scopeIndex: 0,
    }
}

func New_Compiler_With_States(constants []object.Object, symTable *SymTable) *Compiler {
    mainScope := CompilationScope{
        instructions: code.Instructions{},
        lastIns: EmittedInstruction{},
        prevIns: EmittedInstruction{},
    }

    return &Compiler{
        Constants: constants,
        symTable: symTable,
        scopes: []CompilationScope{mainScope},
        scopeIndex: 0,
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
        if c.lastInsIs(code.OpPop) { // remove the pop created by consequence
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

        c.changeOperand(jneInsPos, len(c.currentInstructions()))

        if node.Default == nil {
            c.emit(code.OpNull)
        } else {
            err := c.Compile(node.Default)
            if err != nil {
                return err
            }
            if c.lastInsIs(code.OpPop) { // remove the pop created by consequence
                c.removeLastPop()
            }
        }

        c.changeOperand(jmpPos, len(c.currentInstructions()))

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
    case *ast.HashLiteral:
        keys := []ast.Expression{}
        
        for k := range node.Pairs {
            keys = append(keys, k)
        }

        sort.Slice(keys, func(i, j int) bool {
            return keys[i].String() < keys[j].String()
        })

        for _, k := range keys {
            err := c.Compile(k)
            if err != nil {
                return err
            }

            err = c.Compile(node.Pairs[k])
            if err != nil {
                return err
            }
        }

        c.emit(code.OpHash, len(node.Pairs) * 2)
    case *ast.IndexExpression:
        err := c.Compile(node.Left)
        if err != nil {
            return err
        }

        err = c.Compile(node.Index)
        if err != nil {
            return err
        }

        c.emit(code.OpIndex)
    case *ast.FunctionLiteral:
        c.enterScope()

        err := c.Compile(node.Body)
        if err != nil {
            return err
        }

        if c.lastInsIs(code.OpPop) {
            c.replaceLastPopWithReturn()
        }

        if !c.lastInsIs(code.OpReturnValue) {
            c.emit(code.OpReturn)
        }

        instructions := c.leaveScope()
        compiledFun := &object.CompiledFunction{Instructions: instructions}
        c.emit(code.OpConstant, c.addConstant(compiledFun))
    case *ast.ReturnStatement:
        err := c.Compile(node.ReturnValue)
        if err != nil {
            return err
        }

        c.emit(code.OpReturnValue)
    }
    return nil
}

func (c *Compiler) Bytecode() *Bytecode {
    return &Bytecode{
        Instructions: c.currentInstructions(),
        Constants: c.Constants,
    }
}

func (c *Compiler) addConstant(obj object.Object) int {
    c.Constants = append(c.Constants, obj)
    return len(c.Constants) - 1
}

func (c *Compiler) addInstruction(ins []byte) int {
    pos_new_ins := len(c.currentInstructions())
    updatedInstructions := append(c.currentInstructions(), ins...)
    c.scopes[c.scopeIndex].instructions = updatedInstructions
    return pos_new_ins
}

// NOTE: the new Instruction should be the same width as the old one
func (c *Compiler) replaceInstruction(newIns []byte, pos int) {
    instructions := c.currentInstructions()
    for i := 0; i < len(newIns); i++ {
        instructions[pos + i] = newIns[i]
    }
}

func (c *Compiler) changeOperand(opPos int, operand int) {
    op := code.Opcode(c.currentInstructions()[opPos])
    newIns := code.Make(op, operand)
    c.replaceInstruction(newIns, opPos)
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
    ins := code.Make(op, operands...)
    pos := c.addInstruction(ins)

    c.setLastInstruction(op, pos)

    return pos
}

func (c *Compiler) currentInstructions() code.Instructions {
    return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
    last := EmittedInstruction{Opcode: op, Pos: pos}
    prev := c.scopes[c.scopeIndex].lastIns

    c.scopes[c.scopeIndex].prevIns = prev
    c.scopes[c.scopeIndex].lastIns = last
}

func (c *Compiler) lastInsIs(op code.Opcode) bool {
    if (len(c.currentInstructions()) == 0) {
        return false
    }

    return c.scopes[c.scopeIndex].lastIns.Opcode == op
}

func (c *Compiler) removeLastPop() {
    last := c.scopes[c.scopeIndex].lastIns
    prev := c.scopes[c.scopeIndex].prevIns

    old := c.currentInstructions()
    new := old[:last.Pos]

    c.scopes[c.scopeIndex].instructions = new
    c.scopes[c.scopeIndex].lastIns = prev
}

func (c *Compiler) replaceLastPopWithReturn() {
    last := c.scopes[c.scopeIndex].lastIns.Pos
    c.replaceInstruction(code.Make(code.OpReturnValue), last)

    c.scopes[c.scopeIndex].lastIns.Opcode = code.OpReturnValue
}

func (c *Compiler) enterScope() {
    scope := CompilationScope{
        instructions: code.Instructions{},
        lastIns: EmittedInstruction{},
        prevIns: EmittedInstruction{},
    }
    c.scopes = append(c.scopes, scope)
    c.scopeIndex++
}

func (c *Compiler) leaveScope() code.Instructions {
    instructions := c.currentInstructions()
    c.scopes = c.scopes[:len(c.scopes)-1]
    c.scopeIndex--

    return instructions
}
