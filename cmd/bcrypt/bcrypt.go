package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	//for i, arg := range os.Args {
	//	fmt.Println(i, arg)
	//}

	if len(os.Args) < 2 {
		fmt.Println("Usage: bcrypt <hash|compare> <password> [<hash>]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "hash":
		if len(os.Args) < 3 {
			fmt.Println("Usage: bcrypt hash <password>")
			os.Exit(1)
		}
		hash(os.Args[2])
	case "compare":
		if len(os.Args) < 4 {
			fmt.Println("Usage: bcrypt compare <password> <hash>")
			os.Exit(1)
		}
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Println("Usage: bcrypt <hash|compare> <password> [<hash>]")
	}
}

func hash(password string) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing:", password, err)
	}
	fmt.Println("Hashed password:", string(hashedBytes))
}

func compare(password, hash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Password does not match hash:", password, hash, err)
		return
	}
	fmt.Println("Password matches hash:", password, hash)
}
