package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"

	// "monkey/evaluator"
	// "monkey/object"
	"monkey/compiler"
	"monkey/vm"
)

const PROMPT = "$ "

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    // env := object.NewEnvironment()
    constants := []object.Object{}
    globals := make([]object.Object, vm.GlobalSize)
    symTable := compiler.NewSymTable()

    for {
        fmt.Print(PROMPT)
        scanned := scanner.Scan()
        if !scanned {
            return
        }
        line := scanner.Text()

        lex := lexer.NewLexer(line)
        p := parser.NewParser(lex)
        program := p.ParseProgram()
        if len(p.Errors()) != 0 {
            PrintParserErrors(out, p.Errors())
            continue
        }

        comp := compiler.New_Compiler_With_States(constants, symTable)
        err := comp.Compile(program)
        if err != nil {
            fmt.Fprintf(out, "Compilation failed:\n %s\n", err)
            continue
        }

        virt_machine := vm.New_VM_With_Global_Store(comp.Bytecode(), globals)
        err = virt_machine.Run()
        if err != nil {
            fmt.Fprintf(out, "Execution failed:\n %s\n", err)
            continue
        }

        stack_top := virt_machine.LastPopped()
        if stack_top != nil {
            io.WriteString(out, stack_top.Inspect())
            io.WriteString(out, "\n")
        }
        

        // evalueated := evaluator.Eval(program, env)
        // if evalueated != nil {
        //     io.WriteString(out, evalueated.Inspect())
        //     io.WriteString(out, "\n")
        // }
        //
        // io.WriteString(out, program.String())
        // io.WriteString(out, "\n")
        // for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
        //     fmt.Printf("%+v\n", tok)
        // }
    }
}

func PrintParserErrors(out io.Writer, errors []string) {
    io.WriteString(out, "seems like you suck at writing monkey!\n")
    io.WriteString(out, " parser errors:\n")
    for _, msg := range errors {
        io.WriteString(out, "\t" + msg + "\n")
    }
}
