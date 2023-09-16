package object

Environment :: struct {
    store: map[string]^Object,
    outer: ^Environment,
}

new_env :: proc() -> ^Environment {
    st := make(map[string]^Object)
    env := new(Environment)
    env.store = st

    return env
}

new_enclosed_env :: proc(outer: ^Environment) -> ^Environment {
    env := new_env()
    env.outer = outer

    return env
}

env_get :: proc(env: ^Environment, name: string) -> (^Object, bool) {
    obj, ok := env.store[name]
    if !ok && env.outer != nil {
        obj, ok = env_get(env.outer, name)
    }

    return obj, ok
}

env_set :: proc(env: ^Environment, name: string, value: ^Object) -> ^Object {
    env.store[name] = value

    return value
}
