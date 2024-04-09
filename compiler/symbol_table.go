package compiler

type SymScope string

const (
    GlobalScope SymScope = "GLOBAL"
    LocalScope SymScope = "LOCAL"
)

type Symbol struct {
    Name string
    Scope SymScope
    Index int
}

type SymTable struct {
    Outer *SymTable
    store map[string]Symbol
    num_def int
}

func NewSymTable() *SymTable {
    m := make(map[string]Symbol)
    return &SymTable{store: m, num_def: 0}
}

func NewEnclosedSymTable(outer *SymTable) *SymTable {
    s := NewSymTable()
    s.Outer = outer
    return s
}

func (s *SymTable) Define(name string) Symbol {
    symbol := Symbol{Name: name, Index: s.num_def}
    if s.Outer == nil {
        symbol.Scope = GlobalScope
    } else {
        symbol.Scope = LocalScope
    }
    s.store[name] = symbol
    s.num_def++

    return symbol
}

func (s *SymTable) Resolve(name string) (Symbol, bool) {
    sym, ok := s.store[name]
    if !ok && s.Outer != nil {
        sym, ok = s.Outer.Resolve(name)
        return sym, ok
    }
    return sym, ok
}
