package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DB_URL string

func LoadConfig() {
  err := godotenv.Load()
  if err != nil {
	log.Fatal("Error loading .env file")
  }
  DB_URL = os.Getenv("DB_URL")
  if DB_URL == "" {
	log.Fatal("DB_URL not set in enviromet")
  }
}