package main

import (
	"fmt"
	"monkey/pkg/repl"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic("Cannot find user.")
	}
	fmt.Printf("Hello %s, this is the Monkey REPL!\n", usr.Name)
	// Start REPL here.
	repl.Start(os.Stdin, os.Stdout)
}
