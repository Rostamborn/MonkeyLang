package main

import "lexer"
import "core:fmt"
import "core:os"
import "repl"

main :: proc() {
    fmt.println("Monkey Lang")
    stream := os.stream_from_handle(os.stdin)
    repl.start(stream)
}
