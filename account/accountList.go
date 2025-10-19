package account

import (
	"fmt"
	"sync"
	"time"
)

// структура для хранения списка аккаунтов
type AccountList struct {
	accounts         map[string]*Account //мапа аккаунтов по ID
	accountsbyNumber map[string]*Account //мапа аккаунтов по номеру телефона
	mu               sync.RWMutex        // для потокобезопасности
}

// создание нового списка аккаунтов
func NewAccountList() *AccountList {
	return &AccountList{
		accounts: make(map[string]*Account), accountsbyNumber: make(map[string]*Account),
	}
}

// добавление аккаунта в список
func (al *AccountList) AddAccount(account *Account) error {
	if err := account.Validate(); err != nil {
		return err
	}

	al.mu.Lock()
	defer al.mu.Unlock()

	if _, exists := al.accounts[account.ID]; exists {
		return fmt.Errorf("account with ID %s already exists", account.ID)
	}
	if _, exists := al.accountsbyNumber[account.Phone]; exists {
		return fmt.Errorf("account with phone %s number already exists", account.Phone)
	}
	al.accounts[account.ID] = account
	al.accountsbyNumber[account.Phone] = account

	go func() {
		al.saveToFile()
	}()
	// if al.saveToFile() != nil {
	// 	log.Printf("❌ Failed to save account list after adding account ID %s", account.ID)
	// }
	return nil
}

// получение всех аккаунтов
func (al *AccountList) GetAccounts() []*Account {
	al.mu.RLock()
	defer al.mu.RUnlock()

	accounts := make([]*Account, 0, len(al.accounts))
	for _, acc := range al.accounts {
		accounts = append(accounts, acc)
	}
	return accounts
}

// получение аккаунта по ID или номеру телефона
func (al *AccountList) GetAccount(id string) (*Account, error) {
	al.mu.RLock()
	defer al.mu.RUnlock()

	return al.findAccount(id)
}

// удаление аккаунта по ID
func (al *AccountList) RemoveAccount(id string) error {
	al.mu.Lock()
	defer al.mu.Unlock()

	acc, exists := al.accounts[id]
	if !exists {
		return fmt.Errorf("account not found")
	}

	delete(al.accountsbyNumber, acc.Phone)
	delete(al.accounts, id)

	return al.saveToFile()
}

// пополнение баланса
func (al *AccountList) Deposit(accountID string, amount float64) error {
	al.mu.Lock()
	defer al.mu.Unlock()

	acc, err := al.findAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found")
	}

	if err := acc.Deposit(amount); err != nil {
		return err
	}
	go func() {
		al.saveToFile()
	}()
	// if err := al.saveToFile(); err != nil {
	// 	return err
	// }
	return nil
}

// снятие средств
func (al *AccountList) Withdraw(accountID string, amount float64) error {
	al.mu.Lock()
	defer al.mu.Unlock()

	acc, err := al.findAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found")
	}

	if err := acc.Withdraw(amount); err != nil {
		return err
	}

	go func() {
		al.saveToFile()
	}()
	// if err := al.saveToFile(); err != nil {
	// 	return err
	// }

	return nil
}

// перевод средств между аккаунтами
func (al *AccountList) Transfer(from string, to string, amount float64) error {
	al.mu.Lock()
	defer al.mu.Unlock()

	fromAcc, err := al.findAccount(from)
	if err != nil {
		return fmt.Errorf("source account not found")
	}

	toAcc, err := al.findAccount(to)
	if err != nil {
		return fmt.Errorf("destination account not found")
	}

	if fromAcc.IsExpired() || toAcc.IsExpired() {
		return fmt.Errorf("source account is expired")
	}

	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if fromAcc.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	fromAcc.Balance -= amount
	toAcc.Balance += amount

	transactionOutAcc := Transaction{
		ID:          len(fromAcc.Transactions) + 1,
		Type:        "transfer_out",
		FromAccount: fromAcc.ID,
		ToAccount:   toAcc.ID,
		Amount:      amount,
		Timestamp:   time.Now(),
		Status:      "completed",
	}
	fromAcc.Transactions = append(fromAcc.Transactions, transactionOutAcc)

	transactionInAcc := Transaction{
		ID:          len(toAcc.Transactions) + 1,
		Type:        "transfer_in",
		FromAccount: fromAcc.ID,
		ToAccount:   toAcc.ID,
		Amount:      amount,
		Timestamp:   time.Now(),
		Status:      "completed",
	}
	toAcc.Transactions = append(toAcc.Transactions, transactionInAcc)

	go func() {
		al.saveToFile()
	}()
	// if err := al.saveToFile(); err != nil {
	// 	return err
	// }
	return nil
}

func (al *AccountList) findAccount(id string) (*Account, error) { // внутренний метод для поиска аккаунта по ID
	if acc, exists := al.accounts[id]; exists {
		return acc, nil
	}
	if acc, exists := al.accountsbyNumber[id]; exists {
		return acc, nil
	}
	return nil, fmt.Errorf("account not found")
}

func (al *AccountList) saveToFile() error {
	return al.SaveToFile("accounts.json")
}
