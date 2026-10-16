//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	var password string

	if len(os.Args) >= 2 {
		password = os.Args[1]
	} else {
		fmt.Print("Password: ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading password:", err)
			os.Exit(1)
		}

		password = strings.TrimRight(line, "\r\n")
	}

	if password == "" {
		fmt.Fprintln(os.Stderr, "password cannot be empty")
		os.Exit(1)
	}
	if len(password) > 72 {
		fmt.Fprintln(os.Stderr, "password exceeds 72 bytes")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bcrypt error:", err)
		os.Exit(1)
	}

	fmt.Println(string(hash))
}
