package api

import (
	"context"
	"encoding/json"
	"fmt"
	"mfp/account"
	"mfp/session"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// структура сервера API
type Server struct {
	accountList    *account.AccountList
	SessionManager *session.SessionManager
	RateLimiter    *RateLimiter
}

// создание нового сервера API
func NewServer(accountList *account.AccountList, sessionManager *session.SessionManager) *Server {
	return &Server{
		accountList:    accountList,
		SessionManager: sessionManager,
		RateLimiter:    NewRateLimiter(3, 10*time.Second),
	}
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIPAddress(r)
		fmt.Printf("Rate limit check for IP: %s\n", ip)

		if !s.RateLimiter.Allow(ip) {
			fmt.Printf("Rate limit EXCEEDED for IP: %s\n", ip) // ← ДОБАВЬ ЛОГ
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}
		fmt.Printf("Rate limit ALLOWED for IP: %s\n", ip)

		next.ServeHTTP(w, r)
	})
}

func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		sess, err := s.SessionManager.GetSession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", sess.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// обработчик создания аккаунта
func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	// start := time.Now()
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// log.Printf("❌ JSON decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// log.Printf("✅ JSON decoded in %v", time.Since(start))

	acc := account.NewAccount(req.Password, req.FirstName, req.Phone, req.Age)
	// log.Printf("✅ Account object created in %v", time.Since(start))

	if err := s.accountList.AddAccount(acc); err != nil {
		// log.Printf("❌ AddAccount error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// log.Printf("✅ Account added to list in %v", time.Since(start))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AccountToResponse(acc))
	// log.Printf("🎉 Total account creation time: %v", time.Since(start))
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

func (s *Server) handleMyDeposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Valid amount required", http.StatusBadRequest)
		return
	}

	if err := s.accountList.Deposit(userID, amount); err != nil {
		http.Error(w, "Error depositing amount:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Deposit successful",
		"id":      userID,
		"amount":  fmt.Sprintf("%.2f", amount),
	})
}

func (s *Server) handleMyWithdraw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Valid amount required", http.StatusBadRequest)
		return
	}

	if err := s.accountList.Withdraw(userID, amount); err != nil {
		http.Error(w, "Error withdrawing amount:\n"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Withdrawal successful",
		"id":      userID,
		"amount":  fmt.Sprintf("%.2f", amount),
	})
}

// обработчик удаления аккаунта
func (s *Server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := s.accountList.RemoveAccount(userID); err != nil {
		http.Error(w, "Error removing account:\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account successfully deleted"})
}

// обработчик перевода средств между аккаунтами
func (s *Server) handleMyTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fromID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

func (s *Server) handleMyTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	acc, err := s.accountList.GetAccount(userID)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acc.Transactions)
}

func (s *Server) handleGetMyAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	acc, err := s.accountList.GetAccount(userID)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(AccountToResponse(acc))
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	acc, err := s.accountList.GetAccount(loginReq.Phone)
	if err != nil {
		http.Error(w, "Invalid phone", http.StatusUnauthorized)
		return
	}

	if !account.CheckPasswordHash(loginReq.Password, acc.Password) {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}

	sessionID := s.SessionManager.CreateSession(acc.ID, r.RemoteAddr, r.UserAgent())

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"user_id": acc.ID,
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	s.SessionManager.DeleteSession(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (s *Server) Start() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(s.rateLimitMiddleware)

	r.Post("/login", s.handleLogin)
	r.Post("/logout", s.handleLogout)
	r.Post("/register", s.handleCreateAccount)

	r.Group(func(r chi.Router) {
		r.Use(s.authMiddleware)

		r.Get("/accounts/me", s.handleGetMyAccount)
		r.Get("/accounts/me/transactions", s.handleMyTransactions)
		r.Post("/accounts/me/deposit", s.handleMyDeposit)
		r.Post("/accounts/me/withdraw", s.handleMyWithdraw)
		r.Post("/accounts/me/transfer", s.handleMyTransfer)
		r.Delete("/accounts/me", s.handleDeleteAccount)
	})

	r.Get("/accounts", s.handleGetAccounts)
	r.Get("/accounts/{id}", s.handleGetAccount)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
