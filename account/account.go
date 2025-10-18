package account

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// структура аккаунта
type Account struct {
	ID        string    `json:"id"`
	Password  string    `json:"password"` //хешированный пароль
	CVC2      string    `json:"cvc2"`     // Card Verification Code
	Balance   float64   `json:"balance"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// функция создания нового аккаунта
func NewAccount(password, name, phone string, age int) *Account {
	if err := validatePassword(password); err != nil {
		panic(fmt.Sprintf("Account creation failed: %v", err))
	}

	generator := NewCardGenerator()
	return &Account{
		ID:        generator.GenerateCardNumber(), //генерация номера аккаунта
		Password:  hashPassword(password),         //хеширование пароля
		CVC2:      generator.GenerateCVC(),        //генерация CVC2
		Balance:   0,
		Name:      name,
		Phone:     phone,
		Age:       age,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().AddDate(5, 0, 0), // срок действия аккаунта 5 лет
	}
}

// проверка на истечение срока действия аккаунта
func (a *Account) IsExpired() bool {
	return time.Now().After(a.ExpiredAt)
}

// пополнение баланса
func (a *Account) Deposit(amount float64) error {
	if a.IsExpired() {
		return fmt.Errorf("account is expired")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	a.Balance += amount
	return nil
}

// снятие средств
func (a *Account) Withdraw(amount float64) error {
	if a.IsExpired() {
		return fmt.Errorf("account is expired")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if a.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}
	a.Balance -= amount
	return nil
}

// валидация пароля
func validatePassword(a string) error {
	if len(a) != 4 {
		return fmt.Errorf("password must be exactly 4 digits")
	}
	for _, ch := range a {
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
