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

        i = i + 1 + read
    }

    return out.String()
}

func (ins Instructions) instructionsFormat(def *Definition, operands []int) string {
    operandCount := len(def.OperandWidths)
    if operandCount == 0 {
        return def.Name
    }

    switch operandCount {
    case 1:
        return fmt.Sprintf("%s %d", def.Name, operands[0])
    }

    return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
 
type Opcode byte

// Enum
const (
    OpConstant Opcode = iota
)

type Definition struct {
    Name string
    OperandWidths []int
}

var definitions = map[Opcode]*Definition{
    OpConstant: {"OpConstant", []int{2}}, // 2-byte long single operand
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
        }

        offset += width
    }

    return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
    return binary.BigEndian.Uint16(ins)
}
