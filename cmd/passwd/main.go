package main

import (
	"fmt"
	"github.com/fzakaria/surf-journal/database"
	"github.com/fzakaria/surf-journal/passwords"
	"golang.org/x/term"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "passwd username")
		os.Exit(1)
	}

	password, err := term.ReadPassword(0)
	if err != nil {
		panic(err)
	}

	username := os.Args[1]
	serialized, err := passwords.NewSerializedPassword(string(password))
	if err != nil {
		panic(err)
	}

	err = database.AddPassword(username, serialized)
	if err != nil {
		panic(err)
	}
}
