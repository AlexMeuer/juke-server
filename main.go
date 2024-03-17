package main

import (
	"fmt"

	"github.com/alexmeuer/juke/cmd"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local", ".env"); err != nil {
		fmt.Printf("Failed to load .env file: %v\n", err)
	}
	cmd.Execute()
}
