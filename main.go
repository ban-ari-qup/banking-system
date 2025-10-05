package main

import (
	"fmt"
	"mfp/account"
)

func main() {
	printAccountDetails(createAccount())
}

func createAccount() *account.Account {
	var password, name, phone string
	var age int

	fmt.Println("Write your password:")
	fmt.Scan(&password)

	fmt.Println("Write your name:")
	fmt.Scan(&name)

	fmt.Println("Write your phone:")
	fmt.Scan(&phone)

	fmt.Println("Write your age:")
	fmt.Scan(&age)

	acc := account.NewAccount(password, name, phone, age)

	if err := acc.Validate(); err != nil {
		panic("Error creating account:\n" + err.Error())
	}

	accountList := account.NewAccountList()
	accountList.AddAccount(acc)
	return acc
}
func printAccountDetails(acc *account.Account) {
	fmt.Println("\n=== Account created successfully! ===")
	fmt.Printf("Your Card Number: %s\n", acc.ID)
	fmt.Printf("Your account password: %s\n", acc.Password)
	fmt.Printf("Your account name: %s\n", acc.Name)
	fmt.Printf("Your account phone: %s\n", acc.Phone)
	fmt.Printf("Your account CVV code: %s\n", acc.CVC2)
	fmt.Printf("Your account age: %d\n", acc.Age)
	fmt.Printf("Your account balance: %.2f\n", acc.Balance)
	fmt.Printf("Your account created at: %s\n", acc.CreatedAt)
	fmt.Printf("Your account expired at: %s\n", acc.ExpiredAt)
}
