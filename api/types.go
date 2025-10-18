package api

import "mfp/account"

// запрос на создание аккаунта
type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	Age       int    `json:"age"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

// ответ с информацией об аккаунте
type AccountResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	Phone     string  `json:"phone"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
}

// преобразование аккаунта в ответ API
func AccountToResponse(acc *account.Account) AccountResponse {
	return AccountResponse{
		ID:        acc.ID,
		Name:      acc.Name,
		Age:       acc.Age,
		Phone:     acc.Phone,
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
