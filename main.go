package main

import (
	"fmt"
	"log"
	"monkey/repl"
	"os"
	"os/user"
    "monkey/file"
)

func main() {
    user, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }
    args := os.Args

    if len(args) == 2 {
        file.Run_file(args[1])
    } else {
        fmt.Printf("Hello %s. KYS\n", user.Username)
        fmt.Printf("enter commands:\n")
        repl.Start(os.Stdin, os.Stdout)
    }
}


