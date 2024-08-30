package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvPORT() string {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	return os.Getenv("PORT")
}
