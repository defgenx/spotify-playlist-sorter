package session

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

// Session represents a user session
type Session struct {
	UserID    string
	Token     *oauth2.Token
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Store manages user sessions in memory
type Store struct {
	sessions map[string]*Session // sessionID -> Session
	userSessions map[string]string // userID -> sessionID
	mu       sync.RWMutex
}

// NewStore creates a new session store
func NewStore() *Store {
	store := &Store{
		sessions: make(map[string]*Session),
		userSessions: make(map[string]string),
	}

	// Start cleanup goroutine
	go store.cleanupExpiredSessions()

	return store
}

// Create creates a new session
func (s *Store) Create(sessionID, userID string, token *oauth2.Token) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := &Session{
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Session expires in 24 hours
	}

	s.sessions[sessionID] = session
	s.userSessions[userID] = sessionID

	return session
}

// Get retrieves a session by ID
func (s *Store) Get(sessionID string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return session, nil
}

// GetByUserID retrieves a session by user ID
func (s *Store) GetByUserID(userID string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessionID, exists := s.userSessions[userID]
	if !exists {
		return nil, ErrSessionNotFound
	}

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return session, nil
}

// Update updates a session's token
func (s *Store) Update(sessionID string, token *oauth2.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}

	session.Token = token
	return nil
}

// Delete deletes a session
func (s *Store) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session, exists := s.sessions[sessionID]; exists {
		delete(s.userSessions, session.UserID)
	}
	delete(s.sessions, sessionID)
}

// DeleteByUserID deletes a session by user ID
func (s *Store) DeleteByUserID(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sessionID, exists := s.userSessions[userID]; exists {
		delete(s.sessions, sessionID)
		delete(s.userSessions, userID)
	}
}

// Refresh extends the session expiration
func (s *Store) Refresh(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}

	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	return nil
}

// cleanupExpiredSessions periodically removes expired sessions
func (s *Store) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()

		for sessionID, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.userSessions, session.UserID)
				delete(s.sessions, sessionID)
			}
		}

		s.mu.Unlock()
	}
}

// Count returns the number of active sessions
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.sessions)
}
