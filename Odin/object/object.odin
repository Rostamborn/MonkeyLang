package object

import "../ast"
import "core:hash"
import "core:fmt"
import "core:mem"
import "core:bytes"
import "core:strings"

ObjectType :: enum {
    INTEGER,
    BOOLEAN,
    STRING,
    NULL,
    RETURN_VALUE,
    ERROR,
    FUNCTION,
    BUILTIN,
    ARRAY,
    HASHMAP,
}

any_obj :: union {
    ^Integer,
    ^Boolean,
    ^String,
    ^Null,
    ^Return_Value,
    ^Error,
    ^Function,
    ^Builtin,
    ^Array,
    ^Hash_Map,
}

Builtin_Func :: proc(args: ..^Object) -> ^Object

new_obj :: proc($T: typeid, allocator := context.allocator) -> ^T where ast.has_field(T, "derived") {
    obj := new(T, allocator)
    obj.derived = obj 

    return obj 
}

new_bool_obj :: proc(val: bool, allocator := context.allocator) -> ^Boolean {
    obj := new(Boolean, allocator)
    obj.derived = obj
    obj.value = val

    return obj
}

// Objects

Object :: struct {
    derived: any_obj,
}

Integer :: struct {
    using obj: Object,
    value: int,
}

Boolean :: struct {
    using obj: Object,
    value: bool,
}

String :: struct {
    using obj: Object,
    value: string,
}

Null :: struct {
    using obj: Object,
}

Return_Value :: struct {
    using obj: Object,
    value: ^Object,
}

Error :: struct {
    using obj: Object,
    msg: string,
}

Function :: struct {
    using obj: Object,
    params: [dynamic]^ast.Ident,
    body: ^ast.Block_Stmt,
    env: ^Environment,
}

Builtin :: struct {
    using obj: Object,
    func: Builtin_Func,
}

Array :: struct {
    using obj: Object,
    elems: [dynamic]^Object,
}

Hash_Map :: struct {
    using obj: Object,
    pairs: map[Hash_Key]Hash_Pair, // Hash_Pair contains the real key value pair.
}                                  // Hash_Key is needed so we can access myMap["key"]
                                   // because everytime "key" is a new object and we wouldn't get
Hash_Key :: struct {               // the desired value. Hence we need a hashed key.
    type: ObjectType,
    value: u64,
}

Hash_Pair :: struct {
    key: ^Object,
    value: ^Object,
}

// Helpers

obj_type :: proc(obj: ^Object) -> ObjectType {
    switch v in obj.derived {
        case ^Integer: return .INTEGER
        case ^Boolean: return .BOOLEAN
        case ^String: return .STRING
        case ^Null: return .NULL
        case ^Return_Value: return .RETURN_VALUE
        case ^Error: return .ERROR
        case ^Function: return .FUNCTION
        case ^Builtin: return .BUILTIN
        case ^Array: return .ARRAY
        case ^Hash_Map: return .HASHMAP
        case: return .NULL
    }
}

obj_inspect :: proc(obj: ^Object) -> string {
    #partial switch v in obj.derived {

    }
    return ""
}

integer_inspect :: proc(obj: ^Integer) -> string {
    return fmt.tprintf("%d", obj.value)
}

boolean_inspect :: proc(obj: ^Boolean) -> string {
    return fmt.tprintf("%t", obj.value)
}

string_inspect :: proc(obj: ^String) -> string {
    return obj.value
}

null_inspect :: proc(obj: ^Null) -> string {
    return fmt.tprint("null")
}

return_value_inspect :: proc(obj: ^Return_Value) -> string {
    return obj_inspect(obj.value)
}

error_inspect :: proc(obj: ^Error) -> string {
    return obj.msg 
}

function_inspect :: proc(obj: ^Function) -> string {
    out: bytes.Buffer
    params := [dynamic]string{}
    for p in obj.params {
        append(&params, ast.to_string(p))
    }
    res := strings.join(params[:], ", ", context.temp_allocator)
    defer delete(res, context.temp_allocator)

    bytes.buffer_write(&out, transmute([]u8)string("fn"))
    bytes.buffer_write(&out, transmute([]u8)string("("))
    bytes.buffer_write(&out, transmute([]u8)res)
    bytes.buffer_write(&out, transmute([]u8)string(") {\n"))
    bytes.buffer_write(&out, transmute([]u8)ast.to_string(obj.body))
    bytes.buffer_write(&out, transmute([]u8)string("\n}"))

    return bytes.buffer_to_string(&out)
}

builtin_inspect :: proc(obj: ^Builtin) -> string {
    return "builtin function"
}

array_inspect :: proc(obj: ^Array) -> string {
    out: bytes.Buffer
    params := [dynamic]string{}
    for p in obj.elems {
        append(&params, obj_inspect(p))
    }
    res := strings.join(params[:], ", ", context.temp_allocator)
    defer delete(res, context.temp_allocator)

    bytes.buffer_write(&out, transmute([]u8)string("["))
    bytes.buffer_write(&out, transmute([]u8)res)
    bytes.buffer_write(&out, transmute([]u8)string("]"))

    return bytes.buffer_to_string(&out)
}

hash_map_inspect :: proc(obj: ^Hash_Map) -> string {
    out: bytes.Buffer
    params := [dynamic]string{}
    for _, p in obj.pairs {
        str := fmt.tprintf("%s: %s", obj_inspect(p.key), obj_inspect(p.value))
        defer delete(str, context.temp_allocator)
        append(&params, str)
    }
    res := strings.join(params[:], ", ", context.temp_allocator)
    defer delete(res, context.temp_allocator)

    bytes.buffer_write(&out, transmute([]u8)string("{"))
    bytes.buffer_write(&out, transmute([]u8)res)
    bytes.buffer_write(&out, transmute([]u8)string("]"))

    return bytes.buffer_to_string(&out)
}

// Hash Helpers

// function overloading is not an option since we need to unwrap obj.derived
hash_key :: proc(obj: ^Object) -> Hash_Key {
    #partial switch v in obj.derived {
        case ^Integer: return hash_key_int(v)
        case ^String: return hash_key_str(v)
        case ^Boolean: return hash_key_bool(v)
        case: return Hash_Key{.NULL, 0}
    }
}

hash_key_int :: proc(obj: ^Integer) -> Hash_Key {
    return Hash_Key{obj_type(obj), u64(obj.value)}
}

hash_key_str :: proc(obj: ^String) -> Hash_Key {
    buf := transmute([]u8)obj.value
    h := hash.fnv64a(buf)

    return Hash_Key{obj_type(obj), h}
}

hash_key_bool :: proc(obj: ^Boolean) -> Hash_Key {
    value: u64

    if obj.value {
        value = 1
    } else {
        value = 0
    }

    return Hash_Key{obj_type(obj), value}
}
