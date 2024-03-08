package main

import (
	"github.com/alexmeuer/juke/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local", ".env")
	cmd.Execute()
}
