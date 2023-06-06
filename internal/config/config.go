package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	envFile := fmt.Sprintf(`env/%s.env`, os.Getenv("APP_ENV"))
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error load .env file")
	}
}
