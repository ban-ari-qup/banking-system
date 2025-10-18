package account

import "fmt"

// тип функции для валидации аккаунта
type ValidationRule func(*Account) error

// список правил валидации
var validationRules = []ValidationRule{
	validateAge,
	validatePhone,
}

// метод валидации аккаунта
func (acc *Account) Validate() error {
	for _, rule := range validationRules {
		if err := rule(acc); err != nil {
			return err
		}
	}
	return nil
}

// функция валидации возраста
func validateAge(acc *Account) error {
	if acc.Age < 18 {
		return fmt.Errorf("age must be at least 18")
	}
	return nil
}

// функция валидации номера телефона
func validatePhone(acc *Account) error {
	if len(acc.Phone) != 11 {
		return fmt.Errorf("phone number must be exactly 11 digits")
	}
	for i, ch := range acc.Phone {
		if i < 2 && ch != '7' {
			return fmt.Errorf("phone number must start with '77'")
		}
		if ch < '0' || ch > '9' {
			return fmt.Errorf("phone number must contain only digits")
		}
	}
	return nil
}
