package file

import (
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
    "monkey/repl"
	"os"
)

func Run_file(file_name string) {
    program_text, err := os.ReadFile(file_name)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }

    evn := object.NewEnvironment()
    lex := lexer.NewLexer(string(program_text))
    parser := parser.NewParser(lex)
    program := parser.ParseProgram()

    if len(parser.Errors()) != 0 {
        repl.PrintParserErrors(os.Stdout, parser.Errors())
    }

    _ = evaluator.Eval(program, evn)
}
