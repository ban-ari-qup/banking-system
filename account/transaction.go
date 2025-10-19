package account

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"` // deposit, withdrawal, transfer
	FromAccount string    `json:"from_account"`
	ToAccount   string    `json:"to_account"`
	Amount      float64   `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"` // pending, completed, failed
}
