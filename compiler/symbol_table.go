package compiler

type SymScope string

const (
    GlobalScope SymScope = "GLOBAL"
)

type Symbol struct {
    Name string
    Scope SymScope
    Index int
}

type SymTable struct {
    store map[string]Symbol
    num_def int
}

func NewSymTable() *SymTable {
    m := make(map[string]Symbol)
    return &SymTable{store: m, num_def: 0}
}

func (s *SymTable) Define(name string) Symbol {
    symbol := Symbol{Name: name, Scope: GlobalScope, Index: s.num_def}
    s.store[name] = symbol
    s.num_def++

    return symbol
}

func (s *SymTable) Resolve(name string) (Symbol, bool) {
    sym, ok := s.store[name]
    return sym, ok
}
