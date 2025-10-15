package api

import (
	"encoding/json"
	"fmt"
	"mfp/account"
	"net/http"
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

// func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	id := r.URL.Query().Get("id")
// 	if id == "" {
// 		http.Error(w, "Missing account ID", http.StatusBadRequest)
// 		return
// 	}
// 	acc, err := s.accountList.GetAccountByID(id)
// 	if err != nil {
// 		http.Error(w, "Error retrieving account:\n"+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(AccountToResponse(acc))
// }

func (s *Server) Start() {
	http.HandleFunc("/accounts", s.handleCreateAccount)
	http.HandleFunc("/accounts/list", s.handleGetAccounts)
	http.HandleFunc("/acccounts/delete", s.handleDeleteAccount)
	// http.HandleFunc("")
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
