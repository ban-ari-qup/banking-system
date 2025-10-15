package main

import (
	"fmt"
	"mfp/account"
	"mfp/api"
)

func main() {
	accountList := account.NewAccountList()
	accountList.LoadFromFile("accounts.json") // Загружаем существующие данные
	server := api.NewServer(accountList)
	server.Start()
	// example := accountList.GetAccounts()
	// for _, acc := range example {
	// 	printAccountDetails(acc)
	// }
	// for {
	// 	acc := createAccount(accountList)
	// 	printAccountDetails(acc)
	// 	if err := accountList.AddAccount(acc); err != nil {
	// 		fmt.Println("Error adding account to list:\n" + err.Error())
	// 	} else {
	// 		fmt.Println("Account successfully added to the list.")
	// 	}
	// }
}

func createAccount(accountList *account.AccountList) *account.Account {
	var password, name, phone string
	var age int

	fmt.Println("Write your name:")
	fmt.Scan(&name)

	fmt.Println("Write your age:")
	fmt.Scan(&age)

	fmt.Println("Write your phone:")
	fmt.Scan(&phone)

	fmt.Println("Write your password:")
	fmt.Scan(&password)

	acc := account.NewAccount(password, name, phone, age)

	if err := acc.Validate(); err != nil {
		panic("Error creating account:\n" + err.Error())
	}
	return acc
}
func printAccountDetails(acc *account.Account) {
	fmt.Println("\n=== Account created successfully! ===")
	fmt.Printf("Your account name: %s\n", acc.Name)
	fmt.Printf("Your account phone: %s\n", acc.Phone)
	fmt.Printf("Your account age: %d\n", acc.Age)
	fmt.Printf("Your Card Number: %s\n", acc.ID)
	fmt.Printf("Your account password: %s\n", acc.Password)
	fmt.Printf("Your account CVV code: %s\n", acc.CVC2)
	fmt.Printf("Your account balance: %.2f\n", acc.Balance)
	fmt.Printf("Your account created at: %s\n", acc.CreatedAt)
	fmt.Printf("Your account expired at: %s\n", acc.ExpiredAt)
}
