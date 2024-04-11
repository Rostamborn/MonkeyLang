package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

func (ins Instructions) String() string {
    var out bytes.Buffer

    i := 0
    for i < len(ins) {
        def, err := Lookup(ins[i])
        if err != nil {
            fmt.Fprintf(&out, "ERROR: %s\n", err)
            continue
        }

        operands, read := ReadOperands(def, ins[i+1:])

        fmt.Fprintf(&out, "%04d %s\n", i, ins.instructionsFormat(def, operands))

        i = i + 1 + read // +1 because of OpCode
    }

    return out.String()
}

func (ins Instructions) instructionsFormat(def *Definition, operands []int) string {
    operandCount := len(def.OperandWidths)

    switch operandCount {
    case 0:
        return def.Name
    case 1:
        return fmt.Sprintf("%s %d", def.Name, operands[0])
    }

    return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
 
type Opcode byte

// Enum
const (
    OpConstant Opcode = iota
    OpPop
    OpAdd
    OpSub
    OpMul
    OpDiv
    OpTrue
    OpFalse
    OpEqual
    OpNotEqual
    OpLessThan
    OpMinus
    OpBang
    OpJmp
    OpJNE // jump not equal
    OpNull
    OpSetGlobal
    OpGetGlobal
    OpSetLocal
    OpGetLocal
    OpArray
    OpHash
    OpIndex
    OpCall
    OpReturn
    OpReturnValue
)

type Definition struct {
    Name string
    OperandWidths []int
}

var definitions = map[Opcode]*Definition{
    OpConstant: {"OpConstant", []int{2}}, // 2-byte long single operand
    OpPop: {"OpPop", []int{}},
    OpAdd: {"OpAdd", []int{}},
    OpSub: {"OpSub", []int{}},
    OpMul: {"OpMul", []int{}},
    OpDiv: {"OpDiv", []int{}},
    OpTrue: {"OpTrue", []int{}},
    OpFalse: {"OpFalse", []int{}},
    OpEqual: {"OpEqual", []int{}},
    OpNotEqual: {"OpNotEqual", []int{}},
    OpLessThan: {"OpLessThan", []int{}},
    OpMinus: {"OpMinus", []int{}},
    OpBang: {"OpBang", []int{}},
    OpJmp: {"OpJmp", []int{2}},
    OpJNE: {"OpJNE", []int{2}},
    OpNull: {"OpNull", []int{}},
    OpSetGlobal: {"OpSetGlobal", []int{2}},
    OpGetGlobal: {"OpGetGlobal", []int{2}},
    OpSetLocal: {"OpSetLocal", []int{1}},
    OpGetLocal: {"OpGetLocal", []int{1}},
    OpArray: {"OpArray", []int{2}},
    OpHash: {"OpHash", []int{2}},
    OpIndex: {"OpIndex", []int{}},
    OpCall: {"OpCall", []int{1}},
    OpReturn: {"OpReturn", []int{}},
    OpReturnValue: {"OpReturnValue", []int{}},
}

func Lookup(op byte) (*Definition, error) {
    def, ok := definitions[Opcode(op)]
    if !ok {
        return nil, fmt.Errorf("opcode %d undefined", op)
    }

    return def, nil
}

// making bytecode
func Make(op Opcode, operands ...int) []byte {
    def, ok := definitions[op]
    if !ok {
        return []byte{}
    }

    instruction_len := 1 // we have opcode, so 1 byte

    for _, w := range def.OperandWidths {
        instruction_len += w
    }

    instruction := make([]byte, instruction_len)
    instruction[0] = byte(op)

    offset := 1
    for i, operand := range operands {
        width := def.OperandWidths[i]
        switch width {
        case 2:
            binary.BigEndian.PutUint16(instruction[offset:], uint16(operand)) // one operand which is 16bits
        case 1:
            instruction[offset] = byte(operand)
        }
        offset += width
    }

    return instruction
}

// Opposite of Make()
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
    operands := make([]int, len(def.OperandWidths))
    offset := 0

    for i, width := range def.OperandWidths {
        switch width {
        case 2:
            operands[i] = int(ReadUint16(ins[offset:]))
        case 1:
            operands[i] = int(ReadUint8(ins[offset:]))
        }

        offset += width
    }

    return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
    return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
    return uint8(ins[0])
}
