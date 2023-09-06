package main

import "lexer"
import "core:fmt"
import "core:os"
import "repl"

main :: proc() {
    // lex := lexer.new_lexer("+  ")
    // // fmt.println(lexer.next_token(lex))
    // lexer.next_token(lex)
    fmt.println("Monkey Lang")
    stream := os.stream_from_handle(os.stdin)
    repl.start(stream)
}
