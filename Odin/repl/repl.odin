package repl

import "core:io"
import "core:bufio"
import "core:fmt"
import "../lexer"
import "../token"

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
        tok := lexer.next_token(lex)
        for  tok.type != token.EOF {
            io.write(stream, transmute([]u8)tok.literal)
            io.write(stream, transmute([]u8)string("\n"))
            tok = lexer.next_token(lex)
        }
        fmt.println()
        free_all(context.temp_allocator)
    }
}
