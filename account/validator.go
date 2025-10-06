package account

import "fmt"

type ValidationRule func(*Account) error

var validationRules = []ValidationRule{
	validateAge,
	// validatePassword,
	validatePhone,
}

func validateAge(a *Account) error {
	if a.Age < 18 {
		return fmt.Errorf("age must be at least 18")
	}
	return nil
}

//	func validatePassword(a *Account) error {
//		if len(a.Password) != 4 {
//			return fmt.Errorf("password must be exactly 4 digits")
//		}
//		for _, ch := range a.Password {
//			if ch < '0' || ch > '9' {
//				return fmt.Errorf("password must contain only digits")
//			}
//		}
//		return nil
//	}
func validatePhone(a *Account) error {
	if len(a.Phone) != 11 {
		return fmt.Errorf("phone number must be exactly 11 digits")
	}
	for i, ch := range a.Phone {
		if (i == 0 || i == 1) && ch != '7' {
			return fmt.Errorf("phone number must start with '77'")
		}
		if ch < '0' || ch > '9' {
			return fmt.Errorf("phone number must contain only digits")
		}
	}
	return nil
}

func (a *Account) Validate() error {
	for _, rule := range validationRules {
		if err := rule(a); err != nil {
			return err
		}
	}
	return nil
}
