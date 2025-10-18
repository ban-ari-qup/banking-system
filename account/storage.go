package account

import (
	"encoding/json"
	"fmt"
	"os"
)

//const nameFile = "accounts.json"

// сохранение данных в файл
func (al *AccountList) SaveToFile(filename string) error {
	al.mu.RLock()
	defer al.mu.RUnlock()

	data, err := json.MarshalIndent(al.accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal accounts: %v", err)
	}

	if err = os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

// загрузка данных из файла
func (al *AccountList) LoadFromFile(filename string) error {
	al.mu.RLock()
	defer al.mu.RUnlock()

	data, err := os.ReadFile(filename)

	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(data, &al.accounts); err != nil {
		return fmt.Errorf("failed to unmarshal accounts: %v", err)
	}

	al.accountsbyNumber = make(map[string]*Account)
	for _, acc := range al.accounts {
		al.accountsbyNumber[acc.Phone] = acc
	}

	return nil
}
