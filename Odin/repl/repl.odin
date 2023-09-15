package repl

import "core:io"
import "core:bufio"
import "core:fmt"
import "../lexer"
import "../token"
import "../parser"
import "../ast"

start :: proc(stream: io.Stream) {
    scanner := new(bufio.Scanner)
    scanner = bufio.scanner_init(scanner, stream)

    for {
        fmt.printf(">> ")
        scanned := bufio.scanner_scan(scanner)
        if !scanned {
            return
        }
        line := bufio.scanner_text(scanner)

        lex := lexer.new_lexer(line)
        p := parser.new_parser(lex)
        program := parser.parse_program(p)
        if len(p.errors) > 0 {
            print_parser_errors(stream, p.errors)
        }

        io.write(stream, transmute([]u8)string("parsed program:\n"))
        io.write(stream, transmute([]u8)ast.to_string(program))
        fmt.println()
        
        free_all(context.temp_allocator)
    }
}

print_parser_errors :: proc(stream: io.Stream, errors: [dynamic]string) {
    io.write(stream, transmute([]u8)string("seems like you suck at Monkey-Lang\n"))
    io.write(stream, transmute([]u8)string("parser errors:\n"))
    for msg in errors {
        io.write(stream, transmute([]u8)string("\t"))
        io.write(stream, transmute([]u8)msg)
        io.write(stream, transmute([]u8)string("\n"))
    }
}
