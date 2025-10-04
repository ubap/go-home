package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"goHome/auth"

	"golang.org/x/term"
)

const dbFile = "data/users.json"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cmd/user-manager/main.go <add> <username>")
		os.Exit(1)
	}

	authManager, err := auth.NewBasicAuthManager(dbFile)
	if err != nil {
		log.Fatalf("Failed to initialize auth manager: %v", err)
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run cmd/user-manager/main.go add <username>")
			os.Exit(1)
		}
		username := os.Args[2]
		addUser(authManager, username)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func addUser(authManager *auth.UserManager, username string) {
	fmt.Print("Enter Password: ")
	// Read password securely without echoing to the terminal
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	password := string(bytePassword)
	fmt.Println()

	err = authManager.AddUser(username, password)
	if err != nil {
		log.Fatalf("Error adding user: %v", err)
	}

	fmt.Printf("Successfully added user '%s' to %s\n", username, dbFile)
}
