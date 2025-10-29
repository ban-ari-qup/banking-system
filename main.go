package main

import (
	"log"
	"mfp/api"
	"mfp/database"
	"mfp/session"
)

func main() {
	repo, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	log.Println("Database connected successsfully!")

	sessionManager := session.NewSessionManager()

	// СТАНОВИТСЯ:
	server := api.NewServer(repo, sessionManager) // ← Передаем repo вместо accountList
	server.Start()
}
