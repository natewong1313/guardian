package guardian

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	UpdatedAt time.Time
	ExpiresAt time.Time
}

func GenerateSessionToken() string {
	bytes := make([]byte, 20)
	// never throws an error so we can safely ignore
	rand.Read(bytes)
	return base32.StdEncoding.EncodeToString(bytes)
}

func getSessionId(token string) string {
	hashedToken := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hashedToken[:])
}

func CreateSession(token string, userID string, db DatabaseAdapter, expiry ...int) (*Session, error) {
	sessionID := getSessionId(token)
	expiresIn := 30
	if len(expiry) > 0 {
		expiresIn = expiry[0]
	}
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().AddDate(0, 0, expiresIn),
	}

	err := db.CreateSession(session)
	return session, err
}

func ValidateSessionToken(token string, db DatabaseAdapter) (*Session, error) {
	sessionID := getSessionId(token)
	session, err := db.GetSession(sessionID)
	if err != nil || session == nil {
		return nil, err
	}
	currentTime := time.Now()
	if currentTime.After(session.ExpiresAt) {
		db.DeleteSession(sessionID)
		return nil, errors.New("session has expired")
	}

	expiryDuration := session.ExpiresAt.Sub(session.UpdatedAt)
	// extend the session expiry
	// if theres less than half the expiry duration left
	// this is computed by comparing the updatedAt time and expiresAt
	if session.UpdatedAt.Add(expiryDuration/2).Compare(currentTime) == -1 {
		session.UpdatedAt = time.Now()
		expiryDurationDays := expiryDuration.Hours() / 24
		session.ExpiresAt = currentTime.AddDate(0, 0, int(expiryDurationDays))
		db.UpdateSession(session.ID, session.ExpiresAt)
	}
	return session, nil
}

func InvalidateSession(sessionID string, db DatabaseAdapter) error {
	return db.DeleteSession(sessionID)

}

func InvalidateAllSessions(userID string, db DatabaseAdapter) error {
	return db.DeleteAllSessions(userID)
}
