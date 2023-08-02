package main

import (
	"fmt"
	"log"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
    user, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Hello %s. KYS\n", user.Username)
    fmt.Printf("enter commands:\n")
    repl.Start(os.Stdin, os.Stdout)
}
