package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"github.com/acheong08/obsidian-sync/database"
	"golang.org/x/term"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
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
	db := database.NewDatabase()
	defer db.DBConnection.Close()
	err = db.NewUser(email, string(password), username)
	if err != nil {
		panic(err)
	}
}
