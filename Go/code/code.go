package code

import (
    "fmt"
    "encoding/binary"
)

type Instructions []byte
 
type Opcode byte

const (
    OpConstant Opcode = iota
)

type Definition struct {
    Name string
    OperandWidths []int
}

var definitions = map[Opcode]*Definition{
    OpConstant: {"OpConstant", []int{2}},
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

    instruction_len := 1 // we have opcode

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
        }
        offset += width
    }

    return instruction
}
