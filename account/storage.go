package account

import (
	"encoding/json"
	"fmt"
	"os"
)

func (al *AccountList) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(al.accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal accounts: %v", err)
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	// Реализация сохранения данных в файл
	return nil
}
func (al *AccountList) LoadFromFile(filename string) error {
	al.mu.RLock()
	defer al.mu.RUnlock()
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	err = json.Unmarshal(data, &al.accounts)
	if err != nil {
		return fmt.Errorf("failed to unmarshal accounts: %v", err)
	}

	al.accountsbyNumber = make(map[string]*Account)
	for _, acc := range al.accounts {
		al.accountsbyNumber[acc.Phone] = acc
	}
	// Реализация загрузки данных из файла
	return nil
}
