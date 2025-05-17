package guardian

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

func GenerateSessionToken() []byte {
	bytes := make([]byte, 20)
	// never throws an error
	rand.Read(bytes)
	return bytes
}

func getSessionId(token []byte) string {
	hashed_token := sha256.Sum256(token)
	return hex.EncodeToString(hashed_token[:])
}

func CreateSession(token []byte, user_id int, db DatabaseAdapter) (*Session, error) {
	session_id := getSessionId(token)
	session := &Session{
		ID:        session_id,
		UserID:    user_id,
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	}
	err := db.CreateSession(session)
	return session, err
}

func ValidateSessionToken(token []byte, db DatabaseAdapter) (*Session, error) {
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

func InvalidateAllSessions(user_id int, db DatabaseAdapter) error {
	return db.DeleteAllSessions(user_id)
}
