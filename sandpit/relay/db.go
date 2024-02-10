package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func OpenDatabase() error {

	DB_USER := os.Getenv("DB_USER")
	if DB_USER == "" {
		log.Fatal("DB_USER environment variable not set")
	}
	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		log.Fatal("DB_NAME environment variable not set")
	}
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	if DB_PASSWORD == "" {
		log.Fatal("DB_PASSWORD environment variable not set")
	}
	DB_HOST := os.Getenv("DB_HOST")
	if DB_HOST == "" {
		log.Fatal("DB_HOST environment variable not set")
	}
	DB_PORT := os.Getenv("DB_PORT")
	if DB_PORT == "" {
		log.Fatal("DB_PORT environment variable not set")
	}

	var err error
	DB, err = sql.Open("postgres", "user="+DB_USER+" dbname="+DB_NAME+" password="+DB_PASSWORD+" host="+DB_HOST+" port="+DB_PORT+" sslmode=disable")
	if err != nil {
		return err
	}
	return nil
}

func CloseDatabase() error {
	return DB.Close()
}
