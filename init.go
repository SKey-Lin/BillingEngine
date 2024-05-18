package main

import (
	"log"

	"github.com/joho/godotenv"

	"squalux.com/skey/lending/models"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env file: %s", err)
	}

	models.ConnectDatabase()
}
