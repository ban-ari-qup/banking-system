package api

import (
	"encoding/json"
	"fmt"
	"mfp/account"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// структура сервера API
type Server struct {
	accountList *account.AccountList
}

// создание нового сервера API
func NewServer(accountList *account.AccountList) *Server {
	return &Server{accountList: accountList}
}

// обработчик создания аккаунта
func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	acc := account.NewAccount(req.Password, req.FirstName, req.Phone, req.Age)
	if err := s.accountList.AddAccount(acc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AccountToResponse(acc))
}

// обработчик получения всех аккаунтов
func (s *Server) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accounts := s.accountList.GetAccounts()
	response := make([]AccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		response = append(response, AccountToResponse(acc))
	}

	json.NewEncoder(w).Encode(response)
}

// обработчик получения аккаунта по ID
func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}
	acc, err := s.accountList.GetAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AccountToResponse(acc))
}

func (s *Server) handleDeposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Valid amount required", http.StatusBadRequest)
		return
	}

	acc, err := s.accountList.GetAccount(id)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	if err := acc.Deposit(amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing account ID", http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Valid amount required", http.StatusBadRequest)
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AccountToResponse(acc))
}

// обработчик удаления аккаунта
func (s *Server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := chi.URLParam(r, "id")
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

// обработчик перевода средств между аккаунтами
func (s *Server) handleAccountTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fromID := chi.URLParam(r, "id")
	if fromID == "" {
		http.Error(w, "Missing source account ID", http.StatusBadRequest)
		return
	}

	var req struct {
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.To == "" {
		http.Error(w, "Destination account required", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	if err := s.accountList.Transfer(fromID, req.To, req.Amount); err != nil {
		http.Error(w, "Error transferring amount:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transfer successful",
		"from":    fromID,
		"to":      req.To,
		"amount":  fmt.Sprintf("%.2f", req.Amount),
	})
}

func (s *Server) Start() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/accounts", s.handleCreateAccount)
	r.Get("/accounts/list", s.handleGetAccounts)
	r.Get("/accounts/{id}", s.handleGetAccount)
	r.Post("/accounts/{id}/deposit", s.handleDeposit)
	r.Post("/accounts/{id}/withdraw", s.handleWithdraw)
	r.Post("/accounts/{id}/transfer", s.handleAccountTransfer)
	r.Delete("/accounts/{id}/delete", s.handleDeleteAccount)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
