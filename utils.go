package main

import (
	"os"

	"github.com/joho/godotenv"
)

func GetDBUrl() string {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("env variable db_url is not set")
	}
	return dbUrl
}
