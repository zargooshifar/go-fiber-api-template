package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var AppName = "go-fiber-api-template"

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	// Return the value of the variable
	return os.Getenv(key)
}
