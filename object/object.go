package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
    "hash/fnv"
)

const (
    INTEGER_OBJ = "INTEGER"
    STRING_OBJ  = "STRING"
    BOOLEAN_OBJ = "BOOLEAN"
    NULL_OBJ = "NULL"
    RETURN_VALUE_OBJ = "RETURN_VALUE"
    ERROR_OBJ = "ERROR"
    FUNCTION_OBJ = "FUNCTION"
    BUILTIN_OBJ = "BUILTIN"
    ARRAY_OBJ = "ARRAY"
    HASHMAP_OBJ = "HASHMAP"
    COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
)

type ObjectType string

type BuiltinFunction func(args ...Object) Object

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Hashable interface {
    HashKey() HashKey
}

type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType { 
    return INTEGER_OBJ // I believe Go does an implicit conversion for some types
}

func (i *Integer) Inspect() string {
    return fmt.Sprintf("%d", i.Value)
}

type String struct {
    Value string
}

func (s *String) Type() ObjectType {
    return STRING_OBJ
}

func (s *String) Inspect() string {
    return s.Value
}

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType {
    return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
    return fmt.Sprintf("%t", b.Value)
}

type Null struct {
}

func (n *Null) Type() ObjectType {
    return NULL_OBJ
}

func (n *Null) Inspect() string {
    return "null"
}

type ReturnValue struct {
    Value Object
}

func (rv *ReturnValue) Type() ObjectType {
    return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
    return rv.Value.Inspect()
}

type Error struct {
    Message string
}

func (e *Error) Type() ObjectType {
    return ERROR_OBJ
}

func (e *Error) Inspect() string {
    return "ERROR: " + e.Message
}

type Function struct {
    Parameters []*ast.Identifier
    Body *ast.BlockStatement
    Env *Environment
}

func (f *Function) Type() ObjectType {
    return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
    var out bytes.Buffer

    params := []string{}
    for _, p := range f.Parameters {
        params = append(params, p.String())
    }

    out.WriteString("fn")
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") {\n")
    out.WriteString(f.Body.String())
    out.WriteString("\n}")

    return out.String()
}

type Builtin struct {
    Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
    return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
    return "builtin function"
}

type Array struct {
    Elements []Object
}

func (ao *Array) Type() ObjectType {
    return ARRAY_OBJ
}

func (ao *Array) Inspect() string {
    var out bytes.Buffer

    elements := []string{}
    for _, e := range ao.Elements {
        elements = append(elements, e.Inspect())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}

type HashKey struct {
    Type ObjectType
    Value uint64
}

func (b *Boolean) HashKey() HashKey {
    var value uint64

    if b.Value {
        value = 1
    } else {
        value = 0
    }

    return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
    return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
    h := fnv.New64a()
    h.Write([]byte(s.Value))

    return HashKey{Type: s.Type(), Value: h.Sum64()}
}


type HashPair struct { // we use this for the Inspect method of HashMap
    Key Object
    Value Object
}

type HashMap struct {
    Pairs map[HashKey]HashPair
}

func (hm *HashMap) Type() ObjectType {
    return HASHMAP_OBJ
}

func (hm *HashMap) Inspect() string {
    var out bytes.Buffer

    pairs := []string{}
    for _, pair := range hm.Pairs { // the only usage of HashPair.Key
        pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
    }

    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString("}")

    return out.String()
}

type CompiledFunction struct {
    Instructions []byte
    NumLocals int
    NumParams int
}

func (cf *CompiledFunction) Type() ObjectType {
    return COMPILED_FUNCTION_OBJ
}

func (cf *CompiledFunction) Inspect() string {
    return fmt.Sprintf("CompiledFunction[%p]", cf)
}
