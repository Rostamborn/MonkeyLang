package object

import "core:testing"

@(test)
test_string_hashkey :: proc(t: ^testing.T) {
    hello1 := new_obj(String)
    hello1.value = "Hello World"
    hello2 := new_obj(String)
    hello2.value = "Hello World"

    diff1 := new_obj(String)
    diff1.value = "My name is johnny"
    diff2 := new_obj(String)
    diff2.value = "My name is johnny"

    if hash_key(hello1) != hash_key(hello2) {
        testing.error(t, "strings with same content have different hash keys")
    }

    if hash_key(diff1) != hash_key(diff2) {
        testing.error(t, "strings with same content have different hash keys")
    }

    if hash_key(hello1) == hash_key(diff1) {
        testing.error(t, "strings with different content have same hash keys")
    }
}
