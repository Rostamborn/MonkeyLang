package main

import "lexer"
import "core:fmt"

main :: proc() {
    lex := lexer.new_lexer("+  ")
    // fmt.println(lexer.next_token(lex))
    lexer.next_token(lex)
}
