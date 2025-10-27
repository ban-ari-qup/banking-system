package api

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	attempts map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	go rl.cleanUp()
	return rl
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	var validAttemps []time.Time
	for _, attempt := range rl.attempts[key] {
		if attempt.After(windowStart) {
			validAttemps = append(validAttemps, attempt)
		}
	}
	fmt.Printf("Key: %s, Attempts: %d/%d\n", key, len(validAttemps), rl.limit) // ← ДОБАВЬ ЛОГ

	if len(validAttemps) >= rl.limit {
		return false
	}

	validAttemps = append(validAttemps, now)
	rl.attempts[key] = validAttemps
	return true
}

func (rl *RateLimiter) cleanUp() {
	for {
		time.Sleep(1 * time.Minute)
		rl.mu.Lock()
		for key, attempts := range rl.attempts {
			var validAttemps []time.Time
			windowStart := time.Now().Add(-rl.window)
			for _, attempt := range attempts {
				if attempt.After(windowStart) {
					validAttemps = append(validAttemps, attempt)
				}
			}
			if len(validAttemps) == 0 {
				delete(rl.attempts, key)
			} else {
				rl.attempts[key] = validAttemps
			}
		}
		rl.mu.Unlock()
	}
}
