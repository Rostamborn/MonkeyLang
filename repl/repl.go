package repl

import(
    "bufio"
    "fmt"
    "io"
    "monkey/lexer"
    "monkey/parser"
    "monkey/evaluator"
    "monkey/object"
)

const PROMPT = "$ "

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)
    env := object.NewEnvironment()

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

        evalueated := evaluator.Eval(program, env)
        if evalueated != nil {
            io.WriteString(out, evalueated.Inspect())
            io.WriteString(out, "\n")
        }

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
