package account

import (
	"fmt"
	"sync"
)

type AccountList struct {
	accounts         map[string]*Account
	accountsbyNumber map[string]*Account
	mu               sync.RWMutex // для потокобезопасности
}

func NewAccountList() *AccountList {
	return &AccountList{
		accounts: make(map[string]*Account), accountsbyNumber: make(map[string]*Account),
	}
}

func (al *AccountList) AddAccount(account *Account) error {
	if err := account.Validate(); err != nil {
		return err
	}

	al.mu.Lock()
	defer al.mu.Unlock()

	if _, exists := al.accounts[account.ID]; exists {
		return fmt.Errorf("account already exists")
	}
	if _, exists := al.accountsbyNumber[account.Phone]; exists {
		return fmt.Errorf("account with this phone number already exists")
	}
	al.accounts[account.ID] = account
	al.accountsbyNumber[account.Phone] = account
	if err := al.SaveToFile("accounts.json"); err != nil {
		return fmt.Errorf("failed to save accounts to file: %v", err)
	}
	return nil
}

func (al *AccountList) GetAccounts() []*Account {
	al.mu.RLock()
	defer al.mu.RUnlock()

	accounts := make([]*Account, 0, len(al.accounts))
	for _, acc := range al.accounts {
		accounts = append(accounts, acc)
	}
	return accounts
}

func (al *AccountList) RemoveAccount(id string) error {
	al.mu.Lock()
	defer al.mu.Unlock()

	if _, exists := al.accounts[id]; !exists {
		return fmt.Errorf("account not found")
	}

	delete(al.accountsbyNumber, al.accounts[id].Phone)
	delete(al.accounts, id)
	if err := al.SaveToFile("accounts.json"); err != nil {
		return fmt.Errorf("failed to save accounts to file: %v", err)
	}
	return nil
}

func (al *AccountList) TransferByID(fromID string, toID string, amount float64) error {
	al.mu.Lock()
	defer al.mu.Unlock()
	fromAcc := al.accounts[fromID]
	if fromAcc == nil {
		return fmt.Errorf("source account not found")
	}
	toAcc := al.accounts[toID]
	if toAcc == nil {
		return fmt.Errorf("destination account not found")
	}
	if fromAcc.IsExpired() {
		return fmt.Errorf("source account is expired")
	}
	if toAcc.IsExpired() {
		return fmt.Errorf("destination account is expired")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if fromAcc.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}
	fromAcc.Balance -= amount
	toAcc.Balance += amount
	if err := al.SaveToFile("accounts.json"); err != nil {
		return fmt.Errorf("failed to save accounts to file: %v", err)
	}
	return nil
}

func (al *AccountList) TransferByNumber(fromNumber string, toNumber string, amount float64) error {
	al.mu.Lock()
	defer al.mu.Unlock()
	fromAcc := al.accountsbyNumber[fromNumber]
	if fromAcc == nil {
		return fmt.Errorf("source account not found")
	}
	toAcc := al.accountsbyNumber[toNumber]
	if toAcc == nil {
		return fmt.Errorf("destination account not found")
	}
	if fromAcc.IsExpired() {
		return fmt.Errorf("source account is expired")
	}
	if toAcc.IsExpired() {
		return fmt.Errorf("destination account is expired")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if fromAcc.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	fromAcc.Balance -= amount
	toAcc.Balance += amount
	if err := al.SaveToFile("accounts.json"); err != nil {
		return fmt.Errorf("failed to save accounts to file: %v", err)
	}
	return nil
}
