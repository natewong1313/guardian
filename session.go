package guardian

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

func GenerateSessionToken() string {
	bytes := make([]byte, 20)
	// never throws an error so we can safely ignore
	rand.Read(bytes)
	return base32.StdEncoding.EncodeToString(bytes)
}

func getSessionId(token string) string {
	hashed_token := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hashed_token[:])
}

func CreateSession(token string, user_id string, db DatabaseAdapter) (*Session, error) {
	session_id := getSessionId(token)
	session := &Session{
		ID:        session_id,
		UserID:    user_id,
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	}

	err := db.CreateSession(session)
	return session, err
}

func ValidateSessionToken(token string, db DatabaseAdapter) (*Session, error) {
	session_id := getSessionId(token)
	session, err := db.GetSession(session_id)
	if err != nil || session == nil {
		return session, err
	}
	if time.Now().After(session.ExpiresAt) {
		db.DeleteSession(session_id)
		return nil, errors.New("session has expired.")
	}
	half_expiry := session.ExpiresAt.AddDate(0, 0, 15)
	if time.Now().After(half_expiry) {
		session.ExpiresAt = half_expiry
		db.UpdateSession(session.ID, session.ExpiresAt)
	}
	return session, nil
}

func InvalidateSession(session_id string, db DatabaseAdapter) error {
	return db.DeleteSession(session_id)

}

func InvalidateAllSessions(user_id string, db DatabaseAdapter) error {
	return db.DeleteAllSessions(user_id)
}
