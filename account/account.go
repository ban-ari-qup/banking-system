package account

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// структура аккаунта
type Account struct {
	ID           string        `json:"id"`
	Password     string        `json:"password"` //хешированный пароль
	CVC2         string        `json:"cvc2"`     // Card Verification Code
	Balance      float64       `json:"balance"`
	Name         string        `json:"name"`
	Phone        string        `json:"phone"`
	Age          int           `json:"age"`
	CreatedAt    time.Time     `json:"created_at"`
	ExpiredAt    time.Time     `json:"expired_at"`
	Transactions []Transaction `json:"transactions"`
}

// функция создания нового аккаунта
func NewAccount(password, name, phone string, age int) *Account {
	if err := validatePassword(password); err != nil {
		panic(fmt.Sprintf("Account creation failed: %v", err))
	}

	generator := NewCardGenerator()
	return &Account{
		ID:           generator.GenerateCardNumber(), //генерация номера аккаунта
		Password:     hashPassword(password),         //хеширование пароля
		CVC2:         generator.GenerateCVC(),        //генерация CVC2
		Balance:      0,
		Name:         name,
		Phone:        phone,
		Age:          age,
		CreatedAt:    time.Now(),
		ExpiredAt:    time.Now().AddDate(5, 0, 0), // срок действия аккаунта 5 лет
		Transactions: []Transaction{},
	}
}

// проверка на истечение срока действия аккаунта
func (acc *Account) IsExpired() bool {
	return time.Now().After(acc.ExpiredAt)
}

// пополнение баланса
func (acc *Account) Deposit(amount float64) error {
	if acc.IsExpired() {
		return fmt.Errorf("account is expired")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	acc.Balance += amount

	transaction := Transaction{
		ID:          len(acc.Transactions) + 1,
		Type:        "deposit",
		FromAccount: "",
		ToAccount:   acc.ID,
		Amount:      amount,
		Timestamp:   time.Now(),
		Status:      "completed",
	}
	acc.Transactions = append(acc.Transactions, transaction)

	return nil
}

// снятие средств
func (acc *Account) Withdraw(amount float64) error {
	if acc.IsExpired() {
		return fmt.Errorf("account is expired")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if acc.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}
	acc.Balance -= amount

	transaction := Transaction{
		ID:          len(acc.Transactions) + 1,
		Type:        "withdrawal",
		FromAccount: acc.ID,
		ToAccount:   "",
		Amount:      amount,
		Timestamp:   time.Now(),
		Status:      "completed",
	}
	acc.Transactions = append(acc.Transactions, transaction)

	return nil
}

// валидация пароля
func validatePassword(pass string) error {
	if len(pass) != 4 {
		return fmt.Errorf("password must be exactly 4 digits")
	}
	for _, ch := range pass {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("password must contain only digits")
		}
	}
	return nil
}

// хеширование пароля
func hashPassword(passowrd string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passowrd), 14)
	if err != nil {
		panic(fmt.Sprintf("Password hashing failed: %v", err))
	}
	return string(bytes)
}

// проверка пароля
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
