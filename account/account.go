package account

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID        string    `json:"id"`
	Password  string    `json:"password"`
	CVC2      string    `json:"cvc2"`
	Balance   float64   `json:"balance"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewAccount(password, name, phone string, age int) *Account {
	if err := validatePassword(password); err != nil {
		panic("Error creating account:\n" + err.Error())
	}
	generator := NewCardGenerator()
	return &Account{
		ID:        generator.GenerateCardNumber(),
		Password:  hashPassword(password),
		CVC2:      generator.GenerateCVC(),
		Balance:   0,
		Name:      name,
		Phone:     phone,
		Age:       age,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().AddDate(5, 0, 0),
	}
}

func (a *Account) IsExpired() bool {
	return time.Now().After(a.ExpiredAt)
}

func hashPassword(passowrd string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passowrd), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

//	func Transfer(from, to *Account, amount float64) error {
//		if from.IsExpired() {
//			return fmt.Errorf("account is expired")
//		}
//		if to.IsExpired() {
//			return fmt.Errorf("account is unknown")
//		}
//		if amount <= 0 {
//			return fmt.Errorf("amount must be positive")
//		}
//		if from.Balance < amount {
//			return fmt.Errorf("insufficient funds")
//		}
//		from.Balance -= amount
//		to.Balance += amount
//		return nil
//	}
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
