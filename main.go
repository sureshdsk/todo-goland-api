package main

import (
	"github.com/joho/godotenv"
	"github.com/sureshdsk/todo-goland-api/internal/db"
	"github.com/sureshdsk/todo-goland-api/internal/todo"
	"github.com/sureshdsk/todo-goland-api/internal/transport"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}
	pgUrl := os.Getenv("DB_URL")

	d, err := db.New(pgUrl)
	if err != nil {
		log.Fatal(err)
	}
	svc := todo.NewService(d)
	server := transport.NewServer(svc)
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
