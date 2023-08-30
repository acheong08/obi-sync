package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/acheong08/obsidian-sync/database/vault"
	"golang.org/x/term"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your full name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Print("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}

	// Strip newline characters
	name = strings.Trim(name, "\n")
	email = strings.Trim(email, "\n")
	password = []byte(strings.Trim(string(password), "\n"))

	err = vault.NewUser(email, string(password), name)
	if err != nil {
		panic(err)
	}
}
