package api

import (
	"encoding/json"
	"fmt"
	"mfp/account"
	"net/http"
	"strconv"
)

type Server struct {
	accountList *account.AccountList
}

func NewServer(accountList *account.AccountList) *Server {
	return &Server{accountList: accountList}
}

// Add more API handlers as needed

func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse and validate the request body
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	acc := account.NewAccount(req.Password, req.FirstName, req.Phone, req.Age)
	if err := acc.Validate(); err != nil {
		http.Error(w, "Error creating account:\n"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.accountList.AddAccount(acc); err != nil {
		http.Error(w, "Error adding account to list:\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(AccountToResponse(acc))
}
func (s *Server) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	accounts := s.accountList.GetAccounts()
	fmt.Printf("DEBUG: handleGetAccounts received %d accounts\n", len(accounts))
	response := make([]AccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		response = append(response, AccountToResponse(acc))
		fmt.Printf("DEBUG: Account: ID='%s', Name='%s', Phone='%s'\n", acc.ID, acc.Name, acc.Phone)
	}
	fmt.Printf("DEBUG: handleGetAccounts received %d accounts\n", len(response))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}
	if err := s.accountList.RemoveAccount(id); err != nil {
		http.Error(w, "Error removing account:\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account successfully deleted"})
}

func (s *Server) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}
	acc, err := s.accountList.GetAccount(id)
	if err != nil {
		http.Error(w, "Error retrieving account:\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AccountToResponse(acc))
}

// func (s *Server) handleGetAccountByNumber(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	number := r.URL.Query().Get("number")
// 	if number == "" {
// 		http.Error(w, "Missing account ID", http.StatusBadRequest)
// 		return
// 	}
// 	acc, err := s.accountList.GetAccount(number)
// 	if err != nil {
// 		http.Error(w, "Error retrieving account:\n"+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(AccountToResponse(acc))
// }

func (s *Server) handleDeposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}
	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		http.Error(w, "Missing deposit amount", http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid deposit amount", http.StatusBadRequest)
		return
	}

	acc, err := s.accountList.GetAccount(id)
	if err != nil {
		http.Error(w, "Error retrieving account:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := acc.Deposit(amount); err != nil {
		http.Error(w, "Error depositing amount:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.accountList.SaveToFile("accounts.json"); err != nil {
		http.Error(w, "Error saving account data:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AccountToResponse(acc))

}

func (s *Server) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}
	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		http.Error(w, "Missing withdraw amount", http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid withdraw amount", http.StatusBadRequest)
		return
	}

	acc, err := s.accountList.GetAccount(id)
	if err != nil {
		http.Error(w, "Error retrieving account:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := acc.Withdraw(amount); err != nil {
		http.Error(w, "Error withdrawing amount:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.accountList.SaveToFile("accounts.json"); err != nil {
		http.Error(w, "Error saving account data:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AccountToResponse(acc))
}

func (s *Server) handleAccountTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.accountList.Transfer(req.From, req.To, req.Amount); err != nil {
		http.Error(w, "Error transferring amount:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transfer successful"})
}

func (s *Server) Start() {
	http.HandleFunc("/accounts", s.handleCreateAccount)
	http.HandleFunc("/accounts/list", s.handleGetAccounts)
	http.HandleFunc("/accounts/delete", s.handleDeleteAccount)
	http.HandleFunc("/accounts/{id}", s.handleGetAccountByID)
	// http.HandleFunc("/accounts/{number}", s.handleGetAccountByNumber)
	http.HandleFunc("/accounts/{id}/deposit", s.handleDeposit)
	http.HandleFunc("/accounts/{id}/withdraw", s.handleWithdraw)
	http.HandleFunc("/accounts/transfer", s.handleAccountTransfer)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
