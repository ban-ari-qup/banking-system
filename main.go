package main

import (
	"fmt"
	"log"
	"mfp/account"
	"mfp/api"
)

func main() {
	// Инициализация системы
	accountList := account.NewAccountList()

	// Загружаем существующие данные
	if err := accountList.LoadFromFile("accounts.json"); err != nil {
		fmt.Println("Error loading accounts:", err)
	}

	// Запускаем сервер API
	server := api.NewServer(accountList)
	log.Println("Starting banking system")
	server.Start()
}
