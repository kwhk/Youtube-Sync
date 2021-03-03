package config

import (
	"os"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
)

func ConnectToDB() *pg.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := pg.Connect(&pg.Options{
		User: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Addr: os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Database: "postgres",
	})

	return db
}