package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserAgent    string    `json:"user_agent"`
	IP           string    `json:"ip"`
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
	}

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			sm.CleanupExpiredSessions()
		}
	}()
	return sm
}

func (sm *SessionManager) CreateSession(userID, ip, userAgent string) string {
	sessions := sm.GetUserSessions(userID)
	if len(sessions) > 2 {
		oldestSession := sessions[0]
		for _, sess := range sessions[1:] {
			if sess.CreatedAt.Before(oldestSession.CreatedAt) {
				oldestSession = sess
			}
		}
		sm.DeleteSession(oldestSession.ID)
	}

	sessionID := generateSessionID()
	timestamp := time.Now()

	session := &Session{
		ID:           sessionID,
		UserID:       userID,
		CreatedAt:    timestamp,
		LastActivity: timestamp,
		ExpiresAt:    timestamp.Add(15 * time.Minute),
		UserAgent:    userAgent,
		IP:           ip,
	}

	sm.mu.Lock()
	sm.sessions[sessionID] = session
	sm.mu.Unlock()

	return sessionID
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, error) {
	sm.mu.RLock()
	session, exists := sm.sessions[sessionID]
	sm.mu.RUnlock()

	if session == nil || !exists {
		return nil, fmt.Errorf("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		sm.DeleteSession(sessionID)
		return nil, fmt.Errorf("session expired")
	}

	session.LastActivity = time.Now()
	session.ExpiresAt = session.LastActivity.Add(15 * time.Minute)

	return session, nil
}

func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mu.Lock()
	delete(sm.sessions, sessionID)
	sm.mu.Unlock()
}

func (sm *SessionManager) GetUserSessions(userID string) []*Session {
	result := []*Session{}

	sm.mu.RLock()
	for _, session := range sm.sessions {
		if session.UserID == userID {
			result = append(result, session)
		}
	}
	sm.mu.RUnlock()

	return result
}

func (sm *SessionManager) CleanupExpiredSessions() {
	timestamp := time.Now()
	expiredSessions := []string{}

	sm.mu.RLock()
	for id, session := range sm.sessions {
		if timestamp.After(session.ExpiresAt) {
			expiredSessions = append(expiredSessions, id)
		}
	}
	sm.mu.RUnlock()

	if len(expiredSessions) > 0 {
		sm.mu.Lock()
		for _, id := range expiredSessions {
			delete(sm.sessions, id)
		}
		sm.mu.Unlock()
	}
}

func generateSessionID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
