package database

import "mfp/account"

func (r *Repository) CreateAccount(acc *account.Account) error {
	query := `INSERT INTO accounts (id, password, cvc2, balance, name, phone, age, created_at, expired_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, acc.ID, acc.Password, acc.CVC2, acc.Balance, acc.Name, acc.Phone, acc.Age, acc.CreatedAt, acc.ExpiredAt)
	return err
}

func (r *Repository) GetAccount(id string) error {
	query := `SELECT id, password, cvc2, balance, name, phone, age, created_at, expired_at FROM accounts WHERE id = $1`

	_, err := r.db.Exec(query, id)
	return err
}
