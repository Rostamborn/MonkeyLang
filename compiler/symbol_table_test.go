package compiler 

import (
    "testing"
)

func TestDefine(t *testing.T) {
    expected := map[string]Symbol{
        "a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
        "b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
        "c": Symbol{Name: "c", Scope: LocalScope, Index: 0},
        "d": Symbol{Name: "d", Scope: LocalScope, Index: 1},
        "e": Symbol{Name: "e", Scope: LocalScope, Index: 0},
        "f": Symbol{Name: "f", Scope: LocalScope, Index: 1},
    }
    global := NewSymTable()
    a := global.Define("a")
    if a != expected["a"] {
        t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
    }
    b := global.Define("b")
    if b != expected["b"] {
        t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
    }
    firstLocal := NewEnclosedSymTable(global)
    c := firstLocal.Define("c")
    if c != expected["c"] {
        t.Errorf("expected c=%+v, got=%+v", expected["c"], c)
    }
    d := firstLocal.Define("d")
    if d != expected["d"] {
        t.Errorf("expected d=%+v, got=%+v", expected["d"], d)
    }
    secondLocal := NewEnclosedSymTable(firstLocal)
    e := secondLocal.Define("e")
    if e != expected["e"] {t.Errorf("expected e=%+v, got=%+v", expected["e"], e)
    }
    f := secondLocal.Define("f")
    if f != expected["f"] {
        t.Errorf("expected f=%+v, got=%+v", expected["f"], f)
    }
}

func TestResolveGlobal(t *testing.T) {
    global := NewSymTable()

    global.Define("a")
    global.Define("b")

    expected := []Symbol{
        Symbol{Name: "a", Scope: GlobalScope, Index: 0},
        Symbol{Name: "b", Scope: GlobalScope, Index: 1},
    }

    for _, sym := range expected {
        result, ok := global.Resolve(sym.Name)
        if !ok {
            t.Errorf("name %s not resolvable", sym.Name)
        }

        if result != sym {
            t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
        }
    }
}

func TestResolveLocal(t *testing.T) {
    global := NewSymTable()
    global.Define("a")
    global.Define("b")
    local := NewEnclosedSymTable(global)
    local.Define("c")
    local.Define("d")
    expected := []Symbol{
        Symbol{Name: "a", Scope: GlobalScope, Index: 0},
        Symbol{Name: "b", Scope: GlobalScope, Index: 1},
        Symbol{Name: "c", Scope: LocalScope, Index: 0},
        Symbol{Name: "d", Scope: LocalScope, Index: 1},
    }
    for _, sym := range expected {
        result, ok := local.Resolve(sym.Name)
        if !ok {
            t.Errorf("name %s not resolvable", sym.Name)
            continue
        }
        if result != sym {
            t.Errorf("expected %s to resolve to %+v, got=%+v",
                sym.Name, sym, result)
        }
    }
}

func TestResolveNestedLocal(t *testing.T) {
    global := NewSymTable()
    global.Define("a")
    global.Define("b")
    firstLocal := NewEnclosedSymTable(global)
    firstLocal.Define("c")
    firstLocal.Define("d")
    secondLocal := NewEnclosedSymTable(firstLocal)
    secondLocal.Define("e")
    secondLocal.Define("f")
    tests := []struct {
        table *SymTable
        expectedSymbols []Symbol
    }{
        {
            firstLocal,
            []Symbol{
                Symbol{Name: "a", Scope: GlobalScope, Index: 0},
                Symbol{Name: "b", Scope: GlobalScope, Index: 1},
                Symbol{Name: "c", Scope: LocalScope, Index: 0},
                Symbol{Name: "d", Scope: LocalScope, Index: 1},
            },
        },
        {
            secondLocal,
            []Symbol{
                Symbol{Name: "a", Scope: GlobalScope, Index: 0},
                Symbol{Name: "b", Scope: GlobalScope, Index: 1},
                Symbol{Name: "e", Scope: LocalScope, Index: 0},
                Symbol{Name: "f", Scope: LocalScope, Index: 1},
            },
        },
    }
    for _, tt := range tests {
        for _, sym := range tt.expectedSymbols {
            result, ok := tt.table.Resolve(sym.Name)
            if !ok {
                t.Errorf("name %s not resolvable", sym.Name)
                continue
            }
            if result != sym {
                t.Errorf("expected %s to resolve to %+v, got=%+v",
                    sym.Name, sym, result)
            }
        }
    }
}
