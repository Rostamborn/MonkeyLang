package vm

import (
	"fmt"
    "testing"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

type vmTestCase struct {
    input string
    expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
    tests := []vmTestCase {
        { "1", 1 },
        { "2", 2 },
        { "1 + 2", 3 },
        { "1 - 2", -1 },
        { "1 * 2", 2 },
        { "4 / 2", 2 },
        // { "50 / 2 * 2 + 10 - 5", 55 },
        // { "5 * (2 + 10)", 60 },
        // { "-5", -5 },
        // { "-10 + 5", -5 },
        // { "-10 - 5", -15 },
        // { "-10 * 5", -50 },
        // { "-10 / 5", -2 },
        // { "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50 },
    }

    runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
    tests := []vmTestCase {
        {"true", true},
        {"false", false},
        {"1 < 2", true},
        {"1 > 2", false},
        {"1 < 1", false},
        {"1 > 1", false},
        {"1 == 1", true},
        {"1 != 1", false},
        {"1 == 2", false},
        {"1 != 2", true},
        {"true == true", true},
        {"false == false", true},
        {"true == false", false},
        {"true != false", true},
        {"false != true", true},
        {"(1 < 2) == true", true},
        {"(1 < 2) == false", false},
        {"(1 > 2) == true", false},
        {"(1 > 2) == false", true},
    }

    runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
    t.Helper()

    for _, tt := range tests {
        program := parse(tt.input)
        comp := compiler.New_Compiler()
        err := comp.Compile(program)
        if err != nil {
            t.Errorf("compilation failed: %s", err)
        }

        vm := New_VM(comp.Bytecode())
        err = vm.Run()
        if err != nil {
            t.Fatalf("vm error: %s", err)
        }

        stackElem := vm.LastPopped()

        testExpectedObject(t, tt.expected, stackElem)
    }
}

func parse(input string) *ast.Program {
    l := lexer.NewLexer(input)
    p := parser.NewParser(l)
    return p.ParseProgram()
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
    t.Helper()

    switch expected := expected.(type) {
    case int:
        err := testIntegerObject(int64(expected), actual)
        if err != nil {
            t.Errorf("testIntegerObject failed: %s", err)
        }
    case bool:
        err := testBooleanObject(bool(expected), actual)
        if err != nil {
            t.Errorf("testBooleanObject failed: %s", err)
        }
    }
}

func testIntegerObject(expected int64, actual object.Object) error {
    result, ok := actual.(*object.Integer)
    if !ok {
        return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
    }

    if result.Value != expected {
        return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
    }

    return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
    result, ok := actual.(*object.Boolean)
    if !ok {
        return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
    }

    if result.Value != expected {
        return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
    }

    return nil
}
