package database

import (
	"database/sql"
	"fmt"
	"mfp/account"
	"time"
)

func (r *Repository) CreateAccount(acc *account.Account) error {
	query := `INSERT INTO accounts (id, password, cvc2, balance, name, phone, age, created_at, expired_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, acc.ID, acc.Password, acc.CVC2, acc.Balance, acc.Name, acc.Phone, acc.Age, acc.CreatedAt, acc.ExpiredAt)
	return err
}

func (r *Repository) GetAccount(id string) (*account.Account, error) {
	query := `SELECT id, password, cvc2, balance, name, phone, age, created_at, expired_at FROM accounts WHERE id = $1`

	row := r.db.QueryRow(query, id)
	return scanAccount(row)
}

func (r *Repository) GetAccountByPhone(phone string) (*account.Account, error) {
	query := `
		SELECT id, password, cvc2, balance, name, phone, age, created_at, expired_at
		FROM accounts WHERE phone = $1`

	row := r.db.QueryRow(query, phone)
	return scanAccount(row)
}

func (r *Repository) GetAccounts() ([]*account.Account, error) {
	query := `SELECT id, password, cvc2, balance, name, phone, age, created_at, expired_at FROM accounts`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account
	for rows.Next() {
		acc, err := scanAccountFromRows(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (r *Repository) Deposit(accountID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	result, err := tx.Exec(query, amount, accountID)
	if err != nil {
		return fmt.Errorf("deposit failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	_, err = tx.Exec(`
        INSERT INTO transactions (type, from_account, to_account, amount, timestamp, status, account_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"deposit", "", accountID, amount, time.Now(), "completed", accountID,
	)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %v", err)
	}

	return tx.Commit()
}

func (r *Repository) Withdraw(accountID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw amount must be positive")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = $1", accountID).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("account not found: %v", err)
	}

	if currentBalance < amount {
		return fmt.Errorf("insufficient funds: have %.2f, need %.2f", currentBalance, amount)
	}

	query := `UPDATE accounts SET balance = balance - $1 WHERE id = $2`
	result, err := tx.Exec(query, amount, accountID)
	if err != nil {
		return fmt.Errorf("withdraw failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	_, err = tx.Exec(`
        INSERT INTO transactions (type, from_account, to_account, amount, timestamp, status, account_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"withdraw", accountID, "", amount, time.Now(), "completed", accountID,
	)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %v", err)
	}

	return tx.Commit()
}

func (r *Repository) Transfer(fromAccount, toAccount string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	if fromAccount == toAccount {
		return fmt.Errorf("cannot transfer to the same account")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = $1", fromAccount).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("sender account not found: %v", err)
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds: have %.2f, need %.2f", fromBalance, amount)
	}

	var toAccountExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", toAccount).Scan(&toAccountExists)
	if err != nil || !toAccountExists {
		return fmt.Errorf("receiver account not found")
	}

	queryDeduct := `UPDATE accounts SET balance = balance - $1 WHERE id = $2`
	resultDeduct, err := tx.Exec(queryDeduct, amount, fromAccount)
	if err != nil {
		return fmt.Errorf("transfer deduction failed: %v", err)
	}

	queryAdd := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	resultAdd, err := tx.Exec(queryAdd, amount, toAccount)
	if err != nil {
		return fmt.Errorf("transfer addition failed: %v", err)
	}

	if rows, _ := resultDeduct.RowsAffected(); rows == 0 {
		return fmt.Errorf("sender account not found")
	}
	if rows, _ := resultAdd.RowsAffected(); rows == 0 {
		return fmt.Errorf("receiver account not found")
	}

	_, err = tx.Exec(`
        INSERT INTO transactions (type, from_account, to_account, amount, timestamp, status, account_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"transfer", fromAccount, toAccount, amount, time.Now(), "completed", fromAccount,
	)
	if err != nil {
		return fmt.Errorf("failed to record sender transaction: %v", err)
	}

	_, err = tx.Exec(`
        INSERT INTO transactions (type, from_account, to_account, amount, timestamp, status, account_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		"transfer", fromAccount, toAccount, amount, time.Now(), "completed", toAccount,
	)
	if err != nil {
		return fmt.Errorf("failed to record receiver transaction: %v", err)
	}

	return tx.Commit()
}

func (r *Repository) GetTransactions(accountID string) ([]*account.Transaction, error) {
	query := `
        SELECT id, type, from_account, to_account, amount, timestamp, status 
        FROM transactions 
        WHERE account_id = $1 
        ORDER BY timestamp DESC`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}
	defer rows.Close()

	var transactions []*account.Transaction
	for rows.Next() {
		tx, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *Repository) DeleteAccount(accountID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM sessions WHERE user_id = $1", accountID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions: %v", err)
	}

	_, err = tx.Exec("DELETE FROM accounts WHERE id = $1", accountID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %v", err)
	}

	return tx.Commit()
}

func scanAccount(row *sql.Row) (*account.Account, error) {
	var acc account.Account
	err := row.Scan(
		&acc.ID, &acc.Password, &acc.CVC2, &acc.Balance, &acc.Name,
		&acc.Phone, &acc.Age, &acc.CreatedAt, &acc.ExpiredAt,
	)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func scanTransaction(rows *sql.Rows) (*account.Transaction, error) {
	var tx account.Transaction
	err := rows.Scan(
		&tx.ID,
		&tx.Type,
		&tx.FromAccount,
		&tx.ToAccount,
		&tx.Amount,
		&tx.Timestamp,
		&tx.Status,
	)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func scanAccountFromRows(rows *sql.Rows) (*account.Account, error) {
	var acc account.Account
	err := rows.Scan(
		&acc.ID, &acc.Password, &acc.CVC2, &acc.Balance, &acc.Name,
		&acc.Phone, &acc.Age, &acc.CreatedAt, &acc.ExpiredAt,
	)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
